package routes

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BukhryakovVladimir/vkTest/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const UniqueViolationErr = pq.ErrorCode("23505")

func SignupPerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var person model.Person

	bytes, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bytes, &person)

	person.Username = strings.ToLower(person.Username)

	if !isValidUsername(person.Username) {
		resp, err := json.Marshal("Username should have at least 3 characters and consist only of English letters and digits.")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("Write failed: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if !isValidPassword(person.Password) {
		resp, err := json.Marshal("Password should have at least 8 characters and include both English letters and digits. Special characters optionally.")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("Write failed: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	insertPersonQuery := `INSERT INTO person (username, password, firstName, lastName, sex, birthDate, isAdmin) 
							VALUES ($1::text, $2::text, $3::text, $4::text, $5::text, $6::date, false);`

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(person.Password), 14)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, insertPersonQuery,
		person.Username,
		passwordHash,
		person.FirstName,
		person.LastName,
		person.Sex,
		person.BirthDate,
	)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("signup_userDB QueryRowContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		}

		var pgErr *pq.Error
		if ok := errors.As(err, &pgErr); !ok {
			log.Println("Internal server error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if pgErr.Code == "23505" {
			log.Println("Unique key violation, username already exists: ", err)
			http.Error(w, "Username already exists", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", pgErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal("Signup successful")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Write failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func LoginPerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var person model.Person

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bytes, &person)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	getUserDataQuery := `SELECT id, password FROM person WHERE username = $1::text`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, getUserDataQuery, person.Username)

	var userID, passwordHash string
	if err := row.Scan(&userID, &passwordHash); err != nil {
		resp, err := json.Marshal("Username not found")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("Write failed: %v\n", err)
		}
		return
	}

	if err := row.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(person.Password)); err != nil {
		resp, err := json.Marshal("Incorrect password")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("Write failed: %v\n", err)
		}
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	})

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		resp, err := json.Marshal("Could not login")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("Write failed: %v\n", err)
		}
		return
	}

	tokenCookie := http.Cookie{
		Name:     jwtName,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: false,
	}

	http.SetCookie(w, &tokenCookie)
	resp, err := json.Marshal("Successfully logged in")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

// RegEx. Обязательно латинские буквы, цифры и длина >= 3.
func isValidUsername(username string) bool {
	pattern := "^[a-zA-Z0-9]{3,}$"

	regexpPattern := regexp.MustCompile(pattern)

	return regexpPattern.MatchString(username)
}

// RegEx. Обязательно латинские буквы, цифры и длина >= 8. Опционально специальные символы.
func isValidPassword(password string) bool {
	pattern := `^[a-zA-Z0-9!@#$%^&*()-_=+,.?;:{}|<>]*[a-zA-Z]+[0-9!@#$%^&*()-_=+,.?;:{}|<>]*[0-9]+[a-zA-Z0-9!@#$%^&*()-_=+,.?;:{}|<>]*$`

	regexpPattern := regexp.MustCompile(pattern)

	return regexpPattern.MatchString(password) && len(password) >= 8
}
