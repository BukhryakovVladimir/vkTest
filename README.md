# vkTest
Тестовое задание на стажировку VK

Необходимо разработать бэкенд приложения “Фильмотека”, который предоставляет REST API для управления базой данных фильмов.

Приложение должно поддерживать следующие функции:
* добавление информации об актёре (имя, пол, дата рождения),
* изменение информации об актёре.

Возможно изменить любую информацию об актёре, как частично, так и полностью:
* удаление информации об актёре,
* добавление информации о фильме.

При добавлении фильма указываются его название (не менее 1 и не более 150 символов), описание (не более 1000 символов), дата выпуска, рейтинг (от 0 до 10) и список актёров:
* изменение информации о фильме.

Возможно изменить любую информацию о фильме, как частично, так и полностью:
* удаление информации о фильме,
* получение списка фильмов с возможностью сортировки по названию, по рейтингу, по дате выпуска. По умолчанию используется сортировка по рейтингу (по убыванию),
* поиск фильма по фрагменту названия, по фрагменту имени актёра,
* получение списка актёров, для каждого актёра выдаётся также список фильмов с его участием,
* API должен быть закрыт авторизацией,
* поддерживаются две роли пользователей - обычный пользователь и администратор. Обычный пользователь имеет доступ только на получение данных и поиск, администратор - на все действия. Для упрощения можно считать, что соответствие пользователей и ролей задаётся вручную (например, напрямую через БД).

Требования к реализации:
* язык реализации - go,
* для хранения данных используется реляционная СУБД (предпочтительно - PostgreSQL),
* предоставлена спецификация на API (в формате Swagger 2.0 или OpenAPI 3.0).

Бонус: используется подход api-first (генерация кода из спецификации) или code-first (генерация спецификации из кода).
* Для реализации http сервера разрешается использовать только стандартную библиотеку http (без фреймворков),
* логирование - в лог должна попадать базовая информация об обрабатываемых запросах, ошибки,
* код приложения покрыт юнит-тестами не менее чем на 70%,
* Dockerfile для сборки образа,
* docker-compose файл для запуска окружения с работающим приложением и СУБД.

# Тестирование
Чтобы протестировать введите команду в терминале
```
make test
```
![image](https://github.com/BukhryakovVladimir/vkTest/assets/43881945/e8f22d2f-2df6-4ddb-a405-1c4cb3250b86)


# Запуск
Чтобы запустить введите команду в терминале
```
make run
```
![image](https://github.com/BukhryakovVladimir/vkTest/assets/43881945/5bc6f36a-2301-47be-9bb1-7f855b438684)


# Примеры
**localhost:3000/api/signup**

тело запроса:
```json
{
	"username": "kidala",
	"password": "*****",
	"firstName": "Vladimir",
	"lastName": "Bukhryakov",
	"sex": "Male",
	"birthdate": "2002-04-18T00:00:00Z"
}
```
тело ответа:
```
"Signup successful"
```

**localhost:3000/api/login**

тело запроса:
```json
{
	"username": "kidala",
	"password": "*****"
}
```
тело ответа:
```
"Successfully logged in"
```

**localhost:3000/api/add-actor**

тело запроса:
```json
{
	"firstName": "Cillian", 
	"lastName": "Murphy", 
	"sex": "Male",
	"birthDate": "1976-05-25T00:00:00Z"
}
```
тело ответа:
```
"Actor added successfully"
```

**localhost:3000/api/update-actor**

тело запроса:
```json
{
    "id": 4, 
	"firstName": "Robert", 
	"lastName": "De Niro", 
	"sex": "Male",
	"birthDate": "1943-08-17T00:00:00Z"
}
```
тело ответа:
```
"Actor updated successfully"
```

**localhost:3000/api/delete-actor**

тело запроса:
```json
{
    "id": 44
}
```
вывод:
```
"Actor deleted successfully"
```

**localhost:3000/api/get-actors-with-id**

тело запроса:
```json
{
	"firstName": "Al" 
}
```
тело ответа:
```json
[
    {
        "id": 1,
        "firstName": "Al",
        "lastName": "Pacino",
        "sex": "Male",
        "birthDate": "1940-04-25T00:00:00Z"
    }
]
```

**localhost:3000/api/add-movie**

тело запроса:
```json
{
    "name": "somethin",
    "description": "once upon a time in...",
    "date": "2010-01-01T00:00:00Z",
    "rating": 8,
    "actors":   [
        {
        "firstName": "Al", 
        "lastName": "Pacino", 
        "sex": "Male",
        "birthDate": "1940-04-25T00:00:00Z"
        },
        {
	    "firstName": "Robert", 
	    "lastName": "De Niro", 
	    "sex": "Male",
	    "birthDate": "1943-08-17T00:00:00Z"
        },
        {
	    "firstName": "Robert", 
	    "lastName": "Pattinson", 
	    "sex": "Male",
	    "birthDate": "1986-05-13T00:00:00Z"
        }
    ]
}
```
тело ответа:
```
"Added a movie successfully"
```

**localhost:3000/api/get-movies-with-id**

тело запроса:
```json
{
    "actorFirstName": "Robert"
}
```
тело ответа:
```json
[
    {
        "id": 10,
        "name": "somethin",
        "description": "once upon a time in...",
        "date": "2010-01-01T00:00:00Z",
        "rating": 8,
        "actorFirstName": "",
        "actorLastName": ""
    },
    {
        "id": 8,
        "name": "joska",
        "description": "az",
        "date": "2011-01-01T00:00:00Z",
        "rating": 7,
        "actorFirstName": "",
        "actorLastName": ""
    },
    {
        "id": 9,
        "name": "gaaaa",
        "description": "there was once a gzxczxczxczxczxczxcoy",
        "date": "2010-01-01T00:00:00Z",
        "rating": 8,
        "actorFirstName": "",
        "actorLastName": ""
    }
]
```

**localhost:3000/api/update-movie**

тело запроса:
```json
{
    "id": 9,
    "name": "who",
    "description": "123123123",
    "date": "2011-01-01T00:00:00Z",
    "rating": 7
}
```
тело ответа:
```
"Movie updated successfully"
```

**localhost:3000/api/add-actor-to-movie**

тело запроса:
```json
{
  "movieID": 8,
	"firstName": "Solo", 
	"lastName": "Pacino", 
	"sex": "Male",
	"birthDate": "1990-04-25T00:00:00Z"
}
```
тело ответа:
```
"Added an actor to movie successfully"
```

**localhost:3000/api/delete-actor-from-movie**

тело запроса:
```json
{
    "movieID": 10,
    "actorID": 4
}
```
тело ответа:
```
"Actor deleted from movie successfully"
```

**localhost:3000/api/delete-movie**

тело запроса:
```json
{
    "id": 10
}
```
тело ответа:
```
"Movie deleted successfully"
```

**localhost:3000/api/movies?order=&by=date**

тело ответа:
```json
[
    {
        "id": 0,
        "name": "joska",
        "description": "az",
        "date": "2011-01-01T00:00:00Z",
        "rating": 7,
        "actors": [
            {
                "id": 0,
                "firstName": "Al",
                "lastName": "Pacino",
                "sex": "Male",
                "birthDate": "1940-04-25T00:00:00Z"
            },
            {
                "id": 0,
                "firstName": "Robert",
                "lastName": "De Niro",
                "sex": "Male",
                "birthDate": "1943-08-17T00:00:00Z"
            }
        ]
    },
    {
        "id": 0,
        "name": "who",
        "description": "123123123",
        "date": "2011-01-01T00:00:00Z",
        "rating": 7,
        "actors": [
            {
                "id": 0,
                "firstName": "Al",
                "lastName": "Pacino",
                "sex": "Male",
                "birthDate": "1940-04-25T00:00:00Z"
            },
            {
                "id": 0,
                "firstName": "Robert",
                "lastName": "De Niro",
                "sex": "Male",
                "birthDate": "1943-08-17T00:00:00Z"
            },
            {
                "id": 0,
                "firstName": "Robert",
                "lastName": "Pattinson",
                "sex": "Male",
                "birthDate": "1986-05-13T00:00:00Z"
            }
        ]
    }
]
```


**localhost:3000/api/search-movie**

тело запроса:
```json
{
    "name": "o"
}
```
тело ответа:
```json
[
    {
        "id": 0,
        "name": "joska",
        "description": "az",
        "date": "2011-01-01T00:00:00Z",
        "rating": 7,
        "actors": [
            {
                "id": 0,
                "firstName": "Al",
                "lastName": "Pacino",
                "sex": "Male",
                "birthDate": "1940-04-25T00:00:00Z"
            },
            {
                "id": 0,
                "firstName": "Robert",
                "lastName": "De Niro",
                "sex": "Male",
                "birthDate": "1943-08-17T00:00:00Z"
            }
        ]
    },
    {
        "id": 0,
        "name": "who",
        "description": "123123123",
        "date": "2011-01-01T00:00:00Z",
        "rating": 7,
        "actors": [
            {
                "id": 0,
                "firstName": "Al",
                "lastName": "Pacino",
                "sex": "Male",
                "birthDate": "1940-04-25T00:00:00Z"
            },
            {
                "id": 0,
                "firstName": "Robert",
                "lastName": "De Niro",
                "sex": "Male",
                "birthDate": "1943-08-17T00:00:00Z"
            },
            {
                "id": 0,
                "firstName": "Robert",
                "lastName": "Pattinson",
                "sex": "Male",
                "birthDate": "1986-05-13T00:00:00Z"
            }
        ]
    }
]
```

**localhost:3000/api/actors**

тело ответа:
```json
[
    {
        "id": 0,
        "firstName": "Al",
        "lastName": "Pacino",
        "sex": "Male",
        "birthDate": "1940-04-25T00:00:00Z",
        "movies": [
            {
                "id": 0,
                "name": "joska",
                "description": "az",
                "date": "2011-01-01T00:00:00Z",
                "rating": 7
            },
            {
                "id": 0,
                "name": "who",
                "description": "123123123",
                "date": "2011-01-01T00:00:00Z",
                "rating": 7
            }
        ]
    },
    {
        "id": 0,
        "firstName": "Robert",
        "lastName": "De Niro",
        "sex": "Male",
        "birthDate": "1943-08-17T00:00:00Z",
        "movies": [
            {
                "id": 0,
                "name": "joska",
                "description": "az",
                "date": "2011-01-01T00:00:00Z",
                "rating": 7
            },
            {
                "id": 0,
                "name": "who",
                "description": "123123123",
                "date": "2011-01-01T00:00:00Z",
                "rating": 7
            }
        ]
    },
    {
        "id": 0,
        "firstName": "Robert",
        "lastName": "Pattinson",
        "sex": "Male",
        "birthDate": "1986-05-13T00:00:00Z",
        "movies": [
            {
                "id": 0,
                "name": "who",
                "description": "123123123",
                "date": "2011-01-01T00:00:00Z",
                "rating": 7
            }
        ]
    },
    {
        "id": 0,
        "firstName": "Solo",
        "lastName": "Pacino",
        "sex": "Male",
        "birthDate": "1990-04-25T00:00:00Z",
        "movies": [
            {
                "id": 0,
                "name": "joska",
                "description": "az",
                "date": "2011-01-01T00:00:00Z",
                "rating": 7
            }
        ]
    }
]
```
