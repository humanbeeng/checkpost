package endpoint

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

const (
	pingInterval = 5 * time.Second
	// Pong will be received within few seconds. nextPongWait will mark read deadline till = now() + nextPongWait.
	// so if ping fails to receive a pong, deadline wont be pushed further and deadline will be reached.
	// So ping-pong timeout would be 35 - 30 = 5 seconds.
	nextPongWait = 10 * time.Second
)

type EndpointSession struct {
	sync.RWMutex
	sessionsMap map[string]*WSClient
}

type WSClient struct {
	sessionId string
	endpoint  string
	conn      *websocket.Conn
	manager   *WSManager
	egress    chan EgressMessage
}

func NewWSClient(sessionId string, endpoint string, conn *websocket.Conn, manager *WSManager) *WSClient {
	return &WSClient{
		sessionId: sessionId,
		endpoint:  endpoint,
		conn:      conn,
		manager:   manager,
		egress:    make(chan EgressMessage),
	}
}

func (c *WSClient) readMessages(wg *sync.WaitGroup) {
	slog.Info("Reading from conn", "endpoint", c.endpoint, "session_id", c.sessionId)
	// Cleanup function
	defer func() {
		c.manager.RemoveConn(c.endpoint, c.sessionId)
		wg.Done()
	}()

	c.conn.SetReadDeadline(time.Now().Add(nextPongWait))

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			slog.Error("unable to read message from websocket", "err", err)
			// break and let the cleanup func execute
			break
		}
	}
}

func (c *WSClient) writeMessage(wg *sync.WaitGroup) {
	t := time.NewTicker(pingInterval)

	defer func(t *time.Ticker) {
		t.Stop()
		wg.Done()
		c.manager.RemoveConn(c.endpoint, c.sessionId)
	}(t)

	for {
		select {
		case em, ok := <-c.egress:
			{
				if !ok {
					// Connection has been closed
					if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
						slog.Error("unable to write close message", "endpoint", c.endpoint, "session_id", c.sessionId)
						return
					}
				}

				wm := WSMessage{
					Code:    200,
					Payload: em.Payload,
				}

				err := c.conn.WriteJSON(wm)
				if err != nil {
					slog.Error("unable to write json to connection", "endpoint", c.endpoint, "session_id", c.sessionId)
					return
				}

			}
		case <-t.C:
			{
				if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					slog.Error("unable to send ping", "session_id", c.sessionId, "err", err)
					// Trigger cleanup func
					return
				}
			}
		}
	}
}

func (c *WSClient) pongHandler(data string) error {
	// push deadline further
	c.conn.SetReadDeadline(time.Now().Add(nextPongWait))
	return nil
}

type WSManager struct {
	sync.RWMutex
	endpointSessions map[string]*EndpointSession
}

func NewWSManager() *WSManager {
	return &WSManager{
		endpointSessions: make(map[string]*EndpointSession),
	}
}

func (m *WSManager) AddConn(endpoint string, conn *websocket.Conn) error {
	// Reuse requestId as sessionId
	sessionId := conn.Locals("requestid").(string)

	// Check if there are any existing listeners. No, then create a new sessions map, else just store in existing map
	client := WSClient{
		sessionId: sessionId,
		conn:      conn,
		manager:   m,
		endpoint:  endpoint,
		egress:    make(chan EgressMessage),
	}

	sessions, ok := m.endpointSessions[endpoint]
	if !ok {
		slog.Info("No sessions found", "endpoint", endpoint)
		// No sessions found. Create a new sessions map
		s := EndpointSession{sessionsMap: make(map[string]*WSClient)}
		s.sessionsMap[sessionId] = &client
		m.endpointSessions[endpoint] = &s
	} else {
		// TODO: Plan based limit
		if len(sessions.sessionsMap) > 5 {
			slog.Warn("Number of sessions limit exceeded 5", "endpoint", endpoint)
			return nil
		}

		slog.Info("Found existing sessions", "num_sessions", len(sessions.sessionsMap))
		// Acquire lock on both maps
		sessions.Lock()
		m.Lock()
		sessions.sessionsMap[sessionId] = &client
		m.endpointSessions[endpoint] = sessions
		sessions.Unlock()
		m.Unlock()
	}

	conn.SetPongHandler(client.pongHandler)

	slog.Info("Connection added to manager", "endpoint", endpoint, "session_id", sessionId)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go client.readMessages(&wg)
	go client.writeMessage(&wg)
	wg.Wait()
	return nil
}

func (m *WSManager) RemoveConn(endpoint string, sessionId string) error {
	sessions, ok := m.endpointSessions[endpoint]
	if !ok {
		slog.Warn("No sessions found", "endpoint", endpoint)
		return nil
	}

	if s, ok := sessions.sessionsMap[sessionId]; ok {
		err := s.conn.Close()
		if err != nil {
			slog.Error("unable to close connection", "session_id", sessionId, "err", err)
			return nil
		}

		sessions.Lock()
		defer sessions.Unlock()
		delete(sessions.sessionsMap, sessionId)

		if len(sessions.sessionsMap) == 0 {
			// Remove sessions object itself from endpointSessions
			m.Lock()
			defer m.Unlock()
			delete(m.endpointSessions, endpoint)
		}

		slog.Info("Connection removed", "endpoint", endpoint, "session_id", sessionId, "num_sessions", len(sessions.sessionsMap))
	} else {
		slog.Warn("No session found", "session_id", sessionId)
	}
	return nil
}

type EgressEvent string

const (
	Hook EgressEvent = "hook"
	Err  EgressEvent = "err"
)

type EgressMessage struct {
	Type    EgressEvent     `json:"event"`
	Payload json.RawMessage `json:"payload"`
}
