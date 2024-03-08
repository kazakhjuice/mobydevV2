package auth

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"mobydevLogin/internal/helpers"
	"net/http"
)

func (h *dbHandler) UserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB, log *slog.Logger) {
	id, err := helpers.RetrieveID(r)

	if err != nil {
		helpers.ServeError(err, w, "bad requerst", log, http.StatusBadRequest)
		return
	}
	var userDetails UserDetails

	err = db.QueryRow("SELECT name, phone, date_of_birth FROM user_details WHERE id=?", id).Scan(&userDetails.Name, &userDetails.Phone, &userDetails.DateOfBirth)
	if err != nil {
		helpers.ServeError(err, w, "Internal Server Error", log, http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(userDetails)
	if err != nil {
		helpers.ServeError(err, w, "Internal Server Error", log, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *dbHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB, log *slog.Logger) {
	id, err := helpers.RetrieveID(r)

	if err != nil {
		helpers.ServeError(err, w, "bad request", log, http.StatusBadRequest)
		return
	}

	var updatedUserDetails UserDetails
	err = json.NewDecoder(r.Body).Decode(&updatedUserDetails)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// add validatores and regexps
	query := "UPDATE user_details SET"
	var params []interface{}

	if updatedUserDetails.Name != "" {
		query += " name=?,"
		params = append(params, updatedUserDetails.Name)
	}

	if updatedUserDetails.Phone != "" {
		query += " phone=?,"
		params = append(params, updatedUserDetails.Phone)
	}

	if updatedUserDetails.DateOfBirth != "" {
		query += " date_of_birth=?,"
		params = append(params, updatedUserDetails.DateOfBirth)
	}

	if len(params) > 0 {
		query = query[:len(query)-1]
		query += " WHERE id=?"
		params = append(params, id)

		_, err = db.Exec(query, params...)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("User details updated successfully"))
		log.Info("User updated")

	} else {
		w.Write([]byte("No changes made to user details"))
	}
}
