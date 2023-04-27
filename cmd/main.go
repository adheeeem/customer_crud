package main

import (
	"context"
	"customer_crud/cmd/app"
	"customer_crud/pkg/customers"
	"customer_crud/pkg/security"
	"encoding/hex"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/dig"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	dsn := "postgres://postgres:password@localhost:5432/bankdb"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}

}
func execute(host string, port string, dsn string) (err error) {
	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		customers.NewService,
		security.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		err = container.Provide(dep)
		if err != nil {
			return nil
		}
	}

	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})

	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}

func check() {
	password := "secret"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	log.Print(hex.EncodeToString(hash))

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))

	if err != nil {
		log.Print("invalid password")
		os.Exit(1)
	}

	print("success")
}
