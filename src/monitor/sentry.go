package monitor

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

var client *raven.Client

func init() {
	client, _ = raven.New("http://63302e6e5dc742f5888732dcf7e24019:b163737abe194898846c9b343f9a3028@sentry.xcodecraft.cn/2")
}

func Report(flags map[string]string, err interface{}, errs []*gin.Error) {
	errStr := fmt.Sprint(err)
	packet := raven.NewPacket(errStr, raven.NewException(errors.New(errStr), raven.NewStacktrace(2, 3, nil)))
	client.Capture(packet, flags)
	_, ch := client.Capture(packet, flags)
	if err = <-ch; err != nil {
		fmt.Println("report sentry error")
	}
	for _, item := range errs {
		packet := raven.NewPacket((*item).Error(), &raven.Message{(*item).Error(), []interface{}{item.Meta}})
		client.Capture(packet, flags)
		_, ch := client.Capture(packet, flags)
		if err = <-ch; err != nil {
			fmt.Println("report sentry error")
		}
	}
}
