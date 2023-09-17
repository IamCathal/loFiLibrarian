package dtos

import "time"

type ctxKey int

const (
	REQUEST_ID ctxKey = iota
	BOOK_ID    ctxKey = iota
	TIME_TAKEN ctxKey = iota
	START_TIME ctxKey = iota
	WS         ctxKey = iota
)

type AppConfig struct {
	ApplicationStartUpTime time.Time
}

type UptimeResponse struct {
	Status      string        `json:"status,omitempty"`
	Uptime      time.Duration `json:"uptime,omitempty"`
	StartUpTime int64         `json:"startuptime,omitempty"`
}

type LookupRequest struct {
	ID     string `json:"id"`
	BookId string `json:"bookId"`
}

// Websocket message DTOs

type wsMessage struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`

	Timestamp int64  `json:"time"`
	ID        string `json:"id"`
	BookId    string `json:"bookId"`
}

type wsBookInfo struct {
	Type              string         `json:"type"`
	BookInfo          BookBreadcrumb `json:"bookInfo"`
	IsFromOpenLibrary bool           `json:"isFromOpenLibrary"`

	Timestamp int64  `json:"time"`
	ID        string `json:"id"`
	BookId    string `json:"bookId"`
	TimeTaken int64  `json:"timeTaken"`
}

type wsError struct {
	Type         string `json:"type"`
	ErrorMessage string `json:"errorMessage"`

	Timestamp int64  `json:"time"`
	ID        string `json:"id"`
	BookId    string `json:"bookId"`
	TimeTaken int64  `json:"timeTaken"`
}

type WsLiveStatus struct {
	Type              string `json:"type"`
	ServerSentTime    int64  `json:"serverSentTime"`
	ServerStartupTime int64  `json:"serverStartupTime,omitempty"`
}

// Optional, only used when rabbitMQ is enabled
type MorpheusEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Msg       string `json:"data"`
}
