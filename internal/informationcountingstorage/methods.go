package informationcountingstorage

import (
	"fmt"
	"time"
)

// GetStartTime получить начальное время
func (icr *InformationCountingStorage) GetStartTime() time.Time {
	return icr.startTime
}

// SetStartTime установить начальное время
func (icr *InformationCountingStorage) SetStartTime(v time.Time) {
	icr.startTime = v
}

// GetAllCount получить результаты подсчета всех значений
func (icr *InformationCountingStorage) GetAllCount() map[string]uint {
	icr.mutex.RLock()
	defer icr.mutex.RUnlock()

	newCountingSheet := make(map[string]uint, len(icr.countingSheet))
	for k, v := range icr.countingSheet {
		newCountingSheet[k] = v
	}

	return newCountingSheet
}

// GetCount получить результаты подсчета только выбранного значения
func (icr *InformationCountingStorage) GetCount(elementName string) (uint, error) {
	icr.mutex.RLock()
	defer icr.mutex.RUnlock()

	if v, ok := icr.countingSheet[elementName]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("a value named '%s' was not found.", elementName)
}

// Increase увеличить элемент с заданным наименованием на 1
func (icr *InformationCountingStorage) Increase(elementName string) {
	icr.mutex.Lock()
	defer icr.mutex.Unlock()

	if _, ok := icr.countingSheet[elementName]; !ok {
		icr.countingSheet[elementName] = 1

		return
	}

	icr.countingSheet[elementName] = icr.countingSheet[elementName] + 1
}

// SetCount задать необходимое количество для заданного элемента
func (icr *InformationCountingStorage) SetCount(elementName string, count uint) {
	icr.mutex.Lock()
	defer icr.mutex.Unlock()

	icr.countingSheet[elementName] = count
}
