package api

import (
	"fmt"
	"dawn_api/db"
	"dawn_api/model"
)

// Register 插入用户，先检查是否存在用户，如果没有则存入
func Register(name string, pwd string) error {
	if CheckUser(name) {
		return fmt.Errorf("用户已存在！")
	}
	addUser := model.User{Name:name, Password:pwd}
	err := db.AddUserInfo(addUser)
	return err
}

// CheckUser 检查用户是否存在
func CheckUser(userName string) bool {
	resultUser := db.Query("username", userName)[0]
	result := false
	if (resultUser.Id) != 0 {
		result = true
	}
	return result
}

// LoginCheck 登录验证
func LoginCheck(loginReq model.LoginReq) (bool, model.User) {
	resultBool := false
	resultUser := db.Query("username", loginReq.UserName)[0]
	if (loginReq.PassWord == resultUser.Password) && (loginReq.PassWord != "") {
		resultBool = true
	}
	return resultBool, resultUser
}

// EditUserReq 更新用户信息数据类
type EditUserReq struct {
	UserId     string `json:"userId"`
	UserName   string `json:"userName"`
	UserEmail  string `json:"userEmail"`
}

// UpdateUser 更新用户信息
func UpdateUser(user model.User) bool {
	err := db.UpdateUserInfo(user)
	return err == nil
}
