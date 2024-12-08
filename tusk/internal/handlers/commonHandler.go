package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

type CommonHandler struct {
}

func NewCommonHandler() CommonHandler {
	return CommonHandler{}
}

func (ch *CommonHandler) Handle(handler func(hrw http.ResponseWriter, hreq *http.Request) error) http.Handler {

	operationName := GetFunctionName(handler)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		err := handler(rw, req)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				zap.L().Warn(fmt.Sprintf("%s: ctx canceled", operationName), zap.Error(err))
			} else {
				zap.L().Error(fmt.Sprintf("%s: failed", operationName), zap.Error(err))
			}

			SendError(rw, err)
		}
	})
}

func GetFunctionName(i interface{}) string {
	rv := reflect.ValueOf(i)
	fuc := runtime.FuncForPC(rv.Pointer())
	if fuc == nil {
		return "unknown"
	}
	name := fuc.Name()
	splits := strings.Split(name, ".")
	if len(splits) == 0 {
		return "unknown"
	}
	splits = splits[1:]
	name = strings.Join(splits, ".")
	name = strings.ReplaceAll(name, "(*", "")
	name = strings.ReplaceAll(name, ")", "")
	splits = strings.Split(name, "-")
	if len(splits) == 0 {
		return name
	}
	return splits[0]
}
