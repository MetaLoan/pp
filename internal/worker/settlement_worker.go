package worker

import (
	"log"
	"time"

	"github.com/MetaLoan/pp/internal/config"
	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/service"
)

// StartSettlementWorker runs a ticker to settle due orders without per-order goroutines.
func StartSettlementWorker(tradeService *service.TradeService) {
	cfg := config.GlobalConfig.Settlement
	interval := time.Second
	if d, err := time.ParseDuration(cfg.Interval); err == nil && d > 0 {
		interval = d
	}
	batch := cfg.BatchSize
	if batch <= 0 {
		batch = 200
	}
	lockKey := cfg.LockKey
	if lockKey == 0 {
		lockKey = 4242
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			locked, err := tryAdvisoryLock(lockKey)
			if err != nil {
				log.Printf("settlement lock error: %v", err)
				continue
			}
			if !locked {
				continue
			}
			if err := tradeService.ProcessDueOrders(batch); err != nil {
				log.Printf("settlement worker error: %v", err)
			}
			if err := releaseAdvisoryLock(lockKey); err != nil {
				log.Printf("settlement unlock error: %v", err)
			}
		}
	}()
}

func tryAdvisoryLock(key int64) (bool, error) {
	var ok bool
	err := core.DB.Raw("SELECT pg_try_advisory_lock(?)", key).Scan(&ok).Error
	return ok, err
}

func releaseAdvisoryLock(key int64) error {
	return core.DB.Exec("SELECT pg_advisory_unlock(?)", key).Error
}
