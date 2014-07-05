/**
 * 百度云推送SDK, Channel扩展
 * User: zha_ming@163.com
 * Package: baidupush
 * Date: 14-6-17 11:41
 * Version: 0.1
 */
package baidupush

import (
	"strings"
	"strconv"
	"io/ioutil"
)

type deviceT map[string]int
var deviceMap = deviceT{
	"browser":	1,
	"pc":		2,
	"android":	3,
	"ios":		4,
	"winphone":	5,
}

type Push struct {
	*Channel
}

func (this *Push) PushMessage2(deviceType string, pushType int, pushParam string, messages Opt, deployed bool) bool {
	deviceType = strings.ToLower(deviceType)
	var devices deviceT
	if deviceId,ok := deviceMap[deviceType]; ok {
		devices[deviceType] = deviceId
	} else if deviceType == "all" {
		devices = deviceMap
	}

	opt := make(Opt)
	switch pushType {
	case PUSH_TO_USER:
		opt[USER_ID] = pushParam
	case PUSH_TO_TAG:
		opt[TAG_NAME] = pushParam
	case PUSH_TO_ALL:
		fallthrough
	default:
	}

	rs := true
	for k,v := range devices {
		opt[DEVICE_TYPE] = strconv.Itoa(v)
		if k == "ios" {
			opt[MESSAGE_TYPE] = "1"
			if deployed {
				opt[DEPLOY_STATUS] = "2"
			} else {
				opt[DEPLOY_STATUS] = "1"
			}
		} else {
			opt[MESSAGE_TYPE] = messages["type"]
		}
		_, err := this.PushMessage(pushType, messages["msg"], messages["key"], opt)
		rs = rs && (err == nil)
	}
	return rs
}

func (this *Push) SetTag2(tagName, userId string) bool {
	opt := make(Opt)
	opt[USER_ID] = userId
	_, err := this.SetTag(tagName, opt)
	return err == nil
}

func (this *Push) FetchTag2(tagName string) bool {
	opt := make(Opt)
	if len(tagName) > 0 {opt[TAG_NAME] = tagName}
	_, err := this.FetchTag(opt)
	return err == nil
}

func (this *Push) InitAppIoscert2(name,desc,relCertF,devCertF string, deployed bool) bool {
	opt := make(Opt)
	if deployed {
		opt[DEPLOY_STATUS] = "2"
	} else {
		opt[DEPLOY_STATUS] = "1"
	}
	if _,err := this.QueryAppIoscert(opt); err ==nil {return true}

	relCert,err := ioutil.ReadFile(relCertF)
	if err != nil {println(err);return false}
	devCert,err := ioutil.ReadFile(devCertF)
	if err != nil {println(err);return false}

	_, err = this.InitAppIoscert(name, desc, string(relCert), string(devCert), opt)
	return err == nil
}

func (this *Push) UpdateAppIoscert2(name,desc,relCertF,devCertF string, deployed bool) bool {
	relCert,err := ioutil.ReadFile(relCertF)
	if err != nil {println(err);return false}
	devCert,err := ioutil.ReadFile(devCertF)
	if err != nil {println(err);return false}
	opt := make(Opt)
	if deployed {
		opt[DEPLOY_STATUS] = "2"
	} else {
		opt[DEPLOY_STATUS] = "1"
	}
	opt[NAME] = name
	opt[DESC] = desc
	opt[REL_CERT] = string(relCert)
	opt[DEV_CERT] = string(devCert)
	_, err = this.UpdateAppIoscert(opt)
	return err == nil
}


//SDK入口
func NewPush(apiKey,secretKey string, curlOpts []string) *Push {
	err := NewError(errorMsgMap)
	return &Push{&Channel{
		ApiKey: apiKey,
		SecretKey: secretKey,
		CurlOpts: curlOpts,
		Err: err,
	}}
}
