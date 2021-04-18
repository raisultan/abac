package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/raisultan/abac/pkg/delete"
	"github.com/raisultan/abac/pkg/jwt_refresh"
	"github.com/raisultan/abac/pkg/list"
	"github.com/raisultan/abac/pkg/login"
	"github.com/raisultan/abac/pkg/register"
	"github.com/raisultan/abac/pkg/retrieve"
	"github.com/raisultan/abac/pkg/update"
)

const (
	InvalidUserIDErrMsg     = "Invalid user ID"
	InvalidReqPayloadErrMsg = "Invalid request payload"
	UserNotFoundErrMsg      = "User not found"
	InvalidCredsErrMsg      = "Invalid user credentials"
)

func Handler(
	reg register.Service,
	login login.Service,
	ref jwt_refresh.Service,
	lst list.Service,
	retr retrieve.Service,
	upd update.Service,
	del delete.Service,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", listUsers(lst)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", retrieveUser(retr)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", updateUser(upd)).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", deleteUser(del)).Methods("DELETE")

	r.HandleFunc("/register", registerUser(reg)).Methods("POST")
	r.HandleFunc("/login", loginUser(login)).Methods("POST")
	r.HandleFunc("/refresh", refreshUserJWT(ref)).Methods("POST")

	r.Use(loggingMiddleware)

	return r
}

func listUsers(s list.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateAuth(r); err != nil {
			respondWithErrorMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		limit, _ := strconv.Atoi(r.FormValue("limit"))
		offset, _ := strconv.Atoi(r.FormValue("offser"))

		if limit > 10 || limit < 1 {
			limit = 10
		}
		if offset < 0 {
			offset = 0
		}

		users, err := s.ListUsers(limit, offset)
		if err != nil {
			respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		}

		respondWithJSON(w, http.StatusOK, users)
	}
}

func retrieveUser(s retrieve.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateAuth(r); err != nil {
			respondWithErrorMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidUserIDErrMsg)
			return
		}

		u, err := s.RetrieveUser(id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithErrorMessage(w, http.StatusNotFound, UserNotFoundErrMsg)
			default:
				respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		respondWithJSON(w, http.StatusOK, u)
	}
}

func updateUser(s update.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateAuth(r); err != nil {
			respondWithErrorMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidUserIDErrMsg)
			return
		}

		var ur update.UserUpdateRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&ur); err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidReqPayloadErrMsg)
			return
		}
		defer r.Body.Close()

		vErr := validateRequest(ur)
		if vErr != nil {
			respondWithJSON(w, http.StatusBadRequest, vErr.Error())
			return
		}

		ur.ID = id
		u, err := s.UpdateUser(ur)

		if err != nil {
			respondWithErrorMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, u)

	}
}

func deleteUser(s delete.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateAuth(r); err != nil {
			respondWithErrorMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidUserIDErrMsg)
			return
		}

		if err := s.DeleteUser(id); err != nil {
			respondWithErrorMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	}
}

func registerUser(s register.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ur register.UserRegisterRequest
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&ur); err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidReqPayloadErrMsg)
			return
		}
		defer r.Body.Close()

		vErr := validateRequest(ur)
		if vErr != nil {
			respondWithJSON(w, http.StatusBadRequest, vErr.Error())
			return
		}

		u, err := s.RegisterUser(ur)
		if err != nil {
			respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, u)
	}
}

func loginUser(s login.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ur login.UserLoginRequest
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&ur); err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidReqPayloadErrMsg)
			return
		}
		defer r.Body.Close()

		vErr := validateRequest(ur)
		if vErr != nil {
			respondWithJSON(w, http.StatusBadRequest, vErr.Error())
			return
		}

		u, err := s.LoginUser(ur)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithErrorMessage(w, http.StatusUnauthorized, InvalidCredsErrMsg)
			default:
				respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		respondWithJSON(w, http.StatusOK, u)
	}
}

func refreshUserJWT(s jwt_refresh.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ur jwt_refresh.UserJWTRefreshRequest
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&ur); err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, InvalidReqPayloadErrMsg)
			return
		}
		defer r.Body.Close()

		vErr := validateRequest(ur)
		if vErr != nil {
			respondWithJSON(w, http.StatusBadRequest, vErr.Error())
			return
		}

		at, err := s.RefreshJWT(ur)
		if err != nil {
			respondWithErrorMessage(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, at)
	}
}
