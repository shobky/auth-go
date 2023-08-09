# Authentication with go-gin
Example applications built with golang and gin web freamework. Using postgresql as a DB and GORM for an orm.

See the [gofiber branch](https://github.com/shobky/auth-go/blob/gofiber/README.md) for the more secure authentication built with fiber and uses refresh tokens.

## Setup TODO
* get the packages
  ```b
  go get .
  ```
* add the .env file and populate it with the variables from .env.example.txt
* have a database running on port 5433, You can edit this to use you custom postgres db and edit values in .env file
  ```g
	host := "localhost"
	port := os.Getenv("DB_PORT")
	dbName := "postgres"
	dbUser := "root"
	password := string(os.Getenv("DB_PASSWORD"))
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
  ```
  - For database I used postgres15-alpine docker image and ran it on port 5433
  - If you want to change the Database provider remember to import the right gorm/driver/*
    ```g
    package initializers

    import (
      "fmt"
      "log"
      "os"
    
      // change this to whatever
      "gorm.io/driver/postgres"
      "gorm.io/gorm"
    )
    ```
