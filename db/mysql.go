package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)
var dbPool  = map[string]*sql.DB{}


const ErrorDuplicateEntry = 1062
const MasterName = "master"
const Slave1Name = "slave1"
const Slave2Name = "slave2"

func OpenDB(name string) *sql.DB {

	if _, ok := dbPool[name]; ok {
		return dbPool[name]
	}

	nameForEnv := strings.ToUpper(name)
	fmt.Println("DB_"+ nameForEnv +"_HOST")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_"+ nameForEnv +"_USERNAME"),
		os.Getenv("DB_"+ nameForEnv +"_PASSWORD"),
		os.Getenv("DB_"+ nameForEnv +"_HOST"),
		os.Getenv("DB_"+ nameForEnv +"_PORT"),
		os.Getenv("DB_"+ nameForEnv +"_DATABASE"))
	var err error
	dbPool[name], err = sql.Open("mysql", dataSource)
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

	dbPool[name].SetMaxIdleConns(maxIdleConns)
	dbPool[name].SetMaxOpenConns(maxOpenConns)

	lifeTime := time.Second * time.Duration(maxLifeConns)
	dbPool[name].SetConnMaxLifetime(lifeTime)

	err = dbPool[name].Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return dbPool[name]
}

func Close(name string) {
	if _, ok := dbPool[name]; ok && dbPool[name] != nil {
		err := dbPool[name].Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}