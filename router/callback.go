package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"dawn_api/log"
)

const (
	requestOk = 10000 + iota
	requestParamError
)

type EdgeRequest struct {
	ShopID        string   `json:"shop_id"`
	StartTime     int      `json:"start_time"`
	EndTime       int      `json:"end_time"`
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

//Trajectory event callback
func TrajectoryCallback(c *gin.Context) {
	requestRawData, _ := c.GetRawData()
	request := string(requestRawData)
	log.Info("Traj API Trajectory Request ", zap.ByteString("request", []byte(request)))

	// edgeReq := &EdgeRequest{}
	// err = json.Unmarshal([]byte(request), &edgeReq)
	// if err != nil {
	// 	response, _ := normalResponse(requestParamError, requestID)
	// 	c.String(200, response)
	// 	return
	// }

	// sn_list := api.GetSnListByShopId(edgeReq.ShopID)
	// traj_event := db.QueryTraj(sn_list, edgeReq.StartTime, edgeReq.EndTime)
	// fmt.Println(traj_event)
	// out_traj_map := make(map[string][]db.Traj_Event)
	// for _, event := range traj_event {
	// 	out_traj_map[event.Reid] = append(out_traj_map[event.Reid], event)
	// }
	// out_traj_str, _ := json.Marshal(out_traj_map)
	// fmt.Println(string(out_traj_str))
	// response, err := returnResponse(requestOk, requestID, string(out_traj_str))
	// c.String(200, response)
	return
}

func normalResponse(status int, requestID string) (string, error) {
	responseStr, err := json.Marshal(errorResponse{Code: status, RequestID: requestID})
	if err != nil {
		log.Error("Traj API response Error", zap.Any("error", err))
		return "", err
	}
	log.Info("Traj API normal response ", zap.ByteString("response", responseStr))
	return string(responseStr), nil
}

func returnResponse(status int, requestID string, result string) (string, error) {
	responseStr, err := json.Marshal(Response{Code: status,
		RequestID: requestID,
		Result: result})
	if err != nil {
		log.Error("Traj API response Error", zap.Any("error", err))
		return "", err
	}

	log.Info("Traj API response sucess", zap.ByteString("response", responseStr))
	return string(responseStr), nil
}
