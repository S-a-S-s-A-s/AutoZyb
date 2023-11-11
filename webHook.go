package main

import "github.com/parnurzeal/gorequest"

type reqStruct struct {
	ZhCn struct {
		Title   string `json:"title"`
		Content [][]struct {
			Tag  string `json:"tag"`
			Text string `json:"text"`
			Href string `json:"href,omitempty"`
		} `json:"content"`
	} `json:"zh_cn"`
}

func webHook(urlList []string, text string) {
	request := gorequest.New()

	for i := 0; i < len(urlList); i++ {
		data := reqStruct{}
		data.ZhCn.Title = "有新的题目"
		data.ZhCn.Content = [][]struct {
			Tag  string `json:"tag"`
			Text string `json:"text"`
			Href string `json:"href,omitempty"`
		}{
			{
				{
					Tag:  "text",
					Text: text,
				},
			},
			{
				{
					Tag:  "a",
					Text: "点击查看",
					Href: urlList[i],
				},
			},
		}

		request.Post("https://open.feishu.cn/open-apis/bot/v2/hook/11db536e-f92a-4f1e-a091-cc17463b52fb").Type("json").
			Send(data).
			End()
	}

}
