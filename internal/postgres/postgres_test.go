package postgres

import (
	"os"
	"testing"
)

// Returns a *sql.DB and nil error when all environment variables are set correctly
func TestDial_ValidEnvironmentVariables(t *testing.T) {
	// Set up the test environment
	os.Setenv("DB_HOST", "filmoteka-postgres-test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "filmoteka")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	// Call the Dial function
	db, err := Dial()

	// Check the returned values
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if db == nil {
		t.Error("Expected non-nil *sql.DB, got nil")
	}
}

// Sets the maximum number of open connections, maximum number of idle connections, and connection max lifetime
func TestDial_SetConnectionSettings(t *testing.T) {
	// Set up the test environment
	os.Setenv("DB_HOST", "filmoteka-postgres-test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "filmoteka")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	// Call the Dial function
	db, err := Dial()
	if err != nil {
		t.Errorf("Error starting db")
	}

	// Check the connection settings
	if db != nil {
		if db.Stats().MaxOpenConnections != 10 {
			t.Errorf("Expected MaxOpenConnections to be 10, got %d", db.Stats().MaxOpenConnections)
		}
	} else {
		t.Error("Expected non-nil *sql.DB, got nil")
	}

	// Close the database connection
	db.Close()
}

// Successfully pings the database
func TestDial_PingDatabase(t *testing.T) {
	// Set up the test environment
	os.Setenv("DB_HOST", "filmoteka-postgres-test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "filmoteka")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	// Call the Dial function
	db, err := Dial()

	// Check the database connection
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if db == nil {
		t.Error("Expected non-nil *sql.DB, got nil")
	} else {
		err = db.Ping()
		if err != nil {
			t.Errorf("Expected nil error from Ping, got %v", err)
		}
	}

	// Close the database connection
	db.Close()
}

// Returns an error when DB_HOST environment variable is empty
func TestDial_EmptyDBHost(t *testing.T) {
	// Set up the test environment
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	// Call the Dial function
	db, err := Dial()

	// Check the returned error
	if err == nil {
		t.Error("Expected non-nil error, got nil")
	}
	if db != nil {
		t.Error("Expected nil *sql.DB, got non-nil")
	}
}

// Returns an error when DB_PORT environment variable is empty
func TestDial_EmptyDBPort(t *testing.T) {
	// Set up the test environment
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	// Call the Dial function
	db, err := Dial()

	// Check the returned error
	if err == nil {
		t.Error("Expected non-nil error, got nil")
	}
	if db != nil {
		t.Error("Expected nil *sql.DB, got non-nil")
	}
}

// Returns an error when DB_USER environment variable is empty
func TestDial_EmptyDBUser(t *testing.T) {
	// Set up the test environment
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	// Call the Dial function
	db, err := Dial()

	// Check the returned error
	if err == nil {
		t.Error("Expected non-nil error, got nil")
	}
	if db != nil {
		t.Error("Expected nil *sql.DB, got non-nil")
	}
}

func TestDial_EmptyDBName(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	_, err := Dial()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "environment variable DB_NAME is empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestDial_EmptyDBPassword(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	_, err := Dial()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "environment variable DB_PASSWORD is empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestDial_EmptyMaxOpenConns(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("MAX_OPEN_CONNS", "")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	_, err := Dial()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "environment variable MAX_OPEN_CONNS is empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestDial_MaxIdleConnsEmpty(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "")
	os.Setenv("CONN_MAX_LIFETIME", "30")

	_, err := Dial()

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestDial_ConnMaxLifetimeEmpty(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("MAX_OPEN_CONNS", "10")
	os.Setenv("MAX_IDLE_CONNS", "5")
	os.Setenv("CONN_MAX_LIFETIME", "")

	_, err := Dial()

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}
