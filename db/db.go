package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "dawn_api/log"
    "dawn_api/model"
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

func NewDawnApiDataClient() {
	psqlInfo :=  fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d",
	host,user,dbname,password,port)
    db, err := sql.Open("postgres", psqlInfo)
    checkErr(err)

    err = db.Ping()
    checkErr(err)
    DawnApiDataClient = db
}


func Query(key string, value string) []model.User{
	var users []model.User
    sqlStatement := fmt.Sprintf("SELECT * FROM public.dawn_users WHERE %s = '%s'" , key, value)

    rows, err := DawnApiDataClient.Query(sqlStatement)
    checkErr(err)
    
    for rows.Next() {
        var user model.User
        rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Permission, &user.Last_login_time, &user.Create_time, &user.Avatar_url)
        users = append(users, user)
    }
    
    return users
}

func FilterQuery(param model.FilterParam) []model.User{
    var users []model.User
    sqlStatement := fmt.Sprintf(
        "SELECT * FROM public.dawn_users WHERE username like '%s' and useremail like '%s' and permission in (%s) and create_time > '%s' and create_time < '%s' and last_login_time > '%s' and last_login_time < '%s'", 
        param.UserName, param.UserEmail, param.Permission, param.StartCreateDate,
        param.EndCreateDate, param.StartLoginTime, param.EndLoginTime)
    rows, err := DawnApiDataClient.Query(sqlStatement)
    checkErr(err)
    
    for rows.Next() {
        var user model.User
        rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Permission, &user.Last_login_time, &user.Create_time, &user.Avatar_url)
        users = append(users, user)
    }
    
    return users
}

func UpdateUserInfo(user model.User) error {
    sqlStatement := fmt.Sprintf("UPDATE public.dawn_users SET useremail=$1,usermobile=$2,password=$3,permission=$4,last_login_time=NOW(),avatar_url=$5 WHERE id = $6")
    stmt, err := DawnApiDataClient.Prepare(sqlStatement)
    checkErr(err)

    _, err = stmt.Exec(user.Email, user.Phone, user.Password, user.Permission, user.Avatar_url, user.Id)
    checkErr(err)
    return err
}

func AddUserInfo(user model.User) error {
    sqlStatement := fmt.Sprintf("INSERT INTO public.dawn_users (username,useremail,usermobile,password,create_time,avatar_url) VALUES($1,$2,$3,$4,NOW(),$5)")
    stmt, err := DawnApiDataClient.Prepare(sqlStatement)
    checkErr(err)

    _, err = stmt.Exec(user.Name, user.Email, user.Phone, user.Password, user.Avatar_url)
    checkErr(err)
    return err
}

func CloseConnect(){
	DawnApiDataClient.Close()
}