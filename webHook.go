package main

import "github.com/parnurzeal/gorequest"

func webHook() {
	request := gorequest.New()
	data := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": "已经获得题目，请尽快做题",
		},
	}
	request.Post("https://open.feishu.cn/open-apis/bot/v2/hook/11db536e-f92a-4f1e-a091-cc17463b52fb").Type("json").
		Send(data).
		End()
}
