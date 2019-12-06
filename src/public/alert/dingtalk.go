package alert

import (
	"bytes"
	"crypto/tls"
	"datastructure"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ContentType = "application/json"
	DingTalkUrl = "https://oapi.dingtalk.com/robot/send?access_token=edaef5c4adce3689970145a797da0a33f1da05b12f3573bc9b8196c2f844f6d6"
)

func Ding(a datastructure.Request) (err error) {
	var (
		b, bodyContentByte                    []byte
		subject, textcontent, markdowncontent string
		f                                     [1]string
	)
	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 数据定义
	subject = "        乐湃事件通知"
	markdowncontent = "## " + subject +
		"\n" + "### **" + a.JavaProject + "**项目滚动更新完成 \n" +
		"\n" + "1. 项目版本：" + a.Version +
		"\n" + "2. 镜像版本：" + a.Image +
		"\n" + "3. 更新备注：" + a.Info.UpdateSummary +
		"\n" + "4. 执行人：" + a.Info.RequestMan +
		"\n@" + a.Info.PhoneNumber

	textcontent = subject +
		"\n" + "{ " + a.JavaProject + " } 滚动更新完成" +
		"\n" + "项目版本：" + a.Version +
		"\n" + "镜像版本：" + a.Image +
		"\n" + "更新备注：" + a.Info.UpdateSummary +
		"\n" + "执行人：" + a.Info.RequestMan +
		"\n@" + a.Info.PhoneNumber

	f[0] = a.Info.PhoneNumber
	if a.SendFormat == "text" {
		var d = datastructure.DingText{
			"text",
			datastructure.Text{
				textcontent,
			},
			datastructure.At{
				f,
				"false",
			},
		}
		if b, err = json.Marshal(d); err == nil {
			log.Printf("Send TO DingTalk %v ", string(b))
		}
	} else {
		var d = datastructure.DingMarkDown{
			"markdown",
			datastructure.MarkDown{
				subject,
				markdowncontent,
			},
			datastructure.At{
				f,
				"false",
			},
		}
		if b, err = json.Marshal(d); err == nil {
			log.Printf("Send TO DingTalk %v ", string(b))
		}
	}

	body := new(bytes.Buffer)
	body.ReadFrom(bytes.NewBuffer([]byte(b)))

	client := &http.Client{Transport: tr}
	requestGet, _ := http.NewRequest("POST", DingTalkUrl, body)
	requestGet.Header.Add("Content-Type", ContentType)
	resp, err := client.Do(requestGet)
	if err != nil {
		log.Printf("Get Request Failed ERR:[%s]", err.Error())
		err = fmt.Errorf("Get Request Failed ERR:[%s]", err.Error())
		return
	}
	bodyContentByte, err = ioutil.ReadAll(resp.Body)
	StatusCode := resp.StatusCode
	bodyContent := string(bodyContentByte)
	fmt.Println(StatusCode, bodyContent)
	return
}
