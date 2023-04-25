package middleware

import (
	"log"
	"net/http"
)

func Basic(auth func(login, password string) bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			basicLogin, basicPassword, ok := request.BasicAuth()
			log.Println("BasicLogin:", basicLogin)
			log.Println("BasicPassword:", basicPassword)
			if !ok {
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if !auth(basicLogin, basicPassword) {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(writer, request)
		})
	}
}
