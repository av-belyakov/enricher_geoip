package storagetemporary

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	once sync.Once
	st   *StorageTemporary
	err  error
)

func New(ctx context.Context, opts ...OptionsStorageTemporary) (*StorageTemporary, error) {
	once.Do(func() {
		st = &StorageTemporary{
			cache: make(map[string]storage),
		}

		for _, opt := range opts {
			if err = opt(st); err != nil {
				return
			}
		}

		//автоматическое удаление старых объектов
		st.deleteOldIncomingRequests(ctx)
	})

	return st, err
}

// WithCacheTimeTick интервал проверки данных в кеше, в секундах от 2 до 15
func WithCacheTimeTick(v int) OptionsStorageTemporary {
	return func(st *StorageTemporary) error {
		if v < 2 || v > 15 {
			return errors.New("the time tick interval of a cache entry should be between 2 and 15 seconds")
		}

		st.timeTick = time.Duration(v) * time.Second

		return nil
	}
}

// WithCacheTTL время жизни для объекта в кеше, в секундах от 3 до 3600
func WithCacheTTL(v int) OptionsStorageTemporary {
	return func(st *StorageTemporary) error {
		if v < 3 || v > 3600 {
			return errors.New("the lifetime of a cache entry should be between 3 and 3600 seconds")
		}

		st.timeToLive = time.Duration(v) * time.Second

		return nil
	}
}
