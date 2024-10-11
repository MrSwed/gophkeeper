package sync

import (
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type syncList struct {
	sync.Map
	count *atomic.Int64
}

func (s *syncList) Len() int64 {
	return s.count.Load()
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
