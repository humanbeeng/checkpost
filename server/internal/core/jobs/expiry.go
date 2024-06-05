package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/humanbeeng/checkpost/server/internal/endpoint"
	"github.com/robfig/cron/v3"
)

type ExpiredRequestsRemover struct {
	cron          *cron.Cron
	endpointStore endpoint.EndpointStore
}

func NewExpiredRequestsRemover(cron *cron.Cron, endpointStore endpoint.EndpointStore) *ExpiredRequestsRemover {
	return &ExpiredRequestsRemover{
		cron:          cron,
		endpointStore: endpointStore,
	}
}

func (re *ExpiredRequestsRemover) Start() error {
	slog.Info("Starting cron runner")

	_, err := re.cron.AddFunc("@daily", re.deleteExpiredRequests)
	if err != nil {
		slog.Error("unable to register expire requests expirer", "err", err)
		return err
	}

	re.cron.Start()
	return nil
}

func (re *ExpiredRequestsRemover) Stop() context.Context {
	slog.Info("Stopping all cron runners")
	return re.cron.Stop()
}

func (re *ExpiredRequestsRemover) deleteExpiredRequests() {
	slog.Info("Deleting expired requests", "date", time.Now().Local().String())
	re.endpointStore.ExpireRequests(context.Background())
}
