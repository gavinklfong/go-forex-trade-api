package service

import "github.com/gavinklfong/go-forex-trade-api/dao"

type ForexTradeDealService interface {
}

type ForexTradeDealServiceImpl struct {
	dao              *dao.ForexTradeDealDao
	forexRateService ForexRateService
}

func NewForexTradeDealService(dao *dao.ForexTradeDealDao, forexRateService ForexRateService) ForexTradeDealService {
	return &ForexTradeDealServiceImpl{dao, forexRateService}
}
