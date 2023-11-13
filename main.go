package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"math/rand"
	"time"
)

func main() {
	request := gorequest.New()
	data := map[string]string{
		"code":        Code,
		"userName":    UserName,
		"pwd":         Pwd,
		"fromChannel": "wechat",
		"token":       "",
		"mpversion":   "202105171620",
	}

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
		total, orderId, uniqueId, urlList, text := tasklist(preURL, token.(string))
		if total == 0 {
			fmt.Println("no task not found in response")
			fmt.Println("前方还有", queueposition(preURL, token.(string)), "人")
		} else {
			webHook(urlList, text)
			fmt.Println("获得题目,放弃请输入g")
			//获得输入
			var input string
			fmt.Scanln(&input)
			if input == "g" {
				answer(preURL, token.(string), orderId, uniqueId)
				fmt.Println("已经放弃，请继续答题")
			}
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
func tasklist(preURL string, token string) (float64, string, string, []string, string) {
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
		return 0, "", "", nil, ""
	}
	// 解析 JSON 响应体
	var responseMap map[string]interface{}
	err := json.Unmarshal([]byte(body), &responseMap)
	// 检查是否有错误发生
	if err != nil {
		fmt.Println("Error:", err)
		return 0, "", "", nil, ""
	}
	// 获取 total 值
	total, totalExists := responseMap["data"].(map[string]interface{})["total"]
	if !totalExists {
		//fmt.Println("Total not found in response")
		return 0, "", "", nil, ""
	}
	if total.(float64) < 1.0 {
		return 0, "", "", nil, ""
	}
	list := responseMap["data"].(map[string]interface{})["list"].([]interface{})
	newList := list[0].(map[string]interface{})
	orderId := newList["orderId"]
	uniqueId := newList["uniqueId"]
	content := newList["content"].(map[string]interface{})
	urlList := make([]string, len(content["urlList"].([]interface{})))
	for i, i2 := range content["urlList"].([]interface{}) {
		urlList[i] = fmt.Sprintf("%v", i2)
	}
	text := fmt.Sprintf("%v", content["text"])
	return total.(float64), orderId.(string), uniqueId.(string), urlList, text
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

// 访问 /commitui/api/answer
func answer(preURL string, token string, orderId string, uniqueId string) {
	///访问 commitui/api/answer
	request := gorequest.New()
	//answerId=2327679671&orderId=2328226013&answerAction=2&businessId=839bc51F69&reason=2602&remark=%E6%97%A0&retry=0&token=328edef7259ad8dd1c04bad3fdf497c1f2aa6a9c
	data := map[string]string{
		"answerId":     uniqueId,
		"orderId":      orderId,
		"answerAction": "2",
		"businessId":   GetRandomString(),
		"reason":       "2602",
		"remark":       "无",
		"retry":        "0",
		"token":        token,
	}
	_, _, errs := request.Post(preURL + "/commitui/api/answer").Type("form").
		Send(data).
		End()
	// 检查是否有错误发生
	if errs != nil {
		for _, err := range errs {
			fmt.Println("Error:", err)
		}
	}
}

// 生成十位只包含大小写字母和数字的随机字符串
func GetRandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
