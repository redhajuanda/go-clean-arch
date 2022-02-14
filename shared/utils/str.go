package utils

import (
	"strconv"

	"github.com/pkg/errors"
)

func StringToFloat(s string) (float64, error) {
	if s == "" {
		s = "0"
	}
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return res, errors.Wrapf(err, "cannot parse to float, s:%s", s)
	}
	return res, nil
}

func StringToInt(s string) (int64, error) {
	if s == "" {
		s = "0"
	}
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot parse to float, s:%s", s)
	}
	return res, nil
}
