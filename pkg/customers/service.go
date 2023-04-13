package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx, `
	SELECT id,name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM customers")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		log.Print(item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, nil
}

func (s *Service) GetAllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM customers WHERE active")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, nil
}

func (s *Service) Save(ctx context.Context, item *Customer) (*Customer, error) {
	customer := &Customer{}
	if item.ID == 0 {
		err := s.db.QueryRowContext(ctx, `
INSERT INTO customers(name, phone, active, created) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id, name, phone, active, created;`, item.Name, item.Phone, item.Active, item.Created).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	} else {
		log.Print(item)
		err := s.db.QueryRowContext(ctx, `
UPDATE customers SET name = $2, phone = $3, active = $4, created = $5 WHERE id = $1 RETURNING id, name, phone, active, created;`, item.ID, item.Name, item.Phone, item.Active, item.Created).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}
	return customer, nil
}

func (s *Service) RemoveCustomerById(ctx context.Context, id int64) (*Customer, error) {
	customer := &Customer{}
	err := s.db.QueryRowContext(ctx, `
DELETE FROM customers WHERE id=$1 RETURNING *;`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if err != nil {
		log.Print(err)
		return nil, err
	}
	return customer, nil
}

func (s *Service) BlockCustomerById(ctx context.Context, id int64) (*Customer, error) {
	customer := &Customer{}

	err := s.db.QueryRowContext(ctx, `
UPDATE customers SET active=false WHERE id=$1 RETURNING *;`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return customer, nil
}

func (s *Service) UnblockCustomerById(ctx context.Context, id int64) (*Customer, error) {
	customer := &Customer{}

	err := s.db.QueryRowContext(ctx, `
UPDATE customers SET active=true WHERE id=$1 RETURNING *;`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return customer, nil
}
