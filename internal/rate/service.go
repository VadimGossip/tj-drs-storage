package rate

import (
	"context"
	"github.com/VadimGossip/tj-drs-storage/pkg/util"
	"time"
)

type Service interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, time.Duration, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, time.Duration, error) {
	var totalDur time.Duration
	bRmsgId, dur, err := s.repo.GetBRmsg(ctx, bRmsgKey{
		gwgrId:    gwgrId,
		direction: dir,
		code:      bNumber,
	}, dateAt)
	if err != nil {
		return 0, 0, dur, err
	}

	aRmsgId, dur, err := s.repo.GetARmsg(ctx, aRmsgKey{
		gwgrId:    gwgrId,
		direction: dir,
		bRmsgId:   bRmsgId,
		code:      aNumber,
	}, dateAt)
	totalDur += dur

	rmsrId, rmsvId, dur, err := s.repo.GetRateRmsvRmsr(ctx, rateKey{
		gwgrId:    gwgrId,
		direction: dir,
		aRmsgId:   aRmsgId,
		bRmsgId:   bRmsgId,
	}, dateAt)
	totalDur += dur
	if err != nil {
		return 0, 0, totalDur, err
	}

	rv, dur, err := s.repo.GetRateValue(ctx, rmsvId)
	totalDur += dur
	if err != nil {
		return 0, 0, totalDur, err
	}

	currencyRate, dur, err := s.repo.GetCurrencyRate(ctx, rv.CurrencyId, dateAt)
	totalDur += dur
	if err != nil {
		return 0, 0, totalDur, err
	}

	return rmsrId, util.RoundFloat(rv.Price*currencyRate, 7), totalDur, nil
}
