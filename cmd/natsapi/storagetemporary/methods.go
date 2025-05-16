package storagetemporary

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

// GetReq дескриптор запроса
func (st *StorageTemporary) GetReq(key string) (*nats.Msg, bool) {
	st.mux.RLock()
	defer st.mux.RUnlock()

	if storage, ok := st.cache[key]; ok {
		return storage.requestDescriptor, ok
	}

	return nil, false
}

// SetReq дескриптор в кеш
func (st *StorageTemporary) SetReq(key string, value *nats.Msg) {
	st.mux.Lock()
	defer st.mux.Unlock()

	st.cache[key] = storage{
		timeExpiry:        time.Now().Add(st.timeToLive),
		requestDescriptor: value,
	}
}

// DelReq удалить дескриптор из кеша
func (st *StorageTemporary) DelReq(key string) {
	st.mux.Lock()
	defer st.mux.Unlock()

	delete(st.cache, key)
}

// Cancel выполняет очистку кеша
func (st *StorageTemporary) Cancel() {
	st.mux.Lock()
	defer st.mux.Unlock()

	st.cache = map[string]storage{}
}

// DataSize размер кеша
func (st *StorageTemporary) DataSize() int {
	return len(st.cache)
}

func (st *StorageTemporary) deleteOldIncomingRequests(ctx context.Context) {
	go func(ctx context.Context) {
		tick := time.NewTicker(st.timeTick)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
				st.mux.Lock()

				for key, value := range st.cache {
					if value.timeExpiry.Before(time.Now()) {
						delete(st.cache, key)
					}
				}

				st.mux.Unlock()
			}
		}
	}(ctx)
}
