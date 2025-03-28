package dao

import (
	"database/sql"
	"log/slog"

	"github.com/gavinklfong/go-rest-api-demo/model"
)

type ForexTradeDealDao struct {
	db *sql.DB
}

func NewForexTradeDealDao(db *sql.DB) *ForexTradeDealDao {
	return &ForexTradeDealDao{db: db}
}

func (dao *ForexTradeDealDao) Insert(deal *model.ForexTradeDeal) (int64, error) {
	result, err := dao.db.Exec(`
	INSERT INTO forex_trade_deal(id, ref, timestamp, base_currency, counter_currency, rate, 
	trade_action, base_currency_amount, customer_id) VALUES "
	(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		deal.ID, deal.Ref, deal.Timestamp, deal.BaseCurrency, deal.CounterCurrency, deal.Rate,
		deal.TradeAction, deal.BaseCurrencyAmount, deal.CustomerId)
	if err != nil {
		slog.Error("insert deal: %v", err)
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		slog.Error("rows affected error: %v", err)
		return 0, err
	}

	return count, nil
}

func (dao *ForexTradeDealDao) FindByID(id string) (*model.ForexTradeDeal, error) {
	var deal model.ForexTradeDeal
	err := dao.db.QueryRow(`SELECT id, ref, timestamp, base_currency, counter_currency, rate, 
	trade_action, base_currency_amount, customer_id
	FROM forex_trade_deal WHERE id=?`, id).Scan(&deal.ID, &deal.Ref, &deal.Timestamp,
		&deal.BaseCurrency, &deal.CounterCurrency, &deal.Rate, &deal.TradeAction,
		&deal.BaseCurrencyAmount, &deal.CustomerId)
	switch {
	case err == sql.ErrNoRows:
		slog.Info("no deal record with id %v", id)
		return nil, nil
	case err != nil:
		slog.Error("query error: %v\n", err)
		return nil, err
	default:
		return &deal, nil
	}
}
