package routes

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BukhryakovVladimir/vkTest/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func AddMovie(w http.ResponseWriter, r *http.Request) {
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

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to add movies", http.StatusUnauthorized)
		return
	}

	var movie model.Movie

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if len([]rune(movie.Name)) < 1 || len([]rune(movie.Name)) > 150 {
		http.Error(w, "Movie name must be between 1 and 150 characters long", http.StatusBadRequest)
		return
	}

	if len([]rune(movie.Description)) > 1000 {
		http.Error(w, "Movie description maximum length is 1000 characters", http.StatusBadRequest)
		return
	}

	if movie.Rating < 0 || movie.Rating > 10 {
		http.Error(w, "Movie rating must be between 0 and 10", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	var actorsID []int
	var actorID int

	addMissingActorsQuery := `
	INSERT INTO actor (firstName, lastName, sex, birthDate)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (firstName, lastName, birthDate) DO UPDATE
	SET firstName = EXCLUDED.firstName
	RETURNING id;
	`

	for _, actor := range movie.Actors {
		err = tx.QueryRowContext(ctx, addMissingActorsQuery, actor.FirstName,
			actor.LastName, actor.Sex, actor.BirthDate).Scan(&actorID)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("AddMovie Failed to rollback transaction: %v\n", rollbackErr)
			} else {
				log.Println("AddMovie transaction rollback")
			}

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Println("AddMovie ExecContext deadline exceeded while adding movie: ", err)
				http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
				return
			}
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if actorID != 0 {
			actorsID = append(actorsID, actorID)
		}
	}

	addMovieQuery := `
	INSERT INTO movie (name, description, date, rating)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	var movieID int

	err = tx.QueryRowContext(ctx, addMovieQuery, movie.Name, movie.Description, movie.Date, movie.Rating).Scan(&movieID)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("AddMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("AddMovie transaction rollback")
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("AddMovie ExecContext deadline exceeded while adding movie: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		}

		var errPQ *pq.Error
		if errors.As(err, &errPQ) {
			if errPQ.Code == "23505" {
				log.Println("Movie already exists: ", errPQ)
				http.Error(w, "Movie already exits", http.StatusBadRequest)
				return
			}
		}

		log.Println("Database error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return

	}

	addActorMovieRelQuery := `
	INSERT INTO ActorMovie (actor_id, movie_id)
	VALUES ($1, $2);
	`

	for _, actorID = range actorsID {
		_, err := tx.ExecContext(ctx, addActorMovieRelQuery, actorID, movieID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("AddMovie Failed to rollback transaction: %v\n", rollbackErr)
			} else {
				log.Println("AddMovie transaction rollback")
			}

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Println("AddMovie ExecContext deadline exceeded while adding movie: ", err)
				http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
				return
			}
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("AddMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("AddMovie transaction rollback")
		}

		log.Println("AddMovie error committing transaction")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode("Added a movie successfully")

}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
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

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to update movies", http.StatusUnauthorized)
		return
	}

	var movie model.Movie

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Println("UpdateMovie error reading request body: ", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if len([]rune(movie.Name)) > 150 {
		http.Error(w, "Movie name maximum length is 150 characters", http.StatusBadRequest)
		return
	}

	if len([]rune(movie.Description)) > 1000 {
		http.Error(w, "Movie description maximum length is 1000 characters", http.StatusBadRequest)
		return
	}

	if movie.Rating > 10 {
		http.Error(w, "Movie rating maximum value is 10", http.StatusBadRequest)
		return
	}

	updateMovieQuery := `
	UPDATE movie
	SET 
	name = COALESCE(NULLIF($1, ''), name),
	description = COALESCE(NULLIF($2, ''), description),
	rating = $3, 
	date = CASE WHEN $4::date = '0001-01-01' THEN date ELSE $4::date END 
	WHERE id = $5;
`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, updateMovieQuery, movie.Name, movie.Description,
		movie.Rating, movie.Date, movie.ID)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("UpdateMovie ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		}

		var errPQ *pq.Error
		if errors.As(err, &errPQ) {
			if errPQ.Code == "23505" {
				log.Println("UpdateMovie movie already exists: ", errPQ)
				http.Error(w, "Movie already exits", http.StatusBadRequest)
				return
			}
		}

		log.Println("Database error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return

	}

	resp, err := json.Marshal("Movie updated successfully")
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

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
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

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to delete movies", http.StatusUnauthorized)
		return
	}

	var movie model.Movie
	err = json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if movie.ID == 0 {
		http.Error(w,
			"id is not set, movie is deleted based on id. Please set id and make a request again",
			http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	deleteActorMovieQuery := `DELETE FROM actormovie WHERE movie_id = $1;`

	result, err := tx.ExecContext(ctx, deleteActorMovieQuery, movie.ID)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("DeleteMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("DeleteMovie transaction rollback")
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("DeleteMovie ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("DeleteMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("DeleteMovie transaction rollback")
		}

		log.Println("DeleteMovie error checking rows affected : ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("DeleteMovie ActorMovie table. ActorMovie relations don't exist, no rows affected")
	} else {
		log.Println("DeleteMovie Deleted movie with id = ", movie.ID, " from ActorMovie table")
	}

	deleteMovieQuery := `DELETE FROM movie WHERE id = $1;`

	result, err = tx.ExecContext(ctx, deleteMovieQuery, movie.ID)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("DeleteMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("DeleteMovie transaction rollback")
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("DeleteMovie ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("DeleteMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("DeleteMovie transaction rollback")
		}

		log.Println("DeleteMovie error checking rows affected : ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("DeleteMovie movie table. Movie doesn't exist, no rows affected")
		http.Error(w, "Movie doesn't exist. Nothing deleted", http.StatusBadRequest)
		return
	}

	log.Println("DeleteMovie Deleted movie with id = ", movie.ID, " from Movie table")

	resp, err := json.Marshal("Movie deleted successfully")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("DeleteMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("DeleteMovie transaction rollback")
		}

		log.Println("DeleteMovie error committing transaction")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

func AddActorToMovie(w http.ResponseWriter, r *http.Request) {
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

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to add actor to movie", http.StatusUnauthorized)
		return
	}

	var actorMovie model.ActorMovie

	err = json.NewDecoder(r.Body).Decode(&actorMovie)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if actorMovie.MovieID == 0 {
		http.Error(w, "Movie id cannot be empty", http.StatusBadRequest)
		return
	}

	if len([]rune(actorMovie.FirstName)) > 255 {
		http.Error(w, "Maximum firstName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actorMovie.LastName)) > 255 {
		http.Error(w, "Maximum lastName string length is 255 symbols", http.StatusBadRequest)
		return
	}

	if len([]rune(actorMovie.Sex)) > 10 {
		http.Error(w, "Maximum sex string length is 10 symbols", http.StatusBadRequest)
		return
	}

	if actorMovie.BirthDate.After(time.Now()) {
		http.Error(w, "Birth date cannot be in the future", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	var actorID int

	addMissingActorQuery := `
	INSERT INTO actor (firstName, lastName, sex, birthDate)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (firstName, lastName, birthDate) DO UPDATE
	SET firstName = EXCLUDED.firstName
	RETURNING id;
	`

	err = tx.QueryRowContext(ctx, addMissingActorQuery, actorMovie.FirstName,
		actorMovie.LastName, actorMovie.Sex, actorMovie.BirthDate).Scan(&actorID)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("AddActorToMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("AddActorToMovie transaction rollback")
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("AddActorToMovie ExecContext deadline exceeded while adding movie: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		}
		log.Println("Database error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	addActorMovieRelQuery := `
	INSERT INTO ActorMovie (actor_id, movie_id)
	VALUES ($1, $2);
	`

	_, err = tx.ExecContext(ctx, addActorMovieRelQuery, actorID, actorMovie.MovieID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("AddActorToMovieRel Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("AddActorToMovieRel transaction rollback")
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("AddActorToMovieRel ExecContext deadline exceeded while adding movie: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		}
		log.Println("Database error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("AddActorToMovie Failed to rollback transaction: %v\n", rollbackErr)
		} else {
			log.Println("AddActorToMovie transaction rollback")
		}

		log.Println("AddActorToMovie error committing transaction")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode("Added an actor to movie successfully")
}

func DeleteActorFromMovie(w http.ResponseWriter, r *http.Request) {
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

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to delete actor from movie", http.StatusUnauthorized)
		return
	}

	var actorMovie model.ID

	err = json.NewDecoder(r.Body).Decode(&actorMovie)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if actorMovie.MovieID == 0 {
		http.Error(w, "Movie id cannot be empty", http.StatusBadRequest)
		return
	}

	if actorMovie.ActorID == 0 {
		http.Error(w, "Actor id cannot be empty", http.StatusBadRequest)
		return
	}

	deleteActorFromMovieQuery := `
	DELETE FROM ActorMovie
	WHERE actor_id = $1
	AND movie_id = $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, deleteActorFromMovieQuery, actorMovie.ActorID, actorMovie.MovieID)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("DeleteActorFromMovie ExecContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal("Actor deleted from movie successfully")
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

func GetMoviesWithID(w http.ResponseWriter, r *http.Request) {
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

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isAdmin, err := isAdmin(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking administrator privileges", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have administrator privileges to get movies with id", http.StatusUnauthorized)
		return
	}

	var movie model.SearchMovie

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if len([]rune(movie.Name)) > 150 {
		http.Error(w, "Movie name maximum length is 150 characters", http.StatusBadRequest)
		return
	}

	if len([]rune(movie.Description)) > 1000 {
		http.Error(w, "Movie description maximum length is 1000 characters", http.StatusBadRequest)
		return
	}

	if movie.Rating > 10 {
		http.Error(w, "Movie rating maximum value is 10", http.StatusBadRequest)
		return
	}

	getMoviesQuery := `
	SELECT DISTINCT m.* FROM movie m 
	JOIN ActorMovie am on m.id = am.movie_id
	JOIN actor a ON am.actor_id = a.id
	WHERE ($1 <> '' AND m.name LIKE '%' || $1 || '%')
    OR ($2 <> '' AND m.description LIKE '%' || $2 || '%')
	OR ($3 <> '' AND m.date = $3::date)
	OR (m.rating = $4)
	OR (($5 <> '' AND a.firstName LIKE '%' || $5 || '%')
    OR ($6 <> '' AND a.lastName LIKE '%' || $6 || '%'));
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, getMoviesQuery, movie.Name, movie.Description,
		movie.Date, movie.Rating, movie.ActorFirstName, movie.ActorLastName)
	defer rows.Close()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("GetMoviesWithID QueryContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	var moviesWithID []model.SearchMovie
	var movieWithID model.SearchMovie
	for rows.Next() {
		if err := rows.Scan(&movieWithID.ID, &movieWithID.Name,
			&movieWithID.Description, &movieWithID.Date, &movieWithID.Rating); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		moviesWithID = append(moviesWithID, movieWithID)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(moviesWithID)
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
