package repository

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

func GetSomeVacancies(offset, num int) ([]BD.Vacancy, error) {
	vacanciesTable := BD.Vacancies
	leftBound := offset
	rightBound := offset + num
	// covering cases when offset is out of slice bounds
	if leftBound > int(vacanciesTable.Count) {
		rightBound = leftBound
	} else if rightBound > int(vacanciesTable.Count) {
		rightBound = int(vacanciesTable.Count)
	}
	vacanciesTable.Mutex.RLock()
	vacancies := vacanciesTable.Vacancy[leftBound:rightBound]
	vacanciesTable.Mutex.RUnlock()
	return vacancies, nil
}
