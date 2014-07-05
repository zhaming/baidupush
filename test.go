/**
 * 测试百度推送SDK
 * User: zha_ming@163.com
 * Package: main
 * Date: 14-6-17 10:54
 * Version: 0.1
 */
package main

import (
	."./baidupush"
	"time"
)

func main() {
	var curlOpt []string
	opt := make(Opt)
	apiKey := ""
	//apiKey := "123"
	secretKey := ""
	//userId := "813292011209507601";
	//userId := "1148340223945189762";
	userId := "868725677998357477"; //iPhone
	//userId := "856690259366666776"; //HUAWEI
	//channelId := "3891701151187579171";
	var messages = Opt{
		"type": "1",  //0:消息 1:通知
		"key": "msg_key",
		"msg": "{\"title\": \"来自GOLANG的通知\",\"description\": \"爆料：小米论坛数据库800万用户资料泄漏！尽情期待\"," +
			"\"notification_basic_style\":7,\"open_type\":1,\"user_confirm\":1," +
			"\"url\":\"http://www.mi.com\",\"aps\":{\"sound\":\"\",\"badge\":1}}",
	}
	//opt["channel_id"] = "4480248167080268496"

	push := NewPush(apiKey, secretKey, curlOpt)
	push.QueryBindList(userId, opt)
	time.Sleep(100 * time.Millisecond)
	push.VerifyBind(userId, opt)
	time.Sleep(100 * time.Millisecond)
	push.SetTag2("group5", userId)
	time.Sleep(100 * time.Millisecond)
	push.QueryUserTags(userId, opt)  //用户自定义标签
	time.Sleep(100 * time.Millisecond)
	//push.DeleteTag("group5", opt)
	time.Sleep(100 * time.Millisecond)
	push.FetchTag2("")
	time.Sleep(100 * time.Millisecond)
	push.FetchMessage(userId, opt)
	time.Sleep(100 * time.Millisecond)
	push.FetchMessageCount(userId, opt)
	time.Sleep(100 * time.Millisecond)
	//push.InitAppIoscert2("APNs-go-test", "the demo for APNs", "./APNs-pro.pem", "./APNs-dev.pem", false)
	time.Sleep(100 * time.Millisecond)
	push.QueryAppIoscert(opt);
	time.Sleep(100 * time.Millisecond)
	push.PushMessage2("all", 2, "group3", messages, false)
}
