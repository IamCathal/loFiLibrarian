package util

import (
	"context"
	"errors"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamcathal/lofilibrarian/dtos"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func WriteWsMessage(ctx context.Context, msg string) {
	ws := ctx.Value(dtos.WS).(*websocket.Conn)
	wsMessage := dtos.NewWsMessage(ctx, msg)
	if err := ws.WriteJSON(wsMessage); err != nil && !isClosedErr(err) {
		logger.Sugar().Panicf("failed to write msg ' %s ' to websocket: %+v", msg, err)
	}
	logger.Sugar().Infof("Write partial ws message: %+v", wsMessage)
}

func WriteBookDetailsBreadcrumb(ctx context.Context, bookBreadcrumb dtos.BookBreadcrumb, isFromOpenLibrary bool) {
	if ctx.Value(dtos.WS) == nil {
		logger.Info("No websocket connection present, skipped write message")
		return
	}
	ws := ctx.Value(dtos.WS).(*websocket.Conn)

	wsBookInfo := dtos.NewWsBookInfo(ctx, bookBreadcrumb, isFromOpenLibrary)
	if err := ws.WriteJSON(wsBookInfo); err != nil && !isClosedErr(err) {
		logger.Sugar().Panicf("failed to write bookBreadcrumb for ' %s ' to websocket: %+v", bookBreadcrumb.ISBN, err)
	}
}

func WriteWsError(ctx context.Context, message string) {
	if ctx.Value(dtos.WS) == nil {
		logger.Info("No websocket connection present, skipped write message")
		return
	}
	ws := ctx.Value(dtos.WS).(*websocket.Conn)

	wsError := dtos.NewWsError(ctx, message)
	if err := ws.WriteJSON(wsError); err != nil && !isClosedErr(err) {
		logger.Sugar().Panicf("failed to write errorMessage ' %s ' to websocket: %+v", message, err)
	}

	logger.Sugar().Infof("Write error ws message: %+v", wsError)
}

func WriteWsLiveStatus(ctx context.Context, appStartTime time.Time) {
	if ctx.Value(dtos.WS) == nil {
		logger.Info("No websocket connection present, skipped write message")
		return
	}
	ws := ctx.Value(dtos.WS).(*websocket.Conn)

	wsLiveStatus := dtos.NewWsLiveStatus(appStartTime)
	if err := ws.WriteJSON(wsLiveStatus); err != nil && !isClosedErr(err) {
		logger.Sugar().Panicf("failed to write wsLiveStatus ' %s ' to websocket: %+v", wsLiveStatus, err)
	}

}

func isClosedErr(err error) bool {
	return errors.Is(err, syscall.EPIPE)
}
