.PHONY: test

test:
	-docker rm -f filmoteka_golang_test
	-docker rm -f filmoteka_db_postgres_test
	-docker volume rm vktest_database_postgres_test
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

.PHONY: run

run:
	docker-compose up --build
