package main

import (
	"database/sql"

	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

const defaultPort = ":8080"

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbUrl string) {
	var err error
	a.DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(defaultPort, a.Router))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	u := user{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var uCreds userLogin
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&uCreds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	u := userLogin{Email: uCreds.Email}
	if err := u.getUserByEmail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(uCreds.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	var uReq userRegisterRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&uReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	passwdMinLength := 6
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > passwdMinLength
	})

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} is not strong enough", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})

	if err := v.Struct(uReq); err != nil {
		fieldErrors := []fieldValidationError{}
		for _, e := range err.(validator.ValidationErrors) {
			fieldError := fieldValidationError{Field: e.Namespace(), Error: e.Translate(trans)}
			fieldErrors = append(fieldErrors, fieldError)
		}
		respondWithJSON(w, http.StatusBadRequest, validationError{Details: fieldErrors})
		return
	}

	exists, err := uReq.checkUserExists(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		if exists == true {
			respondWithError(w, http.StatusBadRequest, "User already exists")
			return
		}
	}

	if err := uReq.register(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	uResp := userRegisterResponse{
		ID:        uReq.ID,
		Email:     uReq.Email,
		FirstName: uReq.FirstName,
		LastName:  uReq.LastName,
	}
	respondWithJSON(w, http.StatusCreated, uResp)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

	a.Router.HandleFunc("/register", a.register).Methods("POST")
	a.Router.HandleFunc("/login", a.login).Methods("POST")
}
