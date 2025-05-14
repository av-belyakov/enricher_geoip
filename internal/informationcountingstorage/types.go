package informationcountingstorage

import (
	"sync"
	"time"
)

type InformationCountingStorage struct {
	mutex         sync.RWMutex
	startTime     time.Time       //время инициализации счетчика
	countingSheet map[string]uint //лист подсчета
}
