package dtos

import (
	"context"
	"time"
)

func NewWsMessage(ctx context.Context, msg string) wsMessage {
	return wsMessage{
		Type: "message",
		Msg:  msg,

		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value(REQUEST_ID).(string),
		BookId:    ctx.Value(BOOK_ID).(string),
	}
}

func NewWsBookInfo(ctx context.Context, bookInfo BookBreadcrumb) wsBookInfo {
	startTime := ctx.Value(START_TIME).(int64)
	ctx = context.WithValue(ctx, TIME_TAKEN, time.Now().UnixMilli()-startTime)

	return wsBookInfo{
		Type:     "bookInfo",
		BookInfo: bookInfo,

		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value(REQUEST_ID).(string),
		BookId:    ctx.Value(BOOK_ID).(string),
		TimeTaken: ctx.Value(TIME_TAKEN).(int64),
	}
}

func NewWsError(ctx context.Context, msg string) wsError {
	startTime := ctx.Value(START_TIME).(int64)
	ctx = context.WithValue(ctx, TIME_TAKEN, time.Now().UnixMilli()-startTime)

	return wsError{
		Type:         "error",
		ErrorMessage: msg,

		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value(REQUEST_ID).(string),
		BookId:    ctx.Value(BOOK_ID).(string),
		TimeTaken: ctx.Value(TIME_TAKEN).(int64),
	}
}
