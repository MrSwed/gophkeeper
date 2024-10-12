package sync

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type syncList struct {
	sync.Map
	count     *atomic.Int64
	startTime time.Time
}

func (s *syncList) Len() int64 {
	return s.count.Load()
}

func (s *syncList) KeyQueue() chan string {
	keysQueue := make(chan string)
	go func() {
		syncCount := int64(0)
		s.Range(func(key, _ interface{}) bool {
			keysQueue <- key.(string)
			syncCount++
			if syncCount >= s.Len() {
				close(keysQueue)
			}
			return true
		})
	}()
	return keysQueue
}

// ToSync
// must call for store each key all times, from list of server and from local
// if it already exists, at set from second source, check updated_at and delete
// from list if equal
func (s *syncList) ToSync(key any, updatedAt *timestamp.Timestamp) {
	if value, ok := s.Map.Load(key); ok && value != nil {
		// update: keep or drop from list
		if updatedAt.IsValid() && value == updatedAt {
			s.Map.Delete(key)
			s.count.Add(-1)
			return
		}
		// keep at syncList: need init update
		return
	}
	s.Map.Store(key, updatedAt)
	s.count.Add(1)
}
