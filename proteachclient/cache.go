package proteachclient

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

// Cache holds the latest fetched Config as an immutable snapshot. Readers
// (Telegram update handlers, the reminder cron job) always see either the
// fully old or fully new config, never a partially-updated one, because
// atomic.Value.Store/Load swap the whole pointer in one step.
type Cache struct {
	client  *Client
	current atomic.Value // holds *Config
}

func NewCache(client *Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get() *Config {
	cfg, _ := c.current.Load().(*Config)
	return cfg
}

func (c *Cache) Refresh(ctx context.Context) error {
	cfg, err := c.client.Fetch(ctx)
	if err != nil {
		return err
	}
	c.current.Store(cfg)
	return nil
}

// StartAutoRefresh launches a background goroutine that refreshes the cache
// every interval. Failures are logged and the previous good snapshot keeps
// serving reads — only the initial Refresh (called separately, before this)
// is expected to fail fast.
func (c *Cache) StartAutoRefresh(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.Refresh(ctx); err != nil {
					log.Printf("proteachclient: refresh failed, keeping last known config: %v", err)
				}
			}
		}
	}()
}
