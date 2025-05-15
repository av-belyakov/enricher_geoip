package storagetemporary

import (
	"context"
	"time"
)

// GetData получить, по ключу, хранящиеся в кеше данные
func (st *StorageTemporary) GetData(key string) (IncomingRequest, bool) {
	st.mux.RLock()
	defer st.mux.RUnlock()

	if storage, ok := st.cache[key]; ok {
		return storage.incomingRequest, ok
	}

	return IncomingRequest{}, false
}

// SetData добавить данные в кеш
func (st *StorageTemporary) SetData(key string, value IncomingRequest) {
	st.mux.Lock()
	defer st.mux.Unlock()

	st.cache[key] = storage{
		timeExpiry:      time.Now().Add(st.timeToLive),
		incomingRequest: value,
	}
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
