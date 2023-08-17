package dtos

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
