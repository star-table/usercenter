package logger

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/client"
)

// log wrapper logs every time a request is made
type LogWrapper struct {
	client.Client
}

func (l *LogWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[wrapper] client request service: %s method: %s\n", req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp)
}

// Implements client.Wrapper as logWrapper
func LogWrap(c client.Client) client.Client {
	return &LogWrapper{c}
}
