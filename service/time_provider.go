package service

import "time"

type TimeProviderImpl struct {
}

func (t *TimeProviderImpl) Now() time.Time {
	return time.Now()
}
