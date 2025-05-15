package storagetemporary

import (
	"sync"
	"time"
)

// StorageTemporary временное хранилище
type StorageTemporary struct {
	mux        sync.RWMutex
	cache      map[string]storage
	timeTick   time.Duration
	timeToLive time.Duration
}

type storage struct {
	timeExpiry      time.Time
	incomingRequest IncomingRequest
}

// IncomingRequest подробное описание входящих запросов
type IncomingRequest struct {
	IP     []string
	Source string
	TaskId string
}

type OptionsStorageTemporary func(*StorageTemporary) error
