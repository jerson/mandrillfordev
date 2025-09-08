package store

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jerson/mandrillfordev/internal/types"
)

type Store struct {
	mu        sync.RWMutex
	messages  map[string]*types.MessageRecord
	scheduled map[string]*types.MessageRecord
	templates map[string]*types.Template
}

func NewStore() *Store {
	return &Store{
		messages:  make(map[string]*types.MessageRecord),
		scheduled: make(map[string]*types.MessageRecord),
		templates: make(map[string]*types.Template),
	}
}

func (s *Store) SaveMessage(m *types.MessageRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[m.ID] = m
}

func (s *Store) GetMessage(id string) (*types.MessageRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.messages[id]
	return m, ok
}

func (s *Store) AddScheduled(m *types.MessageRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scheduled[m.ID] = m
	s.messages[m.ID] = m
}

func (s *Store) RemoveScheduled(id string) (*types.MessageRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.scheduled[id]
	if ok {
		delete(s.scheduled, id)
	}
	return m, ok
}

func (s *Store) GetScheduled(id string) (*types.MessageRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.scheduled[id]
	return m, ok
}

func (s *Store) ListScheduled(to string) []*types.MessageRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*types.MessageRecord, 0, len(s.scheduled))
	for _, m := range s.scheduled {
		if to == "" || containsAddress(m.To, to) {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ScheduledAt.Before(*out[j].ScheduledAt) })
	return out
}

func containsAddress(addresses []string, target string) bool {
	t := strings.ToLower(strings.TrimSpace(target))
	for _, a := range addresses {
		if strings.ToLower(strings.TrimSpace(a)) == t {
			return true
		}
	}
	return false
}

func (s *Store) Search(q string, from, to *time.Time, tags []string, senders []string, limit int) []*types.MessageRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*types.MessageRecord, 0, 64)
	ql := strings.ToLower(q)
	tagset := make(map[string]struct{})
	for _, t := range tags {
		tagset[strings.ToLower(t)] = struct{}{}
	}
	senderset := make(map[string]struct{})
	for _, se := range senders {
		senderset[strings.ToLower(se)] = struct{}{}
	}

	for _, m := range s.messages {
		if from != nil && m.CreatedAt.Before(*from) {
			continue
		}
		if to != nil && m.CreatedAt.After(*to) {
			continue
		}
		if ql != "" {
			if !strings.Contains(strings.ToLower(m.Subject), ql) &&
				!strings.Contains(strings.ToLower(m.From), ql) {
				// search to addresses too
				hit := false
				for _, a := range m.To {
					if strings.Contains(strings.ToLower(a), ql) {
						hit = true
						break
					}
				}
				if !hit {
					continue
				}
			}
		}
		if len(tagset) > 0 {
			ok := false
			for _, t := range m.Tags {
				if _, ex := tagset[strings.ToLower(t)]; ex {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		if len(senderset) > 0 {
			if _, ok := senderset[strings.ToLower(m.From)]; !ok {
				continue
			}
		}
		out = append(out, m)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}

// Messages returns a snapshot of all messages
func (s *Store) Messages() []*types.MessageRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*types.MessageRecord, 0, len(s.messages))
	for _, m := range s.messages {
		out = append(out, m)
	}
	return out
}

// Template store ops
func (s *Store) SaveTemplate(t *types.Template) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.templates[strings.ToLower(t.Name)] = t
}

func (s *Store) GetTemplate(name string) (*types.Template, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.templates[strings.ToLower(name)]
	return t, ok
}

func (s *Store) DeleteTemplate(name string) (*types.Template, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := strings.ToLower(name)
	t, ok := s.templates[key]
	if ok {
		delete(s.templates, key)
	}
	return t, ok
}

func (s *Store) ListTemplates(label string) []*types.Template {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*types.Template, 0, len(s.templates))
	for _, t := range s.templates {
		if label == "" {
			out = append(out, t)
			continue
		}
		// filter by label case-insensitive
		for _, l := range t.Labels {
			if strings.EqualFold(l, label) {
				out = append(out, t)
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
