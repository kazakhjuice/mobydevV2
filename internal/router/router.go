package router

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"mobydevLogin/internal/auth"
	"mobydevLogin/internal/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("your-secret-key") // TODO: keep in env

func NewRouter(db *sql.DB, log *slog.Logger) *mux.Router {
	router := mux.NewRouter()

	dbHandler := auth.NewDBHandler()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		dbHandler.Login(w, r, db, log)
	}).Methods("POST")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		dbHandler.Register(w, r, db, log)
	}).Methods("POST")

	router.Handle("/user/{id:[0-9]+}", authorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbHandler.UserProfile(w, r, db, log)
	}), log)).Methods("GET")

	router.Handle("/user/{id:[0-9]+}", authorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbHandler.UpdateUserProfile(w, r, db, log)
	}), log)).Methods("PUT")

	return router
}

func authorize(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			helpers.ServeError(err, w, "Unauthorized", log, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			helpers.ServeError(err, w, "Unauthorized", log, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.ServeError(err, w, "Unauthorized", log, http.StatusUnauthorized)
			return
		}

		id, err := helpers.RetrieveID(r)
		if err != nil {
			helpers.ServeError(err, w, "Unauthorized", log, http.StatusUnauthorized)
			return
		}

		if int(claims["id"].(float64)) != id {
			helpers.ServeError(errors.New("wrong ID"), w, "Unauthorized", log, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
