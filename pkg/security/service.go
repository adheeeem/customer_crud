package security

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) Auth(login, password string) (ok bool) {
	var managerPassword string
	ctx := context.Background()

	err := s.pool.QueryRow(ctx, `SELECT password FROM managers WHERE login = $1`, login).Scan(&managerPassword)

	if err != nil {
		log.Print(err)
		return false
	}
	if password != managerPassword {
		log.Print(err)
		return false
	}
	return true
}
