// Package usecase contains usecase for vacancies
package usecase

import (
	"fmt"
	"strconv"
)

//вынуть ошибки в отдельный файл
// Is it really necessary?

var ErrOffsetIsNotANumber = fmt.Errorf("query parameter offset isn't a number")
var ErrNumIsNotANumber = fmt.Errorf("query parameter num isn't a number")

func ValidateRequestParams(offsetStr, numStr string) (uint64, uint64, error) {
	var err error
	offset, err1 := strconv.Atoi(offsetStr)

	if err1 != nil {
		offset = 0
		err = ErrOffsetIsNotANumber
	}

	num, err2 := strconv.Atoi(numStr)
	if err2 != nil {
		num = 0
		err = ErrNumIsNotANumber  // previous err will be overwritten
	}
	return uint64(offset), uint64(num), err
}
