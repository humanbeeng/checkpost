package endpoint

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const pongTimeout = 5 * time.Second
const pingInterval = 6 * time.Second

type EndpointSession struct {
	sessionsMap map[string]*WSClient
}

type WSClient struct {
	sessionId string
	endpoint  string
	conn      *websocket.Conn
	manager   *WSManager
	egress    chan WSMessage
}

func NewWSClient(sessionId string, endpoint string, conn *websocket.Conn, manager *WSManager) *WSClient {
	conn.Conn.ReadMessage()
	return &WSClient{
		sessionId: sessionId,
		endpoint:  endpoint,
		conn:      conn,
		manager:   manager,
		egress:    make(chan WSMessage),
	}
}

func (c *WSClient) readMessages() {
	slog.Info("Reading from conn", "endpoint", c.endpoint, "session_id", c.sessionId)
	// Cleanup function
	defer func() {
		c.manager.RemoveConn(c.endpoint, c.sessionId)
	}()

	for {
		_, msg, err := c.conn.Conn.ReadMessage()
		if err != nil {
			slog.Error("unable to read message from websocket", "err", err)
			// break and let the cleanup func execute
			break
		}
		fmt.Println("Received message", string(msg))
	}
}

func (c *WSClient) writeMessage() {
	t := time.NewTicker(pingInterval)
	defer func(t *time.Ticker) {
		t.Stop()
		c.manager.RemoveConn(c.endpoint, c.sessionId)
	}(t)

	for {
		select {
		case <-t.C:
			{
				fmt.Println("Sending ping", "session_id", c.sessionId)
				if err := c.conn.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					slog.Error("unable to send ping", "session_id", c.sessionId, "err", err)
					// Trigger cleanup func
					return
				}
			}
		}
	}
}

func (c *WSClient) pongHandler(data string) error {
	fmt.Println("Pong")
	// reset deadline
	c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	return nil
}

type WSManager struct {
	endpointSessions map[string]*EndpointSession
}

func NewWSManager() *WSManager {
	return &WSManager{
		endpointSessions: make(map[string]*EndpointSession),
	}
}

func (m *WSManager) AddConn(endpoint string, conn *websocket.Conn) error {
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("unable to generate uuid", err)
		return err
	}
	sessionId := id.String()

	// Check if endpoint exists
	// Check if there are any existing listeners. No, then create a new sessions map, else just store in existing map

	client := WSClient{
		sessionId: sessionId,
		conn:      conn,
		manager:   m,
		endpoint:  endpoint,
		egress:    make(chan WSMessage),
	}

	sessions, ok := m.endpointSessions[endpoint]
	if !ok {
		slog.Info("No sessions found")
		// No sessions found. Create a new sessions map
		s := EndpointSession{sessionsMap: make(map[string]*WSClient)}
		s.sessionsMap[sessionId] = &client
		m.endpointSessions[endpoint] = &s
	} else {
		slog.Info("Found existing sessions map")
		sessions.sessionsMap[sessionId] = &client
		m.endpointSessions[endpoint] = sessions
	}

	conn.SetPingHandler(client.pongHandler)
	slog.Info("Connection added to manager", "endpoint", endpoint, "session_id", sessionId)
	go client.readMessages()
	go client.writeMessage()
	return nil
}

func (m *WSManager) RemoveConn(endpoint string, sessionId string) error {
	slog.Info("Connection removed", "endpoint", endpoint, "session_id", sessionId)
	return nil
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
