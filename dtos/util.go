package dtos

import (
	"context"
	"time"
)

var (
	MGS_WS_MESSAGE_TYPE            = "message"
	ERROR_WS_MESSAGE_TYPE          = "error"
	BOOK_INFO_WS_MESSAGE_TYPE      = "bookInfo"
	WS_LIVE_STATUS_WS_MESSAGE_TYPE = "liveStatus"
	MORPHEUS_EVENT_TYPE            = "lofilibrarian"
)

func NewWsMessage(ctx context.Context, msg string) wsMessage {
	return wsMessage{
		Type: MGS_WS_MESSAGE_TYPE,
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
		Type:              BOOK_INFO_WS_MESSAGE_TYPE,
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
		Type:         ERROR_WS_MESSAGE_TYPE,
		ErrorMessage: msg,

		Timestamp: time.Now().UnixMilli(),
		ID:        ctx.Value(REQUEST_ID).(string),
		BookId:    ctx.Value(BOOK_ID).(string),
		TimeTaken: ctx.Value(TIME_TAKEN).(int64),
	}
}

func NewWsLiveStatus(appStartTime time.Time) WsLiveStatus {
	return WsLiveStatus{
		Type:              WS_LIVE_STATUS_WS_MESSAGE_TYPE,
		ServerSentTime:    time.Now().UnixMilli(),
		ServerStartupTime: appStartTime.UnixMilli(),
	}
}

func NewMorpheusEvent(msg, level string) MorpheusEvent {
	return MorpheusEvent{
		//ID:
		Timestamp: time.Now().UnixMilli(),
		Type:      MORPHEUS_EVENT_TYPE,
		Level:     level,
		Msg:       msg,
	}
}
