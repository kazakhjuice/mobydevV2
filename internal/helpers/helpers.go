package helpers

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

func ServeError(err error, w http.ResponseWriter, errText string, log *slog.Logger, code int) {
	http.Error(w, errText, code)

	log.Error(errText, Err(err))

}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(email)
}

func RetrieveID(r *http.Request) (id int, err error) {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		return 0, errors.New("cannot retrieve ID")
	}

	id, err = strconv.Atoi(idStr)

	if err != nil {
		return 0, errors.New("cannot retrieve ID")
	}

	return id, nil
}
