/**
 * 百度云推送SDK，异常处理
 * User: zha_ming@163.com
 * Package: baidupush
 * Date: 14-6-17 14:23
 * Version: 0.1
 */
package baidupush

import (
	"os"
	"fmt"
	"time"
)

const (
	LAYOUT = "2006-01-02 15:04:05"
	PRODUCTED = false
)

type PushError struct {
	Msg string
	Code int
	Map []string
}

func (this PushError) Error() string {
	if this.Code >= 0 && len(this.Map) > 0 && this.Code < len(this.Map) {
		return fmt.Sprintf("[%v] %v %v", time.Now().Format(LAYOUT), this.Msg, this.Map[this.Code])
	}
	return fmt.Sprintf("[%v] %v", time.Now().Format(LAYOUT), this.Msg)
}

func NewError(_map []string) PushError {
	return PushError{Map: _map}
}

func pError(s string) {
	if PRODUCTED {fmt.Printf(s);os.Exit(1)}
}
