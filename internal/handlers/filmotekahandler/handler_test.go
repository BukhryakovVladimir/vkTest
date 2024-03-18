package filmotekahandler

import (
	"github.com/BukhryakovVladimir/vkTest/internal/routes"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// All routes are correctly registered with their respective HTTP methods
func TestSetupRoutes_RegisterRoutesWithCorrectHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	// Test POST /api/signup
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", nil)
	route, _ := mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/signup not registered")
	}

	// Test POST /api/login
	req, _ = http.NewRequest(http.MethodPost, "/api/login", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/login not registered")
	}

	// Test POST /api/add-actor
	req, _ = http.NewRequest(http.MethodPost, "/api/add-actor", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/add-actor not registered")
	}

	// Test PUT /api/update-actor
	req, _ = http.NewRequest(http.MethodPut, "/api/update-actor", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route PUT /api/update-actor not registered")
	}

	// Test DELETE /api/delete-actor
	req, _ = http.NewRequest(http.MethodDelete, "/api/delete-actor", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route DELETE /api/delete-actor not registered")
	}

	// Test POST /api/get-actors-with-id
	req, _ = http.NewRequest(http.MethodPost, "/api/get-actors-with-id", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/get-actors-with-id not registered")
	}
}

// All routes are correctly registered with their respective endpoints
func TestSetupRoutes_RegisterRoutesWithCorrectEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	// Test POST /api/signup
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", nil)
	route, _ := mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/signup not registered")
	}

	// Test POST /api/login
	req, _ = http.NewRequest(http.MethodPost, "/api/login", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/login not registered")
	}

	// Test POST /api/add-actor
	req, _ = http.NewRequest(http.MethodPost, "/api/add-actor", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/add-actor not registered")
	}

	// Test PUT /api/update-actor
	req, _ = http.NewRequest(http.MethodPut, "/api/update-actor", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route PUT /api/update-actor not registered")
	}

	// Test DELETE /api/delete-actor
	req, _ = http.NewRequest(http.MethodDelete, "/api/delete-actor", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route DELETE /api/delete-actor not registered")
	}

	// Test POST /api/get-actors-with-id
	req, _ = http.NewRequest(http.MethodPost, "/api/get-actors-with-id", nil)
	route, _ = mux.Handler(req)
	if route == nil {
		t.Errorf("Route POST /api/get-actors-with-id not registered")
	}
}

func TestSetupRoutes_RegisterRoutesWithCorrectHandlers(t *testing.T) {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	routes := map[string]http.HandlerFunc{
		"POST /api/signup":                    routes.SignupPerson,
		"POST /api/login":                     routes.LoginPerson,
		"POST /api/add-actor":                 routes.AddActor,
		"PUT /api/update-actor":               routes.UpdateActor,
		"DELETE /api/delete-actor":            routes.DeleteActor,
		"POST /api/get-actors-with-id":        routes.GetActorsWithID,
		"GET /api/actors":                     routes.GetActors,
		"POST /api/add-movie":                 routes.AddMovie,
		"PUT /api/update-movie":               routes.UpdateMovie,
		"DELETE /api/delete-movie":            routes.DeleteMovie,
		"POST /api/get-movies-with-id":        routes.GetMoviesWithID,
		"POST /api/add-actor-to-movie":        routes.AddActorToMovie,
		"DELETE /api/delete-actor-from-movie": routes.DeleteActorFromMovie,
		"GET /api/movies":                     routes.GetMoviesOrdered,
		"POST /api/search-movie":              routes.SearchMovie,
	}

	for route, handler := range routes {
		s := strings.Split(route, " ")
		method, path := s[0], s[1]
		req, _ := http.NewRequest(method, path, nil)
		route, _ := mux.Handler(req)
		if route == nil {
			t.Errorf("Route %s not registered", route)
		} else if reflect.ValueOf(route).Pointer() != reflect.ValueOf(handler).Pointer() {
			t.Errorf("Route %s does not have the correct handler", route)
		}
	}
}
