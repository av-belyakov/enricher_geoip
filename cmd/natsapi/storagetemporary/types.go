package storagetemporary

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// StorageTemporary временное хранилище
type StorageTemporary struct {
	mux        sync.RWMutex
	cache      map[string]storage
	timeTick   time.Duration
	timeToLive time.Duration
}

type storage struct {
	timeExpiry        time.Time
	requestDescriptor *nats.Msg
}

type OptionsStorageTemporary func(*StorageTemporary) error
