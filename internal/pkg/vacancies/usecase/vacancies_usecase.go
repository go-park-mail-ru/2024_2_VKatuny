// Package usecase contains usecase for vacancies
package usecase

import (
	"fmt"
	"strconv"
)

//вынуть ошибки в отдельный файл
// Is it really necessary?

var ErrOffsetIsEmpty = fmt.Errorf("query parameter offset is empty")
var ErrOffsetIsNotANumber = fmt.Errorf("query parameter offset isn't a number")
var ErrNumIsEmpty = fmt.Errorf("query parameter num is empty")
var ErrNumIsNotANumber = fmt.Errorf("query parameter num isn't a number")

func GetVacanciesWithOffsetInputCheck(offsetStr, numStr string) (uint64, uint64, error) {
	if offsetStr == "" {
		return 0, 0, ErrOffsetIsEmpty
	}
	offset, err := strconv.Atoi(offsetStr)

	if err != nil {
		return 0, 0, ErrOffsetIsNotANumber
	}

	if numStr == "" {
		return 0, 0, ErrNumIsEmpty
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, 0, ErrNumIsNotANumber
	}
	return uint64(offset), uint64(num), nil
}
