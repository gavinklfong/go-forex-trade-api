package service

import "github.com/gavinklfong/go-rest-api-demo/dao"

type ForexTradeDealService struct {
	dao              *dao.ForexTradeDealDao
	forexRateService *ForexRateService
}

func NewForexTradeDealService(dao *dao.ForexTradeDealDao, forexRateService *ForexRateService) *ForexTradeDealService {
	return &ForexTradeDealService{dao, forexRateService}
}
