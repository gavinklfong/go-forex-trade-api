package dao

import (
	"database/sql"
	"log"

	"github.com/gavinklfong/go-forex-trade-api/model"
)

type CustomerDao struct {
	db *sql.DB
}

func NewCustomerDao(db *sql.DB) *CustomerDao {
	return &CustomerDao{db: db}
}

func (dao *CustomerDao) Insert(customer *model.Customer) (int64, error) {
	result, err := dao.db.Exec("INSERT INTO customer(id, name, tier, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		customer.ID, customer.Name, customer.Tier, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		log.Fatalf("insert customer: %v", err)
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("rows affected error: %v", err)
		return 0, err
	}

	return count, nil
}

func (dao *CustomerDao) FindByID(id string) (*model.Customer, error) {
	var customer model.Customer
	err := dao.db.QueryRow("SELECT id, name, tier, created_at, updated_at "+
		"FROM customer WHERE id=?", id).Scan(&customer.ID, &customer.Name, &customer.Tier, &customer.CreatedAt, &customer.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no customer record with id %v\n", id)
		return nil, nil
	case err != nil:
		log.Fatalf("query error: %v\n", err)
		return nil, err
	default:
		return &customer, nil
	}
}

func (dao *CustomerDao) FindByTier(tier int) (result []*model.Customer, err error) {
	rows, err := dao.db.Query("SELECT id, name, tier, created_at, updated_at "+
		"FROM customer WHERE tier=?", tier)
	defer rows.Close()

	if err != nil {
		log.Fatalf("query error: %v\n", err)
		return nil, err
	}

	var customer model.Customer
	for rows.Next() {
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Tier, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			log.Fatalf("query error: %v\n", err)
			return nil, err
		}
		result = append(result, &customer)
	}

	return result, nil
}
