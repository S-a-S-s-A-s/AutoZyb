package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

func main() {
	request := gorequest.New()
	//code=0d3uuvFa1db6dG0nT1Ga1z1ERv0uuvFE&userName=17773144837&pwd=bfa2591265651fc0b0e85afe29524244&fromChannel=wechat&token=&mpversion=202105171620
	data := map[string]string{
		"code":        Code,
		"userName":    UserName,
		"pwd":         Pwd,
		"fromChannel": "wechat",
		"token":       "",
		"mpversion":   "202105171620",
	}
	//data := map[string]string{
	//	"code":        "0d3uuvFa1db6dG0nT1Ga1z1ERv0uuvFE",
	//	"userName":    "17773144837",
	//	"pwd":         "bfa2591265651fc0b0e85afe29524244",
	//	"fromChannel": "wechat",
	//	"token":       "",
	//	"mpversion":   "202105171620",
	//}
	preURL := "https://wenda.zuoyebang.com/"
	_, body, errs := request.Post(preURL + "/commitui/session/login").Type("form").
		Send(data).
		End()

	// 检查是否有错误发生
	if errs != nil {
		for _, err := range errs {
			fmt.Println("Error:", err)
		}
	}

	// 解析 JSON 响应体
	var responseMap map[string]interface{}
	err := json.Unmarshal([]byte(body), &responseMap)
	// 检查是否有错误发生
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 获取 token 值
	token, tokenExists := responseMap["data"].(map[string]interface{})["token"]
	if !tokenExists {
		fmt.Println("Token not found in response")
		return
	}
	// 获取 logId 值
	_, logIdExists := responseMap["logId"]
	if !logIdExists {
		fmt.Println("LogId not found in response")
		return
	}
	go escExit(preURL, token.(string))
	//defer dequeueorder(preURL, token.(string))

	//访问 /commitui/question/rushorder
	errors := rushorder(preURL, token.(string))
	// 检查是否有错误发生
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Error:", err)
		}
		return
	}
	// 启动一个协程，每隔十秒发送一次心跳
	go func() {
		for {
			time.Sleep(9 * time.Second)
			heartbeat(preURL, token.(string))
		}
	}()
	//访问 rui/ask/tasklist
	for {
		time.Sleep(2 * time.Second)
		total := tasklist(preURL, token.(string))
		if total == 0 {
			fmt.Println("no task not found in response")
			fmt.Println("前方还有", queueposition(preURL, token.(string)), "人")
		} else {
			webHook()
			fmt.Println("获得题目")
			return
		}
	}

}

// 访问 /commitui/question/rushorder
func rushorder(preURL string, token string) []error {
	//访问 /commitui/question/rushorder
	request := gorequest.New()
	//businessId=857b64aaBc&taskFrom=1&pools=1&priority=&duration=3599&containPL=1&token=835af0db712aeae4a57842056a2ff3d304a9c46d
	data := map[string]string{
		"businessId": "857b64aaBc",
		"taskFrom":   "1",
		"pools":      "1",
		"priority":   "",
		"duration":   "9999",
		"containPL":  "1",
		"token":      token,
	}
	_, _, errs := request.Post(preURL + "/commitui/question/rushorder").Type("form").
		Send(data).
		End()
	// 检查是否有错误发生
	return errs
}

// 访问 /rui/ask/tasklist
func tasklist(preURL string, token string) float64 {
	// 访问 /rui/ask/tasklist
	request := gorequest.New()
	//taskType=1&token=835af0db712aeae4a57842056a2ff3d304a9c46d
	data := map[string]string{
		"taskType": "1",
		"token":    token,
	}
	_, body, errs := request.Post(preURL + "/rui/ask/tasklist").Type("form").
		Send(data).
		End()
	// 检查是否有错误发生
	if errs != nil {
		fmt.Println("Error:", errs)
		return 0
	}
	// 解析 JSON 响应体
	var responseMap map[string]interface{}
	err := json.Unmarshal([]byte(body), &responseMap)
	// 检查是否有错误发生
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	// 获取 total 值
	total, totalExists := responseMap["data"].(map[string]interface{})["total"]
	if !totalExists {
		fmt.Println("Total not found in response")
		return 0
	}
	return total.(float64)
}

// 访问 /commitui/question/dequeueorder
func dequeueorder(preURL string, token string) {
	///访问 commitui/question/dequeueorder
	request := gorequest.New()
	data := map[string]string{
		"token": token,
	}
	_, _, errs := request.Post(preURL + "/commitui/question/dequeueorder").Type("form").
		Send(data).
		End()
	// 检查是否有错误发生
	if errs != nil {
		for _, err := range errs {
			fmt.Println("Error:", err)
		}
	}
	fmt.Println("已经退出队列")
}

// 访问 /commitui/question/queueposition
func queueposition(preURL string, token string) int {
	///访问 commitui/question/queueposition
	request := gorequest.New()
	data := map[string]string{
		"token": token,
	}
	_, body, errs := request.Post(preURL + "/commitui/question/queueposition").Type("form").
		Send(data).
		End()
	// 检查是否有错误发生
	if errs != nil {
		for _, err := range errs {
			fmt.Println("Error:", err)
		}
		return -1
	}
	// 解析 JSON 响应体
	var responseMap map[string]interface{}
	err := json.Unmarshal([]byte(body), &responseMap)
	// 检查是否有错误发生
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}
	// 获取 numInFront 值
	numInFront, numInFrontExists := responseMap["data"].(map[string]interface{})["numInFront"]
	if !numInFrontExists {
		fmt.Println("NumInFront not found in response")
		return -1
	}
	return int(numInFront.(float64))
}

// 访问 /commitui/regist/heartbeat
func heartbeat(preURL string, token string) {
	///访问 commitui/regist/heartbeat
	request := gorequest.New()
	//token=5b8cace2ae379b3d169371dc099baebf3afcd0b7&mpversion=202105171620
	data := map[string]string{
		"token":     token,
		"mpversion": "202105171620",
	}
	_, _, errs := request.Post(preURL + "/commitui/regist/heartbeat").Type("form").
		Send(data).
		End()
	// 检查是否有错误发生
	if errs != nil {
		for _, err := range errs {
			fmt.Println("Error:", err)
		}
	}
}
