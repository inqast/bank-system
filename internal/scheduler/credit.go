package scheduler

import (
	"context"
	"time"
)

type processor interface {
	ProcessPayments(ctx context.Context) error
}

type Scheduler struct {
	processor processor
}

func New(processor processor) *Scheduler {
	return &Scheduler{processor: processor}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.processor.ProcessPayments(context.Background())
			if err != nil {
				return
			}
		}
	}
}
