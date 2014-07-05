/**
 * 百度云推送SDK
 * User: zha_ming@163.com
 * Package: baidupush
 * Date: 14-6-17 11:42
 * Version: 0.1
 */
package baidupush

import (
	//"os"
	"fmt"
	"io"
	"time"
	"bytes"
	"strings"
	"strconv"
	"sort"
	"crypto/md5"
	u "net/url"
	json "github.com/bitly/go-simplejson"
)

const (
	USER_ID = "user_id"  //用户类型
	TAG_NAME = "tag"  //消息标签名，可按标签分组
	PUSH_TYPE = "push_type"
	MESSAGES = "messages"  //消息
	MSG_KEYS = "msg_keys"  //消息key，可以按key去除重复消息
	MSG_IDS = "msg_ids"  //消息id
	DEVICE_TYPE = "device_type" //设备类型 1：浏览器设备；2：PC设备；3：andorid设备
	MESSAGE_TYPE = "message_type" //消息类型 0：默认类型
	DEPLOY_STATUS = "deploy_status"
	PUSH_TO_USER = 1
	PUSH_TO_TAG = 2
	PUSH_TO_ALL = 3
	NAME = "name"
	DESC = "description"
	REL_CERT = "release_cert"
	DEV_CERT = "dev_cert"

	BASEURL = "http://channel.api.duapp.com/rest/2.0/channel/"
	CHANNEL_ID = "channel_id"  //消息通道ID号
	TIMESTAMP = "timestamp"  //发起请求时的时间戳
	API_KEY = "apikey"  //应用key，从百度开发者中心获得,是创建Channel的必须参数
	SECRET_KEY = "secret_key"  //从百度开发者中心获得，是创建Channel的必须参数
	SIGN = "sign"  //Channel常量，用户不必关注
	METHOD = "method"
	HTTP_METHOD = "POST"
	USERAGENT = "RequestCore/1.4.2"
	DEBUG = true

	//Error Code
	CHANNEL_SDK_ERR = 0
	CHANNEL_SDK_SYS = 1
	CHANNEL_SDK_INIT_FAIL = 2
	CHANNEL_SDK_PARAM = 3
	CHANNEL_SDK_HTTP_STATUS_ERROR_AND_RESULT_ERROR = 4
	CHANNEL_SDK_HTTP_STATUS_OK_BUT_RESULT_ERROR = 5
)

var errorMsgMap = []string{
	CHANNEL_SDK_ERR: "sdk error",
	CHANNEL_SDK_SYS: "sdk error",
	CHANNEL_SDK_INIT_FAIL: "sdk init error",
	CHANNEL_SDK_PARAM: "lack param",
	CHANNEL_SDK_HTTP_STATUS_ERROR_AND_RESULT_ERROR: "http status is error, and the body returned is not a json string",
	CHANNEL_SDK_HTTP_STATUS_OK_BUT_RESULT_ERROR: "http status is ok, but the body returned is not a json string",
}

type MethodChannelArray []string
func (this MethodChannelArray) InArray(s string) bool {
	if len(s) == 0 {return false}
	for _, v := range this {
		if v == s {return true}
	}
	return false
}
var methodChannelInBody = MethodChannelArray{"push_msg","set_tag","fetch_tag","delete_tag","query_user_tags"}

var headersMap = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"User-Agent": "Baidu Channel Service Gosdk Client",
}

type Opt map[string]string
type Res map[string]interface{}
type Arg map[interface{}]interface{}

// for sort
type SortT []interface{}
func (this SortT) Len() int {return len(this)}
func (this SortT) Less(i, j int) bool {
	stri := this[i].(string)
	strj := this[j].(string)
	if bytes.Compare([]byte(stri), []byte(strj)) <= 0 {return true}
	return false
}
func (this SortT) Swap(i, j int) {this[i], this[j] = this[j], this[i]}

//TODO:some comment
type Channel struct {
	ApiKey		string
	SecretKey	string
	CurlOpts	[]string
	Err			PushError
	RequestId 	int
}

func (this *Channel) SetApiKey(apiKey string) {
	this.ApiKey = apiKey
}

func (this *Channel) SetSecretKey(secretKey string) {
	this.SecretKey = secretKey
}

//服务器端根据userId, 查询绑定信息
func (this *Channel) QueryBindList(userId string, opt Opt) (Res, error) {
	arrNeed := []string{USER_ID}
	tmpArgs := Arg{0:userId, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "query_bindlist"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//校验userId是否已经绑定
func (this *Channel) VerifyBind(userId string, opt Opt) (Res, error) {
	arrNeed := []string{USER_ID}
	tmpArgs := Arg{0:userId, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "verify_bind"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//推送消息
func (this *Channel) PushMessage(pushType int, messages,messageKeys string, opt Opt) (Res, error) {
	arrNeed := []string{PUSH_TYPE, MESSAGES, MSG_KEYS}
	tmpArgs := Arg{0:strconv.Itoa(pushType), 1:messages, 2:messageKeys, 3:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}

	switch pushType {
	case PUSH_TO_USER:
		if uid,ok := arrArgs[USER_ID]; !ok || uid == nil {
			this.Err.Msg = "userId should be specified in optional when pushType is " + strconv.Itoa(pushType)
			this.Err.Code = CHANNEL_SDK_PARAM
			pError(this.Err.Error())
		}
	case PUSH_TO_TAG:
		if tag,ok := arrArgs[TAG_NAME]; !ok || tag == nil {
			this.Err.Msg = "tagName should be specified in optional when pushType is " + strconv.Itoa(pushType)
			this.Err.Code = CHANNEL_SDK_PARAM
			pError(this.Err.Error())
		}
	case PUSH_TO_ALL:
		fallthrough
	default:
		this.Err.Msg = "pushType value is invalid"
		this.Err.Code = CHANNEL_SDK_PARAM
		pError(this.Err.Error())
	}
	arrArgs[METHOD] = "push_msg"
	arrArgs[PUSH_TYPE] = pushType
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//根据userId查询消息
func (this *Channel) FetchMessage(userId string, opt Opt) (Res, error) {
	arrNeed := []string{USER_ID}
	tmpArgs := Arg{0:userId, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "fetch_msg"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//根据userId查询消息个数
func (this *Channel) FetchMessageCount(userId string, opt Opt) (Res, error) {
	arrNeed := []string{USER_ID}
	tmpArgs := Arg{0:userId, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "fetch_msgcount"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//根据userId, msgIds删除消息
func (this *Channel) DeleteMessageCount(userId string, msgId []string, opt Opt) (Res, error) {
	//array to json
	/*var msgIds string
	if len(msgId) > 0 {
		js := json.New()
		for k,v := range msgId {js.Set(strconv.Itoa(k), v)}
		msgIds, _ = js.Encode()
	}*/
	//join array
	msgIds := strings.Join(msgId, ",")

	arrNeed := []string{USER_ID, MSG_IDS}
	tmpArgs := Arg{0:userId, 1:msgIds, 2:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "delete_msg"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//设置消息标签
func (this *Channel) SetTag(tagName string, opt Opt) (Res, error) {
	arrNeed := []string{TAG_NAME}
	tmpArgs := Arg{0:tagName, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "set_tag"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//查询消息标签信息
func (this *Channel) FetchTag(opt Opt) (Res, error) {
	arrNeed := []string{}
	tmpArgs := Arg{0:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "fetch_tag"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//删除消息标签
func (this *Channel) DeleteTag(tagName string, opt Opt) (Res, error) {
	arrNeed := []string{TAG_NAME}
	tmpArgs := Arg{0:tagName, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "delete_tag"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//查询用户相关的标签
func (this *Channel) QueryUserTags(userId string, opt Opt) (Res, error) {
	arrNeed := []string{USER_ID}
	tmpArgs := Arg{0:userId, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "query_user_tags"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//根据channelId查询设备类型
func (this *Channel) QueryDeviceType(channelId string, opt Opt) (Res, error) {
	arrNeed := []string{CHANNEL_ID}
	tmpArgs := Arg{0:channelId, 1:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "query_device_type"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//初始化应用ios证书
func (this *Channel) InitAppIoscert(name,desc,relCert,devCert string, opt Opt) (Res, error) {
	arrNeed := []string{NAME, DESC, REL_CERT, DEV_CERT}
	tmpArgs := Arg{0:name, 1:desc, 2:relCert, 3:devCert, 4:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "init_app_ioscert"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//修改ios证书内容
func (this *Channel) UpdateAppIoscert(opt Opt) (Res, error) {
	arrNeed := []string{}
	tmpArgs := Arg{0:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "update_app_ioscert"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//查询ios证书内容
func (this *Channel) QueryAppIoscert(opt Opt) (Res, error) {
	arrNeed := []string{}
	tmpArgs := Arg{0:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "query_app_ioscert"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}

//删除ios证书内容
func (this *Channel) DeleteAppIoscert(opt Opt) (Res, error) {
	arrNeed := []string{}
	tmpArgs := Arg{0:opt}
	arrArgs, err := this.mergeArgs(arrNeed, tmpArgs)
	if err != nil {pError(err.Error())}
	arrArgs[METHOD] = "delete_app_ioscert"
	rs, err := this.commonProccess(arrArgs)
	if err != nil {pError(err.Error())}
	this.output(arrArgs[METHOD].(string), rs, err)
	return rs, err
}


/* ----------------------------内部方法------------------------------ */
func (this *Channel) mergeArgs(arrNeed []string, tmpArgs Arg) (Arg, error) {
	var arrArgs = make(Arg)
	var idx int = 0
	if arrNeed == nil && tmpArgs == nil {
		this.Err.Msg = "no params are set"
		this.Err.Code = CHANNEL_SDK_PARAM
		goto Failed
	}
	if len(tmpArgs)-1 != len(arrNeed) && len(tmpArgs) != len(arrNeed) {
		keys := "(" + strings.Join(arrNeed, ",") + ")"
		this.Err.Msg = "invalid sdk, params, params" + keys + "are need"
		this.Err.Code = CHANNEL_SDK_PARAM
		goto Failed
	}
	if _, ok := tmpArgs[len(tmpArgs)-1].(Opt); !ok && len(tmpArgs)-1 == len(arrNeed) {
		this.Err.Msg = "invalid sdk params, optional param must be an array"
		this.Err.Code = CHANNEL_SDK_PARAM
		goto Failed
	}

	if len(arrNeed) > 0 {
		for _, v := range arrNeed {
			if tmpArgs[idx] == nil {
				this.Err.Msg = "lack param " + v
				this.Err.Code = CHANNEL_SDK_PARAM
				goto Failed
			}
			arrArgs[v] = tmpArgs[idx]
			idx += 1
		}
	}
	if len(tmpArgs) == idx + 1 && tmpArgs[idx] != nil {
		for k, v := range tmpArgs[idx].(Opt) {
			if _,ok := arrArgs[k]; !ok && len(v) > 0 {
				arrArgs[k] = v
			}
		}
	}
	return arrArgs, nil

Failed:
	return nil, this.Err
}

func (this *Channel) commonProccess(paramOpt Arg) (Res, error) {
	this.adjustOpt(paramOpt)
	ret := this.baseControl(paramOpt)
	if len(ret.Body) == 0 {
		this.Err.Msg = "base control returned None object"
		this.Err.Code = CHANNEL_SDK_SYS
		return nil, this.Err
	}

	js, err := json.NewJson([]byte(ret.Body))
	if err != nil {
		if ret.IsOK() {
			this.Err.Code = CHANNEL_SDK_HTTP_STATUS_OK_BUT_RESULT_ERROR
		} else {
			this.Err.Msg = "ret body:" + ret.Body + err.Error()
			this.Err.Code = CHANNEL_SDK_HTTP_STATUS_ERROR_AND_RESULT_ERROR
		}
		return nil, this.Err
	}
	this.RequestId = js.Get("request_id").MustInt()
	result, err := js.Map()
	if err != nil {
		this.Err.Msg = err.Error()
		return nil, this.Err
	} else {
		code := js.Get("error_code").MustInt()
		msg := js.Get("error_msg").MustString()
		if code > 0 || msg != "" {
			this.Err.Code = code
			this.Err.Msg = msg
			return nil, this.Err
		}
	}
	return result, nil
}

func (this *Channel) adjustOpt(opt Arg) {
	if _,ok := opt[TIMESTAMP]; !ok {
		opt[TIMESTAMP] = time.Now().Unix()
	}
	opt[API_KEY] = this.ApiKey
	delete(opt, SECRET_KEY)
}

func (this *Channel) baseControl(opt Arg) ResponseCore {
	resource := "channel"
	if _,ok := opt[CHANNEL_ID]; ok {
		if m,ok := opt[METHOD]; ok && !methodChannelInBody.InArray(m.(string)) && opt[CHANNEL_ID] != nil {
			resource = opt[CHANNEL_ID].(string)
			delete(opt, CHANNEL_ID)
		}
	}
	url := BASEURL + resource
	opt[SIGN] = this.genSign(HTTP_METHOD, url, opt)

	req := RequestCore{
		RequestUrl: url,
		Method: HTTP_METHOD,
		UserAgent: USERAGENT,
	}
	req.InitHeader()
	for k,v := range headersMap {req.AddHeader(k, v)}
	var q = make(u.Values)
	for k,v := range opt {
		kk := fmt.Sprintf("%v", k)
		vv := fmt.Sprintf("%v", v)
		q.Add(kk, vv)
	}
	req.SetBody(q.Encode())
	req.SetCurlOpts(this.CurlOpts)
	req.HandleRequest()
	return ResponseCore{
		req.GetResponseHeader(),
		req.GetResponseBody(),
		req.GetResponseCode(),
	}
}

func (this *Channel) genSign(method,url string, opt Arg) string {
	var keys SortT
	for k,_ := range opt {if k != nil {keys = append(keys, k)}}
	sort.Sort(keys)
	gather := method + url
	for _,v := range keys {gather += fmt.Sprintf("%v=%v", v, opt[v])}
	gather += this.SecretKey
	h := md5.New()
	io.WriteString(h, u.QueryEscape(gather))
	sign := h.Sum(nil)
	return fmt.Sprintf("%x", sign)
}

func (this *Channel) output(method string, res Res, err error) {
	if DEBUG {
		if err == nil {
			fmt.Printf("\033[1;40;32mSUCC, %v OK!\nRESULT: %v\033[0m\n\n", method, res)
		} else {
			fmt.Printf("\033[1;40;31mWRONG, %v ERROR!\nERROR NUMBER: %v\nERROR MESSAGE: %v\nREQUEST ID: %v\033[0m\n\n",
				method, this.Err.Code, this.Err.Error(), this.RequestId)
		}
	}
}
