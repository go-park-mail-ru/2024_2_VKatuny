package repository

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// GetSomeVacancies get some num amount of vacancies from db starting from offset
func GetSomeVacancies(offset, num int) ([]inmemorydb.Vacancy, error) {
	vacanciesTable := inmemorydb.Vacancies
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
