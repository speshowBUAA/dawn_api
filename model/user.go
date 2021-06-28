package model

import (
	"time"
)

type User struct{
	Id               int		`json:"id"`
	Name             string		`json:"username"`
	Email            string		`json:"useremail"`
    Phone            string		`json:"usermobile"`
    Password         string		`json:"password"`
    Permission       bool		`json:"permission"`
    Last_login_time  time.Time	`json:"last_login_time"`
	Create_time      time.Time	`json:"create_time"`
	Avatar_url       string		`json:"avatar_url"`
}

// LoginReq 登录请求参数类
type LoginReq struct {
	UserName  string `json:"userName"`
	PassWord  string `json:"password"`
	AutoLogin bool   `json:"autoLogin"`
	Type_     string `json:"type"`
}

//搜索参数
type FilterParam struct {
	UserName  		string
	UserEmail 		string
	Permission		string
	StartCreateDate string
	EndCreateDate 	string
	StartLoginTime	string
	EndLoginTime	string
}