package service

import "time"

type TimeProviderImpl struct {
}

func NewTimeProvider() TimeProvider {
	return &TimeProviderImpl{}
}

func (t *TimeProviderImpl) Now() time.Time {
	return time.Now()
}
