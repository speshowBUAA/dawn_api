package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "dawn_api/log"
	// "time"
	"fmt"
    "go.uber.org/zap"
)

//数据库地址
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "postgres"
)

var DawnApiDataClient *sql.DB

func checkErr(err error) {
	if err != nil {
        log.Error("Error", zap.Any("error", err))
	}
}

type User struct{
	id               int
	name             string
	phone            string
	email            string
	create_time      string
	last_login_time  string
	avatar_url       string
}

func NewDawnApiDataClient() {
	psqlInfo :=  fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d",
	host,user,dbname,password,port)
    db, err := sql.Open("postgres", psqlInfo)
    checkErr(err)

    err = db.Ping()
    checkErr(err)
    DawnApiDataClient = db
}


func Query(key string, value string) User{
	var user User
    sqlStatement := fmt.Sprintf("SELECT * FROM public.users WHERE %s = '%s' ORDER BY id DESC LIMIT 1" , key, value)
	fmt.Println(sqlStatement)
	row := DawnApiDataClient.QueryRow(sqlStatement)
    row.Scan(&user.id, &user.name, &user.phone, &user.email, &user.create_time, &user.last_login_time, &user.avatar_url)
    return user
}

func CloseConnect(){
	DawnApiDataClient.Close()
}