package utils

import (
	"errors"
	"strconv"
)

func ConvertHexToInt(hex string) (int, error) {
	if len(hex) > 2 {
		value, err := strconv.ParseInt(hex[2:], 16, 64)
		if err != nil {
			return 0, err
		}

		return int(value), nil
	}
	return 0, errors.New("invalid block number")
}
