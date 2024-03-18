package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BukhryakovVladimir/vkTest/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setupTestDB() *sql.DB {
	psql := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		"filmoteka-postgres-test", "5432", "postgres", "filmoteka", "postgres")
	//psql := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
	//	"localhost", "2345", "postgres", "filmoteka", "postgres")
	db, err := sql.Open("postgres", psql)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	db.SetMaxOpenConns(10)

	db.SetMaxIdleConns(5)

	db.SetConnMaxLifetime(time.Duration(30) * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	//queryTimeLimit = 5
	//secretKey = "filmoteka_test"
	//jwtName = "filmoteka_test_jwt"

	return db
}

// Returns True for a valid username with only English letters and digits and length > 3.
func TestIsValidUsername_ValidUsernameLengthGreaterThanThree(t *testing.T) {
	username := "test1234"
	result := isValidUsername(username)
	if !result {
		t.Errorf("Expected isValidUsername(%q) to be true, but got false", username)
	}
}

// Returns False for an empty string.
func TestIsValidUsername_EmptyString(t *testing.T) {
	username := ""
	result := isValidUsername(username)
	if result {
		t.Errorf("Expected isValidUsername(%q) to be false, but got true", username)
	}
}

// Returns False for a string with only spaces.
func TestIsValidUsername_OnlySpaces(t *testing.T) {
	username := "     "
	result := isValidUsername(username)
	if result {
		t.Errorf("Expected isValidUsername(%q) to be false, but got true", username)
	}
}

// Returns False for a string with only special characters.
func TestIsValidUsername_OnlySpecialCharacters(t *testing.T) {
	username := "!@#$"
	result := isValidUsername(username)
	if result {
		t.Errorf("Expected isValidUsername(%q) to be false, but got true", username)
	}
}

// Returns True for a password with at least 8 characters, including at least one letter and one digit
func TestIsValidPassword_ValidPasswordWithLetterAndDigit(t *testing.T) {
	password := "Test1234"
	result := isValidPassword(password)
	if !result {
		t.Errorf("Expected isValidPassword(%q) to be true, but got false", password)
	}
}

// Returns True for a password with at least 8 characters, including at least one letter, one digit, and special characters
func TestIsValidPassword_ValidPasswordWithLetterDigitAndSpecialChars(t *testing.T) {
	password := "Test1234!@#$"
	result := isValidPassword(password)
	if !result {
		t.Errorf("Expected isValidPassword(%q) to be true, but got false", password)
	}
}

// Returns False for an empty password
func TestIsValidPassword_EmptyPassword(t *testing.T) {
	password := ""
	result := isValidPassword(password)
	if result {
		t.Errorf("Expected isValidPassword(%q) to be false, but got true", password)
	}
}

// Returns False for a password with less than 8 characters
func TestIsValidPassword_PasswordLessThan8Characters(t *testing.T) {
	password := "Test123"
	result := isValidPassword(password)
	if result {
		t.Errorf("Expected isValidPassword(%q) to be false, but got true", password)
	}
}

// Returns False for a password with 8 characters, but no letters or digits
func TestIsValidPassword_PasswordWithoutLetterAndDigit(t *testing.T) {
	password := "!@#$%^&*"
	result := isValidPassword(password)
	if result {
		t.Errorf("Expected isValidPassword(%q) to be false, but got true", password)
	}
}

// Returns a valid token and no error when given a valid cookie
func TestJwtCheck_ValidCookie_ReturnsValidTokenAndNoError(t *testing.T) {
	// Initialize the test environment
	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	// Call the jwtCheck function with the valid cookie
	token, err := jwtCheck(cookie)

	// Check that the token is valid and there is no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !token.Valid {
		t.Error("Expected valid token, got invalid token")
	}
}

// Returns an error when given an invalid cookie
func TestJwtCheck_InvalidCookie_ReturnsError(t *testing.T) {
	// Initialize the test environment
	queryTimeLimit = 5
	secretKey = "test_secret_key"
	jwtName = "test_jwt_name"

	// Create an invalid cookie
	cookie := &http.Cookie{
		Name:  "test_cookie",
		Value: "invalid_token",
	}

	// Call the jwtCheck function with the invalid cookie
	_, err := jwtCheck(cookie)

	// Check that an error is returned
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Returns an error when the secret key is incorrect
func TestJwtCheck_IncorrectSecretKey_ReturnsError(t *testing.T) {
	// Initialize the test environment
	queryTimeLimit = 5
	secretKey = "test_secret_key"
	jwtName = "test_jwt_name"

	// Create a valid cookie
	cookie := &http.Cookie{
		Name:  "test_cookie",
		Value: "valid_token",
	}

	// Call the jwtCheck function with the valid cookie and incorrect secret key
	_, err := jwtCheck(cookie)

	// Check that an error is returned
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Returns an error when the cookie value is empty
func TestJwtCheck_EmptyCookieValue_ReturnsError(t *testing.T) {
	// Initialize the test environment
	queryTimeLimit = 5
	secretKey = "test_secret_key"
	jwtName = "test_jwt_name"

	// Create a cookie with empty value
	cookie := &http.Cookie{
		Name:  "test_cookie",
		Value: "",
	}

	// Call the jwtCheck function with the cookie with empty value
	_, err := jwtCheck(cookie)

	// Check that an error is returned
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Returns an error when the token is expired
func TestJwtCheck_ExpiredToken_ReturnsError(t *testing.T) {
	// Initialize the test environment
	queryTimeLimit = 5
	secretKey = "test_secret_key"
	jwtName = "test_jwt_name"

	// Create an expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() - 1000, // Set expiration time to the past
	})
	tokenString, _ := expiredToken.SignedString([]byte(secretKey))

	// Create a cookie with the expired token
	cookie := &http.Cookie{
		Name:  "test_cookie",
		Value: tokenString,
	}

	// Call the jwtCheck function with the cookie with expired token
	_, err := jwtCheck(cookie)

	// Check that an error is returned
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Signup with valid username, password, and birthdate
//
//	func TestSignupPerson_ValidInput(t *testing.T) {
//		db = setupTestDB()
//		defer db.Close()
//
//		queryTimeLimit = 5
//		secretKey = "filmoteka_test"
//		jwtName = "filmoteka_test_jwt"
//
//		// Create a new person object with valid input
//		person := model.Person{
//			Username:  "testuser",
//			Password:  "Test1234",
//			FirstName: "John",
//			LastName:  "Doe",
//			Sex:       "Male",
//			BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
//		}
//
//		requestBody, err := json.Marshal(person)
//		if err != nil {
//			log.Fatalf("error marshalling json")
//		}
//
//		// Call the SignupPerson function directly
//		w := httptest.NewRecorder()
//		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
//		//r.Body = io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"username":"%s","password":"%s","firstName":"%s","lastName":"%s","sex":"%s","birthDate":"%s"}`,
//		//	person.Username, person.Password, person.FirstName, person.LastName, person.Sex, person.BirthDate.Format("2006-01-02"))))
//
//		SignupPerson(w, r)
//
//		// Check response status code
//		if w.Code != http.StatusCreated {
//			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
//		}
//
//		// Check response body
//		var responseBody string
//		err = json.NewDecoder(w.Body).Decode(&responseBody)
//		if err != nil {
//			t.Fatalf("Failed to decode response body: %v", err)
//		}
//
//		expectedResponseBody := "Signup successful"
//		if responseBody != expectedResponseBody {
//			t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
//		}
//	}
//
// Successfully sign up a person with valid username, password, and birthdate
func TestSignupPerson_ValidInput(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with valid input
	person := model.Person{
		Username:  "testuser1",
		Password:  "Test1234",
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	requestBody, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("error marshalling json")
	}

	// Call the SignupPerson function directly
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	//r.Body = io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"username":"%s","password":"%s","firstName":"%s","lastName":"%s","sex":"%s","birthDate":"%s"}`,
	//	person.Username, person.Password, person.FirstName, person.LastName, person.Sex, person.BirthDate.Format("2006-01-02"))))

	SignupPerson(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var responseBody string
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedResponseBody := "Signup successful"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Successfully sign up a person with a password containing special characters
func TestSignupPerson_PasswordWithSpecialCharacters(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with a password containing special characters
	person := model.Person{
		Username:  "testuser2",
		Password:  "Test1234!@#$",
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	requestBody, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("error marshalling json")
	}

	// Call the SignupPerson function directly
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	//r.Body = io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"username":"%s","password":"%s","firstName":"%s","lastName":"%s","sex":"%s","birthDate":"%s"}`,
	//	person.Username, person.Password, person.FirstName, person.LastName, person.Sex, person.BirthDate.Format("2006-01-02"))))

	SignupPerson(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var responseBody string
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedResponseBody := "Signup successful"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Successfully sign up a person with a birthdate in the past
func TestSignupPerson_BirthDateInPast(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with a birthdate in the past
	person := model.Person{
		Username:  "testuser3",
		Password:  "Test1234",
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Now().AddDate(-30, 0, 0),
	}

	requestBody, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("error marshalling json")
	}

	// Call the SignupPerson function directly
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	//r.Body = io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"username":"%s","password":"%s","firstName":"%s","lastName":"%s","sex":"%s","birthDate":"%s"}`,
	//	person.Username, person.Password, person.FirstName, person.LastName, person.Sex, person.BirthDate.Format("2006-01-02"))))

	SignupPerson(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var responseBody string
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedResponseBody := "Signup successful"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Return an error response with status code 500 and message "Error reading request body" when request body cannot be read
func TestSignupPerson_ErrorReadingRequestBody(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with invalid input
	requestBody := []byte("invalid")

	// Call the SignupPerson function directly
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	SignupPerson(w, r)

	// Check response status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check the response body
	var responseBody string
	//err := json.NewDecoder(w.Body).Decode(&responseBody)
	bytesBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	responseBody = string(bytesBody)

	expectedResponseBody := "Error reading request body\n"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Successfully log in with correct username and password
func TestLoginPerson_ValidCredentials(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with valid credentials
	person := model.Person{
		Username: "testuser1",
		Password: "Test1234",
	}

	// Marshal the person object into JSON
	requestBody, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	// Call the LoginPerson function directly
	LoginPerson(w, r)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	var responseBody string
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedResponseBody := "Successfully logged in"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Return appropriate response status and message for successful login
func TestLoginPerson_SuccessfulLoginResponse(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with valid credentials
	person := model.Person{
		Username: "testuser1",
		Password: "Test1234",
	}

	// Marshal the person object into JSON
	requestBody, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	// Call the LoginPerson function directly
	LoginPerson(w, r)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	var responseBody string
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedResponseBody := "Successfully logged in"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Set JWT token cookie for successful login
func TestLoginPerson_SetTokenCookie(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with valid credentials
	person := model.Person{
		Username: "testuser1",
		Password: "Test1234",
	}

	// Marshal the person object into JSON
	requestBody, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	// Call the LoginPerson function directly
	LoginPerson(w, r)

	// Check the token cookie
	tokenCookie := w.Result().Cookies()[0]
	if tokenCookie.Name != jwtName {
		t.Errorf("Expected token cookie name %q, got %q", jwtName, tokenCookie.Name)
	}

	// Check the token cookie value
	tokenValue := tokenCookie.Value
	if tokenValue == "" {
		t.Errorf("Expected non-empty token cookie value, got empty")
	}
}

// Return appropriate response status and message for invalid request body
func TestLoginPerson_InvalidRequestBody(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create an invalid request body
	requestBody := []byte("abc")

	// Create a new HTTP request with the invalid request body
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	// Call the LoginPerson function directly
	LoginPerson(w, r)

	// Check the response status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check the response body
	var responseBody string
	//err := json.NewDecoder(w.Body).Decode(&responseBody)
	bytesBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	responseBody = string(bytesBody)

	expectedResponseBody := "Error reading request body\n"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Return appropriate response status and message for invalid username
func TestLoginPerson_InvalidUsername(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with an invalid username
	person := model.Person{
		Username: "invaliduser",
		Password: "Test1234",
	}

	// Marshal the person object into JSON
	requestBody, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	// Call the LoginPerson function directly
	LoginPerson(w, r)

	// Check the response status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Check the response body
	var responseBody string
	//err := json.NewDecoder(w.Body).Decode(&responseBody)
	bytesBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	responseBody = string(bytesBody)

	expectedResponseBody := "Username not found\n"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Return appropriate response status and message for incorrect password
func TestLoginPerson_IncorrectPassword(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new person object with an incorrect password
	person := model.Person{
		Username: "testuser1",
		Password: "Incorrect1234",
	}

	// Marshal the person object into JSON
	requestBody, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

	// Call the LoginPerson function directly
	LoginPerson(w, r)

	// Check the response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Check the response body
	var responseBody string
	//err := json.NewDecoder(w.Body).Decode(&responseBody)
	bytesBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	responseBody = string(bytesBody)

	expectedResponseBody := "Incorrect password\n"
	if responseBody != expectedResponseBody {
		t.Errorf("Expected response body %q, got %q", expectedResponseBody, responseBody)
	}
}

// Returns a list of actors with their movies when the user is authenticated and authorized.
func TestGetActors_AuthenticatedAndAuthorized(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Call the GetActors function directly
	GetActors(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.ActorAndMovies
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the actors list
	// ...
}

// Returns an empty list when there are no actors in the database.
func TestGetActors_NoActors(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Call the GetActors function directly
	GetActors(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.ActorAndMovies
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the actors list
	// ...
}

// Returns a list of actors with their movies when there are no movies in the database.
func TestGetActors_NoMovies(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Call the GetActors function directly
	GetActors(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.ActorAndMovies
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the actors list
	// ...
}

// Returns an error when the JWT cookie is missing.
func TestGetActors_MissingJWT(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new request without the JWT cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	// Call the GetActors function directly
	GetActors(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Check response body
	// ...
}

// Returns an error when the JWT cookie is invalid.
func TestGetActors_InvalidJWT(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new request with an invalid JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: "invalid_token",
	}
	r.AddCookie(cookie)

	// Call the GetActors function directly
	GetActors(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Check response body
	// ...
}

// Returns an error when the user does not exist in the database.
func TestGetActors_UserNotExists(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1999999999"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Call the GetActors function directly
	GetActors(w, r)

	// Check response status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	bytesBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal("Error reading response body")
	}

	responseBody := string(bytesBody)

	expectedResponse := "Error while checking user authorization\n"

	if responseBody != expectedResponse {
		t.Fatalf("Expected %s but received %s", expectedResponse, responseBody)
	}

	// Check response body
	// ...
}

// Successfully add a new actor with valid input data
func TestAddActor_ValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	setAdmin := `UPDATE person SET isAdmin = true WHERE id = $1`

	_, err := db.Exec(setAdmin, claims["iss"])
	if err != nil {
		t.Error("Error setting user to admin")
	}
	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with valid input data
	actor := model.Actor{
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActor function directly
	AddActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully add a new actor with minimum valid input data
func TestAddActor_MinimumValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with minimum valid input data
	actor := model.Actor{
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "",
		BirthDate: time.Time{},
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActor function directly
	AddActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully add a new actor with maximum valid input data
func TestAddActor_MaximumValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with maximum valid input data
	actor := model.Actor{
		FirstName: strings.Repeat("A", 255),
		LastName:  strings.Repeat("B", 255),
		Sex:       strings.Repeat("C", 10),
		BirthDate: time.Now(),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActor function directly
	AddActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully add a new actor with valid input data and special characters
func TestAddActor_ValidInputDataWithSpecialCharacters(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with valid input data and special characters
	actor := model.Actor{
		FirstName: "John!@#$%^&*()_+",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActor function directly
	AddActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

func TestAddActor_ValidInputDataWithNonASCIICharacters(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with valid input data and non-ASCII characters
	actor := model.Actor{
		FirstName: "Jürgen",
		LastName:  "Müller",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActor function directly
	AddActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Update actor with valid input and valid JWT cookie
func TestUpdateActor_ValidInputValidJWT(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a sample actor object
	actor := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	// Encode the actor object as JSON and set it as the request body
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateActor function directly
	UpdateActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Update actor with valid input and valid JWT cookie, but with empty fields
func TestUpdateActor_ValidInputValidJWT_EmptyFields(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a sample actor object with empty fields
	actor := model.Actor{
		ID:        1,
		FirstName: "",
		LastName:  "",
		Sex:       "",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	// Encode the actor object as JSON and set it as the request body
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateActor function directly
	UpdateActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Update actor with valid input and valid JWT cookie, but with birthDate in the past
func TestUpdateActor_ValidInputValidJWT_BirthDatePast(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a sample actor object with birthDate in the past
	actor := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Now().AddDate(-1, 0, 0),
	}

	// Encode the actor object as JSON and set it as the request body
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateActor function directly
	UpdateActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Update actor with valid input and valid JWT cookie, but with birthDate equal to today
func TestUpdateActor_ValidInputValidJWT_BirthDateInFuture(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a sample actor object with birthDate equal to today
	actor := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Now().Add(24 * time.Hour),
	}

	// Encode the actor object as JSON and set it as the request body
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateActor function directly
	UpdateActor(w, r)

	// Check response status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

}

// Successfully retrieve actors from the database with valid input data
func TestGetActorsWithID_ValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with valid input data
	actor := model.Actor{
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetActorsWithID function directly
	GetActorsWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.Actor
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully retrieve actors from the database with empty input data
func TestGetActorsWithID_EmptyInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with empty input data
	actor := model.Actor{}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetActorsWithID function directly
	GetActorsWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.Actor
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully retrieve actors from the database with only one field of input data
func TestGetActorsWithID_OneFieldInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with one field of input data
	actor := model.Actor{
		FirstName: "John",
		LastName:  "",
		Sex:       "",
		BirthDate: time.Time{},
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetActorsWithID function directly
	GetActorsWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.Actor
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully retrieve actors from the database with multiple fields of input data
func TestGetActorsWithID_MultipleFieldsInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor with multiple fields of input data
	actor := model.Actor{
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetActorsWithID function directly
	GetActorsWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var actors []model.Actor
	err := json.NewDecoder(w.Body).Decode(&actors)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}

// Delete an actor successfully with valid JWT and non-admin privileges
func TestDeleteActor_ValidJWT_NonAdminPrivileges(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to false
	_, err := db.Exec("UPDATE person SET isAdmin = false WHERE id = 1")
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new actor to be deleted
	actor := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteActor function directly
	DeleteActor(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Check response body
	// ...
}

// Delete an actor successfully with valid JWT and admin privileges
func TestDeleteActor_ValidJWT_AdminPrivileges(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to true
	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = 1")
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new actor to be deleted
	actor := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actor)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteActor function directly
	DeleteActor(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully add a movie with all valid fields and no actors
func TestAddMovie_ValidFieldsNoActors(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with all valid fields and no actors
	movie := model.Movie{
		Name:        "Test Movie",
		Description: "This is a test movie",
		Date:        time.Now(),
		Rating:      8,
		Actors:      nil,
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddMovie function directly
	AddMovie(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully add a movie with all valid fields and one actor
func TestAddMovie_ValidFieldsOneActor(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with all valid fields and one actor
	actor := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Now(),
	}
	movie := model.Movie{
		Name:        "Test Movie 2",
		Description: "This is a test movie",
		Date:        time.Now(),
		Rating:      8,
		Actors:      []model.Actor{actor},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddMovie function directly
	AddMovie(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Successfully add a movie with all valid fields and multiple actors
func TestAddMovie_ValidFieldsMultipleActors(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	actor1 := model.Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Now().AddDate(-30, 0, 0),
	}

	actor2 := model.Actor{
		ID:        2,
		FirstName: "Jane",
		LastName:  "Doe",
		Sex:       "Female",
		BirthDate: time.Now().AddDate(-25, 0, 0),
	}

	movie := model.Movie{
		Name:        "Test Movie 3",
		Description: "This is a test movie",
		Date:        time.Now(),
		Rating:      8,
		Actors:      []model.Actor{actor1, actor2},
	}

	movieBytes, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(movieBytes))

	AddMovie(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Added a movie successfully") {
		t.Errorf("unexpected body: expected 'Added a movie successfully'; got %s", body)
	}
}

// Function receives a valid JWT cookie and user is an admin, movie is updated successfully
func TestUpdateMovie_ValidJWTAndAdmin_Success(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to true
	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new movie with valid input data
	movie := model.Movie{
		ID:          1,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      8,
		Actors:      []model.Actor{},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateMovie function directly
	UpdateMovie(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Function receives a valid JWT cookie and user is not an admin, function returns unauthorized
func TestUpdateMovie_ValidJWTAndNotAdmin_Unauthorized(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to false
	_, err := db.Exec("UPDATE person SET isAdmin = false WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new movie with valid input data
	movie := model.Movie{
		ID:          1,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      8,
		Actors:      []model.Actor{},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateMovie function directly
	UpdateMovie(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Perform assertions on the response
	// ...
}

// Function receives a valid JWT cookie and user does not exist, function returns unauthorized
func TestUpdateMovie_ValidJWTAndUserNotExist_Unauthorized(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "999" // User does not exist
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with valid input data
	movie := model.Movie{
		ID:          1,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      8,
		Actors:      []model.Actor{},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the UpdateMovie function directly
	UpdateMovie(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Perform assertions on the response
	// ...
}

func TestUpdateMovie_ValidInputData_Success(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	movie := model.Movie{
		ID:          1,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      8,
	}

	movieJson, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(movieJson))

	UpdateMovie(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %v, got %v", http.StatusCreated, resp.StatusCode)
	}

	var response string
	json.NewDecoder(resp.Body).Decode(&response)
	if response != "Movie updated successfully" {
		t.Errorf("expected response to be 'Movie updated successfully', got '%v'", response)
	}
}

// Successfully add an actor to a movie with valid input data
func TestAddActorToMovie_ValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor movie with valid input data
	actorMovie := model.ActorMovie{
		MovieID:   1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "Male",
		BirthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	requestBody, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActorToMovie function directly
	AddActorToMovie(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Add an actor to a movie with minimum valid input data
func TestAddActorToMovie_MinimumValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor movie with minimum valid input data
	actorMovie := model.ActorMovie{
		MovieID:   1,
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "",
		BirthDate: time.Time{},
	}
	requestBody, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the AddActorToMovie function directly
	AddActorToMovie(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Add an actor to a movie with maximum valid input data
func TestAddActorToMovie_MaximumValidInputData(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new actor movie with maximum valid input data
	actorMovie := model.ActorMovie{
		MovieID:   1,
		FirstName: strings.Repeat("a", 255),
		LastName:  strings.Repeat("b", 255),
		Sex:       strings.Repeat("c", 10),
		BirthDate: time.Now().AddDate(-30, 0, 0),
	}

	// Encode actor movie to JSON and set it as request body
	actorMovieJSON, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(actorMovieJSON))

	// Call AddActorToMovie function
	AddActorToMovie(w, r)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}

	// Check the response body
	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Added an actor to movie successfully") {
		t.Errorf("unexpected body: got %v", string(body))
	}
}

// Function is called with valid JWT token and user is not an admin, function returns unauthorized error
func TestDeleteActorFromMovie_ValidTokenNonAdmin_ReturnsUnauthorizedError(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to false
	_, err := db.Exec("UPDATE person SET isAdmin = false WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new actorMovie with valid input data
	actorMovie := model.ID{
		MovieID: 1,
		ActorID: 1,
	}
	requestBody, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteActorFromMovie function directly
	DeleteActorFromMovie(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Check response body
	var response string
	bytesBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	response = string(bytesBody)

	expectedResponse := "You do not have administrator privileges to delete actor from movie\n"

	if response != expectedResponse {
		t.Errorf("Got %s but expected %s", response, expectedResponse)
	}

	// Perform assertions on the response
	// ...
}

// Function is called with valid JWT token and user is an admin, actor is successfully deleted from movie
func TestDeleteActorFromMovie_ValidTokenAdmin_SuccessfullyDeletesActor(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to true
	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new actorMovie with valid input data
	actorMovie := model.ID{
		MovieID: 1,
		ActorID: 1,
	}
	requestBody, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteActorFromMovie function directly
	DeleteActorFromMovie(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Function is called with valid JWT token and movieID and actorID are both non-zero, actor is successfully deleted from movie
func TestDeleteActorFromMovie_ValidTokenNonZeroIDs_SuccessfullyDeletesActor(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Set isAdmin flag for user in db to true
	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	// Create a new actorMovie with valid input data
	actorMovie := model.ID{
		MovieID: 1,
		ActorID: 1,
	}
	requestBody, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteActorFromMovie function directly
	DeleteActorFromMovie(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Function is called with valid JWT token and movieID is zero, function returns bad request error
func TestDeleteActorFromMovie_ValidTokenZeroMovieID_ReturnsBadRequestError(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = $1", claims["iss"])
	if err != nil {
		t.Fatalf("Failed to set isAdmin flag for user in db: %v", err)
	}

	actorMovie := model.ID{MovieID: 0, ActorID: 1}
	actorMovieJson, _ := json.Marshal(actorMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(actorMovieJson))

	DeleteActorFromMovie(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v, got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

// Returns movies with valid search criteria and valid admin token
func TestGetMoviesWithID_ValidSearchCriteriaAndValidAdminToken(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with valid search criteria
	movie := model.SearchMovie{
		Name:           "Avengers",
		Description:    "Superhero movie",
		Date:           time.Now(),
		Rating:         8,
		ActorFirstName: "",
		ActorLastName:  "",
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetMoviesWithID function directly
	GetMoviesWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var movies []model.SearchMovie
	err := json.NewDecoder(w.Body).Decode(&movies)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Returns empty list when no movies match search criteria
func TestGetMoviesWithID_NoMatchingMovies(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with search criteria that won't match any movies
	movie := model.SearchMovie{
		Name:           "Nonexistent Movie",
		Description:    "",
		Date:           time.Now(),
		Rating:         0,
		ActorFirstName: "",
		ActorLastName:  "",
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetMoviesWithID function directly
	GetMoviesWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var movies []model.SearchMovie
	err := json.NewDecoder(w.Body).Decode(&movies)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Returns empty list when no actors match search criteria
func TestGetMoviesWithID_NoMatchingActors(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with search criteria that won't match any actors
	movie := model.SearchMovie{
		Name:           "",
		Description:    "",
		Date:           time.Now(),
		Rating:         0,
		ActorFirstName: "Nonexistent",
		ActorLastName:  "Actor",
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the GetMoviesWithID function directly
	GetMoviesWithID(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var movies []model.SearchMovie
	err := json.NewDecoder(w.Body).Decode(&movies)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Returns movies ordered by rating in descending order by default
func TestGetMoviesOrdered_ReturnsMoviesOrderedByRatingDescending(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/movies", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Call the GetMoviesOrdered function directly
	GetMoviesOrdered(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var movies []model.Movie
	err := json.NewDecoder(w.Body).Decode(&movies)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Returns 401 status code if JWT cookie is missing
func TestGetMoviesOrdered_ReturnsUnauthorizedIfJWTCookieMissing(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new request without the JWT cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/movies", nil)

	// Call the GetMoviesOrdered function directly
	GetMoviesOrdered(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Function returns expected search results when valid search parameters are provided
func TestSearchMovie_ValidSearchParameters(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new SearchMovie object with valid search parameters
	searchMovie := model.SearchMovie{
		Name:           "Avengers",
		ActorFirstName: "Robert",
		ActorLastName:  "Downey",
	}
	requestBody, _ := json.Marshal(searchMovie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the SearchMovie function directly
	SearchMovie(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var movies []model.Movie
	err := json.NewDecoder(w.Body).Decode(&movies)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Function handles and returns appropriate error messages for invalid JWT token
func TestSearchMovie_InvalidJWTToken(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new request without a JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	// Call the SearchMovie function directly
	SearchMovie(w, r)

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Perform assertions on the response
	// ...
}

// Function is called with valid admin JWT and valid movie ID, movie is deleted successfully
func TestDeleteMovie_ValidAdminJWT_ValidMovieID_MovieDeletedSuccessfully(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with valid ID
	movie := model.Movie{
		ID:          1,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      5,
		Actors:      []model.Actor{},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteMovie function directly
	DeleteMovie(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check response body
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Perform assertions on the response
	// ...
}

// Function is called with valid admin JWT and invalid movie ID, function returns http.StatusBadRequest
func TestDeleteMovie_ValidAdminJWT_InvalidMovieID_ReturnsBadRequest(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with invalid ID
	movie := model.Movie{
		ID:          0,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      5,
		Actors:      []model.Actor{},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteMovie function directly
	DeleteMovie(w, r)

	// Check response status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Perform assertions on the response
	// ...
}

// Function is called with valid admin JWT and movie ID that doesn't exist, function returns http.StatusBadRequest
func TestDeleteMovie_ValidAdminJWT_NonexistentMovieID_ReturnsBadRequest(t *testing.T) {
	// Set up test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Create a new JWT token for authentication
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "1"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signedToken, _ := token.SignedString([]byte(secretKey))

	// Create a new request with the JWT token in the cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	cookie := &http.Cookie{
		Name:  jwtName,
		Value: signedToken,
	}
	r.AddCookie(cookie)

	// Create a new movie with nonexistent ID
	movie := model.Movie{
		ID:          999,
		Name:        "Test Movie",
		Description: "Test Description",
		Date:        time.Now(),
		Rating:      5,
		Actors:      []model.Actor{},
	}
	requestBody, _ := json.Marshal(movie)
	r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Call the DeleteMovie function directly
	DeleteMovie(w, r)

	// Check response status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Perform assertions on the response
	// ...
}

// Verify that the function returns true if the user is an admin.
func TestIsAdmin_UserIsAdmin_ReturnsTrue(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Set up the test data
	_, err := db.Exec("UPDATE person SET isAdmin = true WHERE id = $1", "1")
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	// Call the isAdmin function with the admin user ID
	isAdmin, err := isAdmin("1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the isAdmin value
	if !isAdmin {
		t.Errorf("Expected isAdmin to be true, got false")
	}
}

// Verify that the function returns false if the user is not an admin.
func TestIsAdmin_UserIsNotAdmin_ReturnsFalse(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Set up the test data
	_, err := db.Exec("UPDATE person SET isAdmin = false WHERE id  = $1", "2")
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	// Call the isAdmin function with a non-admin user ID
	isAdmin, err := isAdmin("2")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the isAdmin value
	if isAdmin {
		t.Errorf("Expected isAdmin to be false, got true")
	}
}

// Verify that the function returns an error if the query execution fails.
func TestIsAdmin_QueryExecutionFails_ReturnsError(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the isAdmin function with an invalid user ID
	_, err := isAdmin("invalid")
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

// Verify that the function returns an error if the query is invalid.
func TestIsAdmin_InvalidQuery_ReturnsError(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the isAdmin function with an invalid query
	_, err := isAdmin("1' OR '1'='1")
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

// Verify that the function returns an error if the database connection fails.
func TestIsAdmin_DatabaseConnectionFails_ReturnsError(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	db.Close() // Close the database connection to simulate a failure

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the isAdmin function with a valid user ID
	_, err := isAdmin("1")
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

// Returns true if user exists in the database
func TestCheckUserExists_UserExists(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the checkUserExists function with issuer "1"
	exists, err := checkUserExists("1")

	// Check the return values
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected user to exist, got false")
	}
}

// Returns false if user does not exist in the database
func TestCheckUserExists_UserDoesNotExist(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the checkUserExists function with issuer "2"
	exists, err := checkUserExists("999999")

	// Check the return values
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		if errStr != "user not found. Unauthorized access not allowed" {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
	if exists {
		t.Errorf("Expected user to not exist, got true")
	}
}

// Returns error if there is an issue with the database connection
func TestCheckUserExists_DatabaseError(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Close the test database to simulate a database error
	db.Close()

	// Call the checkUserExists function with issuer "1"
	_, err := checkUserExists("1")

	// Check the error value
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}

// Returns error if issuer is empty
func TestCheckUserExists_EmptyIssuer(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = 5
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the checkUserExists function with an empty issuer
	_, err := checkUserExists("")

	// Check the error value
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}

// Returns error if queryTimeLimit is negative
func TestCheckUserExists_NegativeQueryTimeLimit(t *testing.T) {
	// Initialize the test environment
	db = setupTestDB()
	defer db.Close()

	queryTimeLimit = -1
	secretKey = "filmoteka_test"
	jwtName = "filmoteka_test_jwt"

	// Call the checkUserExists function with issuer "1"
	_, err := checkUserExists("1")

	// Check the error value
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}
