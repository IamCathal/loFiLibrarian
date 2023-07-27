package util

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamcathal/lofilibrarian/dtos"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

type ctxKey int

const (
	REQUEST_ID ctxKey = iota
	BOOK_ID    ctxKey = iota
	TIME_TAKEN ctxKey = iota
	WS         ctxKey = iota
)

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func WriteWsMessage(ctx context.Context, msg string) {
	ws := ctx.Value(WS).(*websocket.Conn)
	wsMessage := dtos.WsMessage{
		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value(REQUEST_ID).(string),
		BookId:    ctx.Value(BOOK_ID).(string),
		Msg:       msg,
	}
	if err := ws.WriteJSON(wsMessage); err != nil {
		logger.Sugar().Panicf("failed to write msg '%s' to websocket: %+v", msg, err)
	}
	logger.Sugar().Infof("Write partial ws message: %+v", wsMessage)
}

func WriteBookDetailsBreadcrumb(ctx context.Context, bookBreadcrumb dtos.BookBreadcrumb) {
	ws := ctx.Value(WS).(*websocket.Conn)
	wsMessage := dtos.WsBookInfo{
		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value(REQUEST_ID).(string),
		TimeTaken: ctx.Value(TIME_TAKEN).(int64),
		BookInfo:  bookBreadcrumb,
	}
	if err := ws.WriteJSON(wsMessage); err != nil {
		logger.Sugar().Panicf("failed to write bookBreadcrumb for '%s' to websocket: %+v", bookBreadcrumb.ISBN, err)
	}
}

func WriteWsError(ctx context.Context, message string) {
	ws := ctx.Value(WS).(*websocket.Conn)
	wsMessage := dtos.WsError{
		Timestamp:    time.Now().UnixMilli(),
		ID:           ctx.Value(REQUEST_ID).(string),
		ErrorMessage: message,
	}
	if err := ws.WriteJSON(wsMessage); err != nil {
		logger.Sugar().Panicf("failed to write errorMessage '%s' to websocket: %+v", message, err)
	}
	logger.Sugar().Infof("Write error ws message: %+v", wsMessage)
}
