package utils

import (
	"fmt"
	"go-clean-arch/shared/times"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// CronTaskWrapper is wrapper for graceful tasks
func CronTaskWrapper(wg *sync.WaitGroup, task func()) {
	wg.Add(1)
	task()
	wg.Done()
}

func CronCouldRun(startAt, endAt string) (bool, error) {
	now := times.NowJkt()
	startAtTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", times.NowJkt().Format(times.DATE_FORMAT), startAt), times.LocJkt())
	if err != nil {
		return false, errors.Wrap(err, "cannot parse start at")
	}

	endAtTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", times.NowJkt().Format(times.DATE_FORMAT), endAt), times.LocJkt())
	if err != nil {
		return false, errors.Wrap(err, "cannot parse end at")
	}

	if now.Before(startAtTime) || now.After(endAtTime) {
		return false, nil
	}
	return true, nil

}
