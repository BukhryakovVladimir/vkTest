package routes

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BukhryakovVladimir/vkTest/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

func AddActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie(jwtName)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	token, err := jwtCheck(cookie)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to add actors", http.StatusUnauthorized)
		return
	}

	var actor model.Actor
	err = json.NewDecoder(r.Body).Decode(&actor)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if len([]rune(actor.FirstName)) > 255 {
		http.Error(w, "Maximum firstName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actor.LastName)) > 255 {
		http.Error(w, "Maximum lastName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actor.Sex)) > 10 {
		http.Error(w, "Maximum sex string length is 10 symbols", http.StatusBadRequest)
		return
	}

	if actor.BirthDate.After(time.Now()) {
		http.Error(w, "Birth date cannot be in the future", http.StatusBadRequest)
		return
	}

	addActorQuery := `INSERT INTO actor (firstName, lastName, sex, birthDate) 
						VALUES ($1::text, $2::text, $3::text, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, addActorQuery, actor.FirstName, actor.LastName, actor.Sex, actor.BirthDate)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("AddActor ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal("Actor added successfully")
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

func UpdateActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie(jwtName)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	token, err := jwtCheck(cookie)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to add actors", http.StatusUnauthorized)
		return
	}

	var actor model.Actor
	err = json.NewDecoder(r.Body).Decode(&actor)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if len([]rune(actor.FirstName)) > 255 {
		http.Error(w, "Maximum firstName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actor.LastName)) > 255 {
		http.Error(w, "Maximum lastName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actor.Sex)) > 10 {
		http.Error(w, "Maximum sex string length is 10 symbols", http.StatusBadRequest)
		return
	}

	if actor.BirthDate.After(time.Now()) {
		http.Error(w, "Birth date cannot be in the future", http.StatusBadRequest)
		return
	}

	updateActorQuery := `UPDATE actor
							SET 
								firstName = COALESCE(NULLIF($1, ''), firstName),
								lastName = COALESCE(NULLIF($2, ''), lastName),
								sex = COALESCE(NULLIF($3, ''), sex),
								birthDate = CASE WHEN $4::date = '0001-01-01' THEN birthDate ELSE $4::date END
						WHERE id = $5;
						    `

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, updateActorQuery, actor.FirstName, actor.LastName, actor.Sex, actor.BirthDate, actor.ID)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("UpdateActor ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal("Actor updated successfully")
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

func DeleteActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie(jwtName)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	token, err := jwtCheck(cookie)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to add actors", http.StatusUnauthorized)
		return
	}

	var actor model.Actor
	err = json.NewDecoder(r.Body).Decode(&actor)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if actor.ID == 0 {
		http.Error(w,
			"id is not set, actor is deleted based on id. Please set id and make a request again",
			http.StatusBadRequest)
		return
	}

	deleteActorQuery := `DELETE FROM actor WHERE id = $1;`

	// queryTimeLimit is multiplied by 2, because ctx is used by two queries
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*queryTimeLimit)*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, deleteActorQuery, actor.ID)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("DeleteActor ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	log.Println("Deleted actor with id = ", actor.ID, " from Actor table")

	deleteActorMovieQuery := `DELETE FROM actormovie WHERE actor_id = $1;`

	_, err = db.ExecContext(ctx, deleteActorMovieQuery, actor.ID)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("DeleteActor ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	log.Println("Deleted actor with id = ", actor.ID, " from ActorMovie table")

	resp, err := json.Marshal("Actor deleted successfully")
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

func GetActorsWithID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie(jwtName)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	token, err := jwtCheck(cookie)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to add actors", http.StatusUnauthorized)
		return
	}

	var actor model.Actor
	err = json.NewDecoder(r.Body).Decode(&actor)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if len([]rune(actor.FirstName)) > 255 {
		http.Error(w, "Maximum firstName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actor.LastName)) > 255 {
		http.Error(w, "Maximum lastName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actor.Sex)) > 10 {
		http.Error(w, "Maximum sex string length is 10 symbols", http.StatusBadRequest)
		return
	}

	if actor.BirthDate.After(time.Now()) {
		http.Error(w, "Birth date cannot be in the future", http.StatusBadRequest)
		return
	}

	getActorsQuery := `
    SELECT * FROM actor
    WHERE ($1 <> '' AND firstName LIKE '%' || $1 || '%')
    OR ($2 <> '' AND lastName LIKE '%' || $2 || '%')
    OR (sex = $3)
    OR ($4 <> '' AND birthDate = $4::date)
    ORDER BY id;
`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, getActorsQuery, actor.FirstName, actor.LastName, actor.Sex, actor.BirthDate)
	defer rows.Close()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("GetActorsWithID QueryContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	var actorsWithID []model.Actor
	var actorWithID model.Actor
	for rows.Next() {
		if err := rows.Scan(&actorWithID.ID, &actorWithID.FirstName,
			&actorWithID.LastName, &actorWithID.Sex, &actorWithID.BirthDate); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		actorsWithID = append(actorsWithID, actorWithID)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(actorsWithID)
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
