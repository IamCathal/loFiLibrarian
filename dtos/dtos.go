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

// Goodreads

type BookBreadcrumb struct {
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Series       string   `json:"series"`
	MainCover    string   `json:"mainCover"`
	OtherCovers  []string `json:"otherCovers"`
	Pages        int      `json:"pages"`
	Link         string   `json:"link"`
	Rating       float64  `json:"rating"`
	RatingsCount int      `json:"ratingsCount"`
	Genres       []string `json:"genres"`
	ISBN         string   `json:"isbn"`
}

type GoodReadsSearchBookResult struct {
	ImageURL      string               `json:"imageUrl"`
	BookID        string               `json:"bookId"`
	WorkID        string               `json:"workId"`
	BookURL       string               `json:"bookUrl"`
	FromSearch    bool                 `json:"from_search"`
	FromSrp       bool                 `json:"from_srp"`
	Qid           string               `json:"qid"`
	Rank          int                  `json:"rank"`
	Title         string               `json:"title"`
	BookTitleBare string               `json:"bookTitleBare"`
	NumPages      int                  `json:"numPages"`
	AvgRating     string               `json:"avgRating"`
	RatingsCount  int                  `json:"ratingsCount"`
	Author        GoodReadsAuthor      `json:"author"`
	KcrPreviewURL string               `json:"kcrPreviewUrl"`
	Description   GoodReadsDescription `json:"description"`
}

type GoodReadsAuthor struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	IsGoodreadsAuthor bool   `json:"isGoodreadsAuthor"`
	ProfileURL        string `json:"profileUrl"`
	WorksListURL      string `json:"worksListUrl"`
}

type GoodReadsDescription struct {
	HTML           string `json:"html"`
	Truncated      bool   `json:"truncated"`
	FullContentURL string `json:"fullContentUrl"`
}

type InitLookupDto struct {
	ID     string `json:"id"`
	BookId string `json:"bookId"`
}

// .
// Websocket message DTOs
// .

type wsMessage struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`

	Timestamp int64  `json:"time"`
	ID        string `json:"id"`
	BookId    string `json:"bookId"`
}

type wsBookInfo struct {
	Type     string         `json:"type"`
	BookInfo BookBreadcrumb `json:"bookInfo"`

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
