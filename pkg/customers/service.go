package customers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")
var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

type Response struct {
	Status     string
	CustomerID int64
	Reason     string
}

type TokenValidation struct {
	StatusCode string   `json:"statusCode"`
	Info       Response `json:"info"`
}

func (s *Service) ValidateCustomerToken(ctx context.Context, token string) (result TokenValidation) {
	var id int64
	var expirationTime time.Time
	var res TokenValidation
	err := s.pool.QueryRow(ctx, `SELECT customer_id, expire from customers_tokens WHERE token = $1`, token).Scan(&id, &expirationTime)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	if err == pgx.ErrNoRows {
		res.StatusCode = http.StatusText(http.StatusNotFound)
		res.Info = Response{
			Status: "fail",
			Reason: "Not Found",
		}
		return res
	}

	if expirationTime.UnixNano() < time.Now().UnixNano() {
		res.StatusCode = http.StatusText(http.StatusBadRequest)
		res.Info = Response{
			Status: "fail",
			Reason: "expired",
		}
	} else {
		res.StatusCode = http.StatusText(http.StatusOK)
		res.Info = Response{
			Status:     "ok",
			CustomerID: id,
		}
	}

	return res

}

func (s *Service) TokenForCustomer(ctx context.Context, phone string, password string) (token string, err error) {
	var id int64
	var hash string
	err = s.pool.QueryRow(ctx, `SELECT id, password from customers WHERE phone = $1`, phone).Scan(&id, &hash)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	log.Print(err)
	if err != nil {
		return "", ErrInternal
	}

	log.Print("Hash:", hash)
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Print(err)
		return "", ErrInvalidPassword
	}
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO customers_tokens(token, customer_id) VALUES($1, $2)`, token, id)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
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
	rows, err := s.pool.Query(ctx, "SELECT * FROM customers")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer rows.Close()

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
	rows, err := s.pool.Query(ctx, "SELECT * FROM customers WHERE active")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer rows.Close()

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
	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
INSERT INTO customers(name, phone, password) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id, name, phone, active, created, password;`, item.Name, item.Phone, hash).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created, &customer.Password)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	} else {
		log.Print(item)
		err := s.pool.QueryRow(ctx, `
UPDATE customers SET name = $2, phone = $3, password = $4 WHERE id = $1 RETURNING id, name, phone, active, created, password;`, item.ID, item.Name, item.Phone, hash).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created, &customer.Password)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}
	return customer, nil
}

func (s *Service) RemoveCustomerById(ctx context.Context, id int64) (*Customer, error) {
	customer := &Customer{}
	err := s.pool.QueryRow(ctx, `
DELETE FROM customers WHERE id=$1 RETURNING *;`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if err != nil {
		log.Print(err)
		return nil, err
	}
	return customer, nil
}

func (s *Service) BlockCustomerById(ctx context.Context, id int64) (*Customer, error) {
	customer := &Customer{}

	err := s.pool.QueryRow(ctx, `
UPDATE customers SET active=false WHERE id=$1 RETURNING *;`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return customer, nil
}

func (s *Service) UnblockCustomerById(ctx context.Context, id int64) (*Customer, error) {
	customer := &Customer{}

	err := s.pool.QueryRow(ctx, `
UPDATE customers SET active=true WHERE id=$1 RETURNING *;`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return customer, nil
}
