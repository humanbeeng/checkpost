package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/humanbeeng/checkpost/server/internal/url"
	"github.com/robfig/cron/v3"
)

type ExpiredRequestsRemover struct {
	cron     *cron.Cron
	urlStore url.UrlStore
}

func NewExpiredRequestsRemover(cron *cron.Cron, urlStore url.UrlStore) *ExpiredRequestsRemover {
	return &ExpiredRequestsRemover{
		cron:     cron,
		urlStore: urlStore,
	}
}

func (re *ExpiredRequestsRemover) Start() error {
	slog.Info("Starting cron runner")

	_, err := re.cron.AddFunc("@daily", re.deleteExpiredRequests)
	if err != nil {
		slog.Error("Unable to register expire requests expirer", "err", err)
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
	re.urlStore.ExpireRequests(context.Background())
}
