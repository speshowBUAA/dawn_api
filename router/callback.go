package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"dawn_api/log"
	"dawn_api/db"
	"dawn_api/api"
	"dawn_api/model"
	"net/http"
	myjwt "dawn_api/middleware"
	"time"
	jwtgo "github.com/dgrijalva/jwt-go"
	// "fmt"
	"strconv"
)

const (
	requestOk = 10000 + iota
	requestParamError

	per_page = 10
)

type Anno struct {
	src        string   `json:"src"`
	name       string   `json:"name"`
	regions    string   `json:"regions"`
	pixelSize  string	`json:"pixelSize"`
}

type Response struct {
	Code       int       `json:"code"`
	RequestID  string    `json:"requestId"`
	Result     string    `json:"result"`
}

type errorResponse struct {
	Code       int       `json:"code"`
	RequestID  string    `json:"requestId"`
}

func SaveAnnotationCallback(c *gin.Context){
	requestRawData, _ := c.GetRawData()
	request := string(requestRawData)
	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	if claims != nil {
		var m []interface{}
		err := json.Unmarshal([]byte(request), &m)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"msg":    "json解析失败！",
			})
			return
		}
		
		db.InsertMany(m)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    "保存成功！",
		})
	}
}

// //Trajectory event callback
// func TrajectoryCallback(c *gin.Context) {
// 	requestRawData, _ := c.GetRawData()
// 	request := string(requestRawData)
// 	log.Info("Traj API Trajectory Request ", zap.ByteString("request", []byte(request)))

// 	// edgeReq := &EdgeRequest{}
// 	// err = json.Unmarshal([]byte(request), &edgeReq)
// 	// if err != nil {
// 	// 	response, _ := normalResponse(requestParamError, requestID)
// 	// 	c.String(200, response)
// 	// 	return
// 	// }

// 	// sn_list := api.GetSnListByShopId(edgeReq.ShopID)
// 	// traj_event := db.QueryTraj(sn_list, edgeReq.StartTime, edgeReq.EndTime)
// 	// fmt.Println(traj_event)
// 	// out_traj_map := make(map[string][]db.Traj_Event)
// 	// for _, event := range traj_event {
// 	// 	out_traj_map[event.Reid] = append(out_traj_map[event.Reid], event)
// 	// }
// 	// out_traj_str, _ := json.Marshal(out_traj_map)
// 	// fmt.Println(string(out_traj_str))
// 	// response, err := returnResponse(requestOk, requestID, string(out_traj_str))
// 	// c.String(200, response)
// 	return
// }

// func normalResponse(status int, requestID string) (string, error) {
// 	responseStr, err := json.Marshal(errorResponse{Code: status, RequestID: requestID})
// 	if err != nil {
// 		log.Error("Traj API response Error", zap.Any("error", err))
// 		return "", err
// 	}
// 	log.Info("Traj API normal response ", zap.ByteString("response", responseStr))
// 	return string(responseStr), nil
// }

// func returnResponse(status int, requestID string, result string) (string, error) {
// 	responseStr, err := json.Marshal(Response{Code: status,
// 		RequestID: requestID,
// 		Result: result})
// 	if err != nil {
// 		log.Error("Traj API response Error", zap.Any("error", err))
// 		return "", err
// 	}

// 	log.Info("Traj API response sucess", zap.ByteString("response", responseStr))
// 	return string(responseStr), nil
// }

// 注册信息
type RegistInfo struct {
	// 用户名
	Name string `json:"username"`
	// 密码
	Pwd string `json:"password"`
}

// Register 注册用户
func RegisterUser(c *gin.Context) {
	var registerInfo RegistInfo
	if c.BindJSON(&registerInfo) == nil {
		err := api.Register(registerInfo.Name, registerInfo.Pwd)
		registStr, _:= json.Marshal(registerInfo)
		if err == nil {
			log.Info("Regist Event: ", zap.ByteString("regist", registStr))
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"msg":    "注册成功！",
			})
		} else {
			log.Info("Regist Event: ", zap.ByteString("regist", registStr), zap.ByteString("error", []byte(err.Error())))
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg": err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "解析数据失败！",
		})
	}
}

// LoginResult 登录结果结构
type LoginResult struct {
	Token string `json:"token"`
	model.User
}

// Login 登录
func Login(c *gin.Context) {
	var loginReq model.LoginReq
	if c.BindJSON(&loginReq) == nil {
		isPass, user := api.LoginCheck(loginReq)
		if isPass {
			generateToken(c, user)
			isUpdate := api.UpdateUser(user)
			if isUpdate {
				userStr, _:= json.Marshal(user)
				log.Info("Login Event: ", zap.ByteString("user", userStr))
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    "验证失败",
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "json 解析失败",
		})
	}
}

// 生成令牌
func generateToken(c *gin.Context, user model.User) {
	j := &myjwt.JWT{
		[]byte("leapmotor"),
	}
	claims := myjwt.CustomClaims{
		user.Id,
		user.Name,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "leapmotor",                   //签名的发行者
		},
	}

	var auth string
	if user.Name == "admin" {
		auth = "admin"
	} else {
		auth = "user"
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "登录成功！",
		"access_token":   token,
		"currentAuthority": auth,
		"currentUser": user.Name,
	})
	return
}

//获取当前用户
func CurrentUser(c *gin.Context) {
	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	if claims != nil {
		currentUser := db.Query("username", string(claims.Name))[0]
		c.JSON(http.StatusOK, gin.H{
			"userid": string(currentUser.Id),
			"name": currentUser.Name,
			"email": currentUser.Email,
			"phone": currentUser.Phone,
			"permission": currentUser.Permission,
			"last_login_time": currentUser.Last_login_time,
			"create_time": currentUser.Create_time,
			"avatar": currentUser.Avatar_url,
		})
	}
}

type linkinfo struct {
	Previous string		`json:"previous"`
	Next     string		`json:"next"`
}

type pageInfo struct {
	Per_page         int 		`json:"per_page"`
	Current_page     int		`json:"current_page"`
	Total_page       int		`json:"total_page"`
	Links            linkinfo	`json:"links"`
}

type UserListInfo struct {
	Data []model.User   `json:"data"`
	Meta pageInfo		`json:"meta"`
}

//获取用户列表
func GetUsers(c *gin.Context) {
	// current_page := c.Query("current")
	per_page, _:= strconv.Atoi(c.Query("pageSize"))
	var param model.FilterParam
	param.UserName = c.Query("username")
	if (param.UserName == "") {param.UserName = "%"}
	param.UserEmail = c.Query("useremail")
	if (param.UserEmail == "") {param.UserEmail = "%"}
	param.Permission = c.Query("permission")
	if ((param.Permission == "") || (param.Permission == "all")) {
		param.Permission = "true,false"
	} else if (param.Permission == "permission") {
		param.Permission = "true"
	} else if (param.Permission == "forbidden") {
		param.Permission = "false"
	}
	param.StartCreateDate = c.Query("startCreateDate")
	if (param.StartCreateDate == "") {param.StartCreateDate = "2000-01-01"}
	param.EndCreateDate = c.Query("endCreateDate")
	if (param.EndCreateDate == "") {param.EndCreateDate = "NOW()"}
	param.StartLoginTime = c.Query("startLoginTime")
	if (param.StartLoginTime == "") {param.StartLoginTime = "2000-01-01 00:00:00"}
	param.EndLoginTime = c.Query("endLoginTime")
	if (param.EndLoginTime == "") {param.EndLoginTime = "NOW()"}

	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	if (claims != nil) && (claims.Name == "admin") {
		UserList := db.FilterQuery(param)
		total_page := (len(UserList) / (per_page+1)) + 1
		pageInfo := pageInfo{Per_page:per_page, Total_page: total_page}
		userlistinfo := UserListInfo{Data:UserList, Meta:pageInfo}
		userlistinfo_str, _ := json.Marshal(userlistinfo)
		c.String(http.StatusOK, string(userlistinfo_str))
	} else {
		c.String(http.StatusForbidden, "没有权限！")
	}
}

// GetDataByTime 一个需要token认证的测试接口
func GetDataByTime(c *gin.Context) {
	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "token有效",
			"data":   claims,
		})
	}
}