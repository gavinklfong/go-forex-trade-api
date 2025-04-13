package integrationtest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"database/sql"

	"github.com/gavinklfong/go-forex-trade-api/dao"
	"github.com/gavinklfong/go-forex-trade-api/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type ForexRateDaoTestSuite struct {
	suite.Suite
	dao            dao.ForexRateDao
	mysqlContainer *mysql.MySQLContainer
	db             *sql.DB
}

func (suite *ForexRateDaoTestSuite) SetupSuite() {
	log.Println("setting up test")
	ctx := context.Background()

	// start up MySQL testcontainer
	var err error
	suite.mysqlContainer, err = startMySQLContainer()
	if err != nil {
		log.Panic("fail to start up MySQL container\n", err)
	}
	log.Println("MySQL container started")

	suite.db, err = openDB(ctx, suite.mysqlContainer)
	if err != nil {
		log.Panic("fail to connect to MySQL container\n", err)
	}
	suite.dao = dao.NewForexRateDao(suite.db)

}

func (suite *ForexRateDaoTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM forex_rate_booking")
	if err != nil {
		log.Panic("fail to clean up data")
	}
}

func (suite *ForexRateDaoTestSuite) TearDownSuite() {
	cleanUp(suite.db, suite.mysqlContainer)
}

func (suite *ForexRateDaoTestSuite) TestInsert() {
	// Given
	duration, err := time.ParseDuration("10m")
	suite.Require().NoError(err, "Failed to parse duration")
	
	booking := model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       "GBP",
			CounterCurrency:    "USD",
			BaseCurrencyAmount: 1000,
			TradeAction:        "BUY",
			CustomerId:         123,
		},
		ID:         "1f648720-3bd3-4c8e-8d00-294516f64bf7",
		Timestamp:  time.Now().UTC(),
		Rate:       0.25,
		BookingRef: "ABCD100",
		ExpiryTime: time.Now().Add(duration).UTC(),
		CustomerID: "f1440302-01ab-4083-88fd-8864ae83d435",
	}

	// When
	count, err := suite.dao.Insert(&booking)
	
	// Then
	suite.Require().NoError(err, "Failed to insert booking")
	suite.Equal(int64(1), count, "Expected one row to be affected")
}

func (suite *ForexRateDaoTestSuite) TestFindByID() {
	// Given
	bookingID := "1f648720-3bd3-4c8e-8d00-294516f64bf7"

	duration, err := time.ParseDuration("10m")
	suite.Require().NoError(err, "Failed to parse duration")
	
	booking := model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       "GBP",
			CounterCurrency:    "USD",
			BaseCurrencyAmount: 1000,
			TradeAction:        "BUY",
			CustomerId:         123,
		},
		ID:         bookingID,
		Timestamp:  time.Now().UTC(),
		Rate:       0.25,
		BookingRef: "ABCD100",
		ExpiryTime: time.Now().Add(duration).UTC(),
		CustomerID: "f1440302-01ab-4083-88fd-8864ae83d435",
	}

	// Insert booking directly to database
	err = insertBooking(suite.db, &booking)
	suite.Require().NoError(err, "Failed to insert test booking")

	// When
	actual, err := suite.dao.FindByID(bookingID)
	
	// Then
	suite.Require().NoError(err, "Failed to retrieve booking")
	suite.Require().NotNil(actual, "Expected to find a booking")
	assertForexRateBookingEqual(suite.T(), actual, &booking)
}

func TestForexRateDaoTestSuite(t *testing.T) {
	suite.Run(t, new(ForexRateDaoTestSuite))
}

func assertForexRateBookingEqual(t *testing.T, a, b *model.ForexRateBooking) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Timestamp.Round(time.Duration(time.Second)), b.Timestamp.Round(time.Duration(time.Second)))
	assert.Equal(t, a.BaseCurrency, b.BaseCurrency)
	assert.Equal(t, a.CounterCurrency, b.CounterCurrency)
	assert.Equal(t, a.Rate, b.Rate)
	assert.Equal(t, a.TradeAction, b.TradeAction)
	assert.Equal(t, a.BaseCurrencyAmount, b.BaseCurrencyAmount)
	assert.Equal(t, a.BookingRef, b.BookingRef)
	assert.Equal(t, a.ExpiryTime.Round(time.Duration(time.Second)), b.ExpiryTime.Round(time.Duration(time.Second)))
	assert.Equal(t, a.CustomerID, b.CustomerID)
}

func insertBooking(db *sql.DB, booking *model.ForexRateBooking) error {
	_, err := db.Exec("INSERT INTO forex_rate_booking(id, timestamp, base_currency, counter_currency, rate, "+
		"trade_action, base_currency_amount, booking_ref, expiry_time, customer_id) VALUES "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		booking.ID, booking.Timestamp, booking.BaseCurrency, booking.CounterCurrency, booking.Rate,
		booking.TradeAction, booking.BaseCurrencyAmount, booking.BookingRef, booking.ExpiryTime, booking.CustomerID)
	if err != nil {
		log.Printf("insert booking: %v", err)
		return err
	}

	return nil
}
