package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

const defaultPort = ":8080"

var jwtKey = []byte("23a93f6d2673b09c1d3d063cf7a97fc20d0054dfce65c3737455dbec25439938")

type App struct {
	Router     *mux.Router
	DB         *sql.DB
	Validator  *validator.Validate
	Translator ut.Translator
}

func registerCustomTranslation(v *validator.Validate, trans ut.Translator) {
	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} is not strong enough", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})
}

func registerCustomValidation(v *validator.Validate) {
	passwdMinLength := 6
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > passwdMinLength
	})
}

func (a *App) Initialize(dbUrl string) {
	var err error
	a.DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	a.Translator = trans
	a.Validator = validator.New()
	registerCustomValidation(a.Validator)
	registerCustomTranslation(a.Validator, a.Translator)

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(defaultPort, a.Router))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	u := user{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithErrorMessage(w, http.StatusNotFound, "User not found")
		default:
			respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) ValidateRequest(r interface{}) (bool, validationError) {
	if err := a.Validator.Struct(r); err != nil {
		fieldErrors := []fieldValidationError{}
		for _, e := range err.(validator.ValidationErrors) {
			fieldError := fieldValidationError{
				Field: e.Namespace(),
				Error: e.Translate(a.Translator),
			}
			fieldErrors = append(fieldErrors, fieldError)
		}
		return false, validationError{Details: fieldErrors}
	}
	return true, validationError{}
}

func respondWithErrorMessage(w http.ResponseWriter, code int, message string) {
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
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateUser(a.DB); err != nil {
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var uLoginCreds userLoginRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&uLoginCreds); err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	isValid, vError := a.ValidateRequest(uLoginCreds)
	if !isValid {
		respondWithJSON(w, http.StatusBadRequest, vError)
		return
	}

	u := userLoginRequest{Email: uLoginCreds.Email}
	if err := u.getUserByEmail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithErrorMessage(w, http.StatusUnauthorized, "Invalid credentials")
		default:
			respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(uLoginCreds.Password))
	if err != nil {
		respondWithErrorMessage(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	at, err := createAccessToken(uLoginCreds.Email)
	if err != nil {
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	rt, err := createRefreshToken(uLoginCreds.Email)
	if err != nil {
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, userLoginJWTResponse{Access: at, Refresh: rt})
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	var uReq userRegisterRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&uReq); err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	isValid, vError := a.ValidateRequest(uReq)
	if !isValid {
		respondWithJSON(w, http.StatusBadRequest, vError)
		return
	}

	exists, err := uReq.checkUserExists(a.DB)
	if err != nil {
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		if exists == true {
			respondWithErrorMessage(w, http.StatusBadRequest, "User already exists")
			return
		}
	}

	if err := uReq.register(a.DB); err != nil {
		respondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
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

func (a *App) refresh(w http.ResponseWriter, r *http.Request) {
	var jwtRefReq userJWTRefreshRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&jwtRefReq); err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	isValid, vError := a.ValidateRequest(jwtRefReq)
	if !isValid {
		respondWithJSON(w, http.StatusBadRequest, vError)
		return
	}

	at, err := refreshAccessToken(jwtRefReq.Refresh)
	if err != nil {
		respondWithErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, userJWTRefreshResponse{Access: at})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

	a.Router.HandleFunc("/register", a.register).Methods("POST")
	a.Router.HandleFunc("/login", a.login).Methods("POST")
	a.Router.HandleFunc("/refresh", a.refresh).Methods("POST")
}

func createAccessToken(email string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["isAuthorized"] = true
	atClaims["email"] = email
	atClaims["type"] = "access"
	atClaims["exp"] = time.Now().Add(5 * time.Minute).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createRefreshToken(email string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["isAuthorized"] = true
	atClaims["email"] = email
	atClaims["type"] = "refresh"
	atClaims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func refreshAccessToken(a string) (string, error) {
	token, err := jwt.Parse(a, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(
				fmt.Sprintf(
					"unexpected signing method: %v",
					token.Header["alg"],
				),
			)
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", err
		}

		tType, ok := claims["type"].(string)
		if !ok {
			return "", err
		}
		if tType != "refresh" {
			return "", errors.New("Refresh token is expected")
		}

		atClaims := jwt.MapClaims{}
		atClaims["isAuthorized"] = true
		atClaims["email"] = email
		atClaims["type"] = "access"
		atClaims["exp"] = time.Now().Add(5 * time.Minute).Unix()

		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		token, err := at.SignedString(jwtKey)
		if err != nil {
			return "", err
		}

		return token, nil
	}

	return "", errors.New("Invalid refresh token")
}

func extractToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	strArr := strings.Split(bearer, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenStr := extractToken(r)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(
				fmt.Sprintf(
					"unexpected signing method: %v",
					token.Header["alg"],
				),
			)
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func isTokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func extractTokenPayload(r *http.Request) (*tokenPayload, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return nil, err
		}

		exp, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		if err != nil {
			return nil, err
		}

		isAuth, ok := claims["isAuthorized"].(bool)
		if !ok {
			return nil, err
		}

		tType, ok := claims["type"].(string)
		if !ok {
			return nil, err
		}

		tp := tokenPayload{
			Email:        email,
			Expiration:   exp,
			IsAuthorized: isAuth,
			Type:         tType,
		}
		return &tp, nil
	}

	return nil, err
}
