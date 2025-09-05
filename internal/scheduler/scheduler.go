package scheduler

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/jerson/mandrillfordev/internal/config"
	"github.com/jerson/mandrillfordev/internal/mailer"
	"github.com/jerson/mandrillfordev/internal/store"
	"github.com/jerson/mandrillfordev/internal/types"
)

type Scheduler struct {
	cfg   config.Config
	store *store.Store
	stop  chan struct{}
	alive atomic.Bool
}

func NewScheduler(cfg config.Config, st *store.Store) *Scheduler {
	return &Scheduler{cfg: cfg, store: st, stop: make(chan struct{})}
}

func (s *Scheduler) Start() {
	if s.alive.Swap(true) {
		return
	}
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				s.tick()
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	if !s.alive.Swap(false) {
		return
	}
	close(s.stop)
}

func (s *Scheduler) tick() {
	// naive scan of scheduled map
	items := s.store.ListScheduled("")
	now := time.Now()
	for _, m := range items {
		if m.ScheduledAt != nil && !m.ScheduledAt.After(now) {
			// due: attempt send
			go func(mr *types.MessageRecord) {
				if _, ok := s.store.RemoveScheduled(mr.ID); !ok {
					return
				}
				// Send
				if err := mailer.SendMessage(s.cfg, mr.Message, mr.ID, &mr.Raw); err != nil {
					mr.Status = "rejected"
					mr.RejectReason = err.Error()
					log.Printf("scheduled send failed id=%s: %v", mr.ID, err)
				} else {
					now := time.Now()
					mr.SentAt = &now
					mr.Status = "sent"
				}
				s.store.SaveMessage(mr)
			}(m)
		}
	}
}
