DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'filmoteka') THEN
        CREATE DATABASE filmoteka;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS Person (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    firstName VARCHAR(255),
    lastName VARCHAR(255),
    isAdmin BOOL NOT NULL,
    sex VARCHAR(10),
    birthDate DATE
);

CREATE TABLE IF NOT EXISTS Actor (
    id SERIAL PRIMARY KEY,
    firstName VARCHAR(255) NOT NULL,
    lastName VARCHAR(255) NOT NULL,
    sex VARCHAR(10) NOT NULL,
    birthDate DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS Movie (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL, CHECK ( name <> '' ),
    description VARCHAR(1000) NOT NULL,
    date DATE NOT NULL,
    rating SMALLINT NOT NULL, CHECK ( rating BETWEEN 0 AND 10 )
);

CREATE TABLE IF NOT EXISTS ActorMovie (
    actor_id INT,
    movie_id INT,
    FOREIGN KEY (actor_id) REFERENCES Actor(id),
    FOREIGN KEY (movie_id) REFERENCES Movie(id),
    PRIMARY KEY (actor_id, movie_id)
);

