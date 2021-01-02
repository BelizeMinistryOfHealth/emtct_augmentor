package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"moh.gov.bz/mch/emtct/internal/auth"
	"net/http"
)

type userRoutes struct {
	userStore auth.UserStore
}

func (u userRoutes) UserHandlers(w http.ResponseWriter, r *http.Request) {
	handlerName := "UserHandlers"
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		isAdmin := auth.IsAdmin(token.Permissions)
		if !isAdmin {
			log.WithFields(log.Fields{
				"user":    user,
				"method":  r.Method,
				"handler": handlerName,
			}).Error("non-admin user tried to view list of users")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		users, err := u.userStore.ListUsers()
		if err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"handler": handlerName,
				"method":  r.Method,
			}).WithError(err).Error("failed to retrieve list of users")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.WithFields(log.Fields{
				"users":   users,
				"user":    user,
				"handler": handlerName,
				"method":  r.Method,
			}).WithError(err).Error("failed to decode user response")
			http.Error(w, "failed to decode user response", http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		token := r.Context().Value("user").(auth.JwtToken)
		user := token.Email
		isAdmin := auth.IsAdmin(token.Permissions)
		if !isAdmin {
			log.WithFields(log.Fields{
				"user":    user,
				"method":  r.Method,
				"handler": handlerName,
			}).Error("non-admin user tried to view list of users")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		var userToEdit auth.User
		if err := json.NewDecoder(r.Body).Decode(&userToEdit); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"handler": handlerName,
				"method":  r.Method,
			}).WithError(err).Error("failed to decode body")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := u.userStore.UpdateUser(userToEdit); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"handler": handlerName,
				"method":  r.Method,
				"request": userToEdit,
			}).WithError(err).Error("error updating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(userToEdit); err != nil {
			log.WithFields(log.Fields{
				"user":    user,
				"handler": handlerName,
				"method":  r.Method,
				"request": userToEdit,
			}).WithError(err).Error("failed to encode user")
			http.Error(w, "failed to encode user", http.StatusInternalServerError)
			return
		}
	}
}
