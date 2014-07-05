/**
 * 百度云推送SDK，CURL封装
 * User: zha_ming@163.com
 * Package: baidupush
 * Date: 14-6-17 11:41
 * Version: 0.1
 */
package baidupush

import (
	."time"
	"strings"
	"strconv"
	"bytes"
	curl "github.com/andelf/go-curl"
)

const (
	TIMEOUT = 30 * Second
	CONNECTIONTIMEOUT = 5 * Second
)


type RequestCore struct {
	RequestUrl		string
	Method			string
	UserAgent		string
	Debug			bool

	curlOpts		[]string
	requestHeaders	map[string]string
	requestBody		string
	responseHeaders	string
	responseBody	string
	responseCode	string
}

func (this *RequestCore) SetRequestUrl(url string) {
	this.RequestUrl = url
}

func (this *RequestCore) SetMethod(method string) {
	this.Method = strings.ToUpper(method)
}

func (this *RequestCore) SetUserAgent(ua string) {
	this.UserAgent = ua
}

func (this *RequestCore) SetCurlOpts(curlOpts []string) {
	this.curlOpts = curlOpts
}

func (this *RequestCore) InitHeader() {
	this.requestHeaders = make(map[string]string)
}

func (this *RequestCore) AddHeader(key,value string) {
	this.requestHeaders[key] = value
}

func (this *RequestCore) RemoveHeader(key string) {
	delete(this.requestHeaders, key)
}

func (this *RequestCore) SetBody(body string) {
	this.requestBody = body
}

func (this *RequestCore) GetResponseHeader() string {
	return this.responseHeaders
}

func (this *RequestCore) GetResponseBody() string {
	return this.responseBody
}

func (this *RequestCore) GetResponseCode() string {
	return this.responseCode
}

func (this *RequestCore) HandleRequest() {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	if this.Debug {
		easy.Setopt(curl.OPT_VERBOSE, true)
	}
	easy.Setopt(curl.OPT_URL, this.RequestUrl)
	easy.Setopt(curl.OPT_REFERER, this.RequestUrl)
	easy.Setopt(curl.OPT_USERAGENT, this.UserAgent)
	easy.Setopt(curl.OPT_TIMEOUT, int64(TIMEOUT))
	easy.Setopt(curl.OPT_CONNECTTIMEOUT, int64(CONNECTIONTIMEOUT))
	easy.Setopt(curl.OPT_FILETIME, true)
	easy.Setopt(curl.OPT_HEADER, true)
	easy.Setopt(curl.OPT_FOLLOWLOCATION, 1)
	easy.Setopt(curl.OPT_MAXREDIRS, 5)
	if len(this.requestHeaders) > 0 {
		var headers []string
		for key,value := range this.requestHeaders {
			headers = append(headers, key + ":" + value)
		}
		easy.Setopt(curl.OPT_HTTPHEADER, headers)
	}

	switch this.Method {
	case "POST":
		easy.Setopt(curl.OPT_POST, true)
		easy.Setopt(curl.OPT_POSTFIELDS, this.requestBody)
	case "HEAD":
		easy.Setopt(curl.OPT_CUSTOMREQUEST, this.Method)
		easy.Setopt(curl.OPT_NOBODY, 1)
	case "PUT":
		fallthrough
	default:  //"GET"
		easy.Setopt(curl.OPT_CUSTOMREQUEST, this.Method)
		easy.Setopt(curl.OPT_POSTFIELDS, this.requestBody)
	}

	var b []byte
	buffer := bytes.NewBuffer(b)
	writeFunc := func(buf []byte, _ interface{}) bool {
		_, err := buffer.Write(buf)
		if err != nil {
			return false
		}
		return true
	}
	easy.Setopt(curl.OPT_WRITEFUNCTION, writeFunc)

	if this.curlOpts != nil {
		for key,value := range this.curlOpts {
			easy.Setopt(key, value)
		}
	}

	if err := easy.Perform(); err != nil {
		pError("Perform:" + err.Error())
	}

	//Sleep(1 * Second)
	code, _ := easy.Getinfo(curl.INFO_HTTP_CODE)
	this.responseCode = strconv.Itoa(code.(int))
	headerSize, _ := easy.Getinfo(curl.INFO_HEADER_SIZE)
	size, _ := headerSize.(int)

	buf := make([]byte, buffer.Len())
	_, err := buffer.Read(buf)
	if err != nil {
		pError("Read:" + err.Error())
	}
	this.responseHeaders = string(buf[:size])
	this.responseBody = string(buf[size:])
}


type ResponseCore struct {
	Header,Body,Status string
}

func (this *ResponseCore) IsOK() bool {
	codes := ",200,201,204,206,"
	status := "," + this.Status + ","
	return strings.Contains(codes, status)
}
