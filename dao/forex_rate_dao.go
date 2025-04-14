package dao

import (
	"database/sql"
	"log"

	"github.com/gavinklfong/go-forex-trade-api/model"
)

type ForexRateDaoImpl struct {
	db *sql.DB
}

func NewForexRateDao(db *sql.DB) ForexRateDao {
	return &ForexRateDaoImpl{db: db}
}

func (dao *ForexRateDaoImpl) Insert(booking *model.ForexRateBooking) (int64, error) {
	result, err := dao.db.Exec("INSERT INTO forex_rate_booking(id, timestamp, base_currency, counter_currency, rate, "+
		"trade_action, base_currency_amount, booking_ref, expiry_time, customer_id) VALUES "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		booking.ID, booking.Timestamp, booking.BaseCurrency, booking.CounterCurrency, booking.Rate,
		booking.TradeAction, booking.BaseCurrencyAmount, booking.BookingRef, booking.ExpiryTime, booking.CustomerID)
	if err != nil {
		log.Fatalf("insert booking: %v", err)
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("rows affected error: %v", err)
		return 0, err
	}

	return count, nil
}

func (dao *ForexRateDaoImpl) FindByID(id string) (*model.ForexRateBooking, error) {
	var booking model.ForexRateBooking
	err := dao.db.QueryRow("SELECT id, timestamp, base_currency, counter_currency, rate, "+
		"trade_action, base_currency_amount, booking_ref, expiry_time, customer_id "+
		"FROM forex_rate_booking "+
		"WHERE id=?", id).Scan(&booking.ID, &booking.Timestamp, &booking.BaseCurrency, &booking.CounterCurrency, &booking.Rate,
		&booking.TradeAction, &booking.BaseCurrencyAmount, &booking.BookingRef, &booking.ExpiryTime, &booking.CustomerID)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no booking record with id %v\n", id)
		return nil, nil
	case err != nil:
		log.Fatalf("query error: %v\n", err)
		return nil, err
	default:
		return &booking, nil
	}
}

func (dao *ForexRateDaoImpl) FindByBookingRef(bookingRef string) (*model.ForexRateBooking, error) {
	var booking model.ForexRateBooking
	err := dao.db.QueryRow("SELECT id, timestamp, base_currency, counter_currency, rate, "+
		"trade_action, base_currency_amount, booking_ref, expiry_time, customer_id "+
		"FROM forex_rate_booking "+
		"WHERE booking_ref=?", bookingRef).Scan(&booking.ID, &booking.Timestamp, &booking.BaseCurrency, &booking.CounterCurrency, &booking.Rate,
		&booking.TradeAction, &booking.BaseCurrencyAmount, &booking.BookingRef, &booking.ExpiryTime, &booking.CustomerID)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no booking record with id %v\n", bookingRef)
		return nil, nil
	case err != nil:
		log.Fatalf("query error: %v\n", err)
		return nil, err
	default:
		return &booking, nil
	}
}
