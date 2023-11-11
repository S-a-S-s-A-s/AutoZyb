package main

import "github.com/parnurzeal/gorequest"

type req struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCN struct {
				Title   string `json:"title"`
				Content [][]struct {
					Tag  string `json:"tag"`
					Text string `json:"text"`
					Href string `json:"href,omitempty"`
				} `json:"content"`
			} `json:"zh-CN"`
		} `json:"post"`
	} `json:"content"`
}

func webHook(urlList []string, text string) {
	request := gorequest.New()

	for i := 0; i < len(urlList); i++ {
		data := req{}
		data.MsgType = "post"
		data.Content.Post.ZhCN.Title = "有新的题目"
		data.Content.Post.ZhCN.Content = [][]struct {
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

		_, _, errors := request.Post("https://open.feishu.cn/open-apis/bot/v2/hook/11db536e-f92a-4f1e-a091-cc17463b52fb").Type("json").
			Send(data).
			End()
		if errors != nil {
			println(errors)
		}
	}

}
