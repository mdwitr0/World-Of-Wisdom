package hashcash

import (
	"context"
	"strconv"
	"time"
)

const (
	futuristicDays = 2 * 24 * time.Hour
	expiredDays    = 28 * 24 * time.Hour
	zeroByte       = '0'
)

type IService interface {
	CheckStamp(ctx context.Context, stampForValidation Stamp) bool
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (r Service) CheckStamp(ctx context.Context, stampForValidation Stamp) bool {
	if stampForValidation.Date > time.Now().Add(futuristicDays).Unix() {
		return false
	}
	if stampForValidation.Date < time.Now().Add(-expiredDays).Unix() {
		return false
	}

	if !stampForValidation.IsSolved() {
		return false // insufficient zeroes
	}

	v, err := strconv.ParseInt(stampForValidation.Rand, 10, 64)
	if err != nil {
		return false // invalid rand
	}

	_, err = r.repository.GetIndicator(ctx, v)

	return err == nil
}
