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

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func WriteWsMessage(ctx context.Context, msg string) {
	ws := ctx.Value("ws").(*websocket.Conn)
	wsMessage := dtos.WsMessage{
		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value("requestId").(string),
		BookId:    ctx.Value("bookId").(string),
		Msg:       msg,
	}
	if err := ws.WriteJSON(wsMessage); err != nil {
		logger.Sugar().Panicf("failed to write msg '%s' to websocket: %+v", msg, err)
	}
	logger.Sugar().Infof("Write partial ws message: %+v", wsMessage)
}

func WriteWsPartialBookInfo(ctx context.Context, bookInfo dtos.BookBreadcrumb) {
	ws := ctx.Value("ws").(*websocket.Conn)
	wsMessage := dtos.WsBookInfo{
		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value("requestId").(string),
		TimeTaken: ctx.Value("timeTaken").(int64),
		BookInfo:  bookInfo,
	}
	if err := ws.WriteJSON(wsMessage); err != nil {
		logger.Sugar().Panicf("failed to write bookInfo for '%s' to websocket: %+v", bookInfo.Title, err)
	}
	logger.Sugar().Infof("Write partial book info message: %+v", wsMessage)
}

func WriteWsError(ctx context.Context, message string) {
	ws := ctx.Value("ws").(*websocket.Conn)
	wsMessage := dtos.WsError{
		Timestamp:    time.Now().UnixMilli(),
		ID:           ctx.Value("requestId").(string),
		ErrorMessage: message,
	}
	if err := ws.WriteJSON(wsMessage); err != nil {
		logger.Sugar().Panicf("failed to write errorMessage '%s' to websocket: %+v", message, err)
	}
	logger.Sugar().Infof("Write error ws message: %+v", wsMessage)
}
