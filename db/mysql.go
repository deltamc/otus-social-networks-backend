package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)
var dbPull *sql.DB

const ERROR_DUPLICATE_ENTRY  = 1062

func OpenDB() *sql.DB {

	//if dbPull != nil {
	//	return dbPull
	//}

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"))
	var err error
	dbPull, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalln(err)
	}


	maxOpenConns, err := strconv.Atoi(os.Getenv("SQL_MAX_OPEN_CONNECT"))
	if err != nil {
		maxOpenConns = 5
		fmt.Println(fmt.Sprintf("SQL_MAX_OPEN_CONNECT, %s", err.Error()))
	}
	maxIdleConns, err := strconv.Atoi(os.Getenv("SQL_MAX_IDLE_CONNECT"))
	if err != nil {
		maxIdleConns = 5
		fmt.Println(fmt.Sprintf("SQL_MAX_IDLE_CONNECT, %s", err.Error()))
	}
	maxLifeConns, err := strconv.Atoi(os.Getenv("SQL_MAX_LIFE_CONNECT"))
	if err != nil {
		maxLifeConns = 3600
		fmt.Println(fmt.Sprintf("SQL_MAX_LIFE_CONNECT, %s", err.Error()))
	}

	dbPull.SetMaxIdleConns(maxIdleConns)
	dbPull.SetMaxOpenConns(maxOpenConns)

	lifeTime := time.Second * time.Duration(maxLifeConns)
	dbPull.SetConnMaxLifetime(lifeTime)

	err = dbPull.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return dbPull
}

func Close() {
	if dbPull != nil {
		err := dbPull.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}