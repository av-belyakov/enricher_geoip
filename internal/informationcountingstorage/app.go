// Пакет informationcountingstorage позволяет вести подсчет количества значений с уникальными наименованиями
package informationcountingstorage

import (
	"time"
)

// NewInformationMessageCountingStorage новое временное хранилище подсчета информационных сообщений
func NewInformationMessageCountingStorage() *InformationCountingStorage {
	return &InformationCountingStorage{
		startTime:     time.Now(),
		countingSheet: make(map[string]uint),
	}
}
