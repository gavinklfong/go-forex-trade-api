package integrationtest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"database/sql"

	"math/rand"

	"github.com/gavinklfong/go-forex-trade-api/dao"
	"github.com/gavinklfong/go-forex-trade-api/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type CustomerDaoTestSuite struct {
	suite.Suite
	dao            *dao.CustomerDao
	mysqlContainer *mysql.MySQLContainer
	db             *sql.DB
}

func (suite *CustomerDaoTestSuite) SetupSuite() {
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
	suite.dao = dao.NewCustomerDao(suite.db)

}

func (suite *CustomerDaoTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM customer")
	if err != nil {
		log.Panic("fail to clean up data")
	}
}

func (suite *CustomerDaoTestSuite) TearDownSuite() {
	cleanUp(suite.db, suite.mysqlContainer)
}

func (suite *CustomerDaoTestSuite) TestInsert() {
	customer := buildCustomer()
	count, err := suite.dao.Insert(&customer)
	if err != nil {
		fmt.Println("fail to insert", err)
	}

	assert.Equal(suite.T(), int64(1), count)

}

func (suite *CustomerDaoTestSuite) TestFindByID() {
	customer := buildCustomer()
	err := insertCustomer(suite.db, &customer)
	if err != nil {
		fmt.Println("fail to insert", err)
	}

	var actual *model.Customer
	actual, err = suite.dao.FindByID(customer.ID)
	if err != nil {
		suite.T().Error("fail to retrieve record", err)
	}
	if actual == nil {
		suite.T().Error("No record found")
	}

	assertCustomerEqual(suite.T(), actual, &customer)

}

func TestCustomerDaoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerDaoTestSuite))
}

func assertCustomerEqual(t *testing.T, a, b *model.Customer) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Tier, b.Tier)
	assert.Equal(t, a.CreatedAt.Round(time.Duration(time.Second)), b.CreatedAt.Round(time.Duration(time.Second)))
	assert.Equal(t, a.UpdatedAt.Round(time.Duration(time.Second)), b.UpdatedAt.Round(time.Duration(time.Second)))
}

func buildCustomer() model.Customer {
	faker := faker.New()

	timestamp := time.Now().In(time.UTC)
	uuid := uuid.New()
	return model.Customer{ID: uuid.String(), Name: faker.Person().Name(),
		Tier: rand.Intn(5), CreatedAt: timestamp, UpdatedAt: timestamp}
}

func insertCustomer(db *sql.DB, customer *model.Customer) error {
	_, err := db.Exec("INSERT INTO customer(id, name, tier, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		customer.ID, customer.Name, customer.Tier, customer.CreatedAt, customer.UpdatedAt)

	if err != nil {
		log.Fatalf("insert customer: %v", err)
		return err
	}

	return nil
}
