package main

import (
	"customer_crud/cmd/app"
	"customer_crud/pkg/customers"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net"
	"net/http"
	"os"
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
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	mux := http.NewServeMux()
	customersSvc := customers.NewService(db)
	server := app.NewServer(mux, customersSvc)
	server.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}
