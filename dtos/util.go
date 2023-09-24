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

func NewWsBookInfo(ctx context.Context, bookInfo BookBreadcrumb, isFromOpenLibrary bool) wsBookInfo {
	startTime := ctx.Value(START_TIME).(int64)
	ctx = context.WithValue(ctx, TIME_TAKEN, time.Now().UnixMilli()-startTime)

	return wsBookInfo{
		Type:              "bookInfo",
		BookInfo:          bookInfo,
		IsFromOpenLibrary: isFromOpenLibrary,

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

func NewWsLiveStatus(appStartTime time.Time) WsLiveStatus {
	return WsLiveStatus{
		Type:              "liveStatus",
		ServerSentTime:    time.Now().UnixMilli(),
		ServerStartupTime: appStartTime.UnixMilli(),
	}
}

func NewMorpheusEvent(msg, level string) MorpheusEvent {
	return MorpheusEvent{
		//ID:
		Timestamp: time.Now().UnixMilli(),
		Type:      "lofilibrarian",
		Level:     level,
		Msg:       msg,
	}
}
