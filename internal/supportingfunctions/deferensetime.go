package supportingfunctions

import "time"

func GetDifference(a, b time.Time) (days, hours, minutes, seconds int) {
	if a.After(b) {
		a, b = b, a
	}

	//количество дней по месяцам
	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// extracting years, months,
	// days of two dates
	year1, month1, day1 := a.Date()
	year2, month2, day2 := b.Date()

	// extracting hours, minutes,
	// seconds of two times
	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()

	totalDays1 := year1*365 + day1

	for i := 0; i < (int)(month1)-1; i++ {
		totalDays1 += monthDays[i]
	}

	//подсчитываем високосные годы с начала года "а" и добавляем это количество
	// дополнительных дней к общему количеству дней	totalDays1 += leapYears(a)
	totalDays2 := year2*365 + day2

	for i := range (int)(month2) - 1 {
		totalDays2 += monthDays[i]
	}

	totalDays2 += leapYears(b)
	days = totalDays2 - totalDays1

	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1

	// если разница в секундах становится меньше 0, прибавляем 60
	// и уменьшаем количество минут
	if seconds < 0 {
		seconds += 60
		minutes--
	}

	// выполнение аналогичных операций для минут и часов
	if minutes < 0 {
		minutes += 60
		hours--
	}

	// выполнение аналогичных операций для часов и дней
	if hours < 0 {
		hours += 24
		days--
	}

	return days, hours, minutes, seconds
}

func leapYears(date time.Time) (leaps int) {
	y, m, _ := date.Date()

	if m <= 2 {
		y--
	}

	leaps = y/4 + y/400 - y/100

	return leaps
}
