package goodreads

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iamcathal/lofilibrarian/dtos"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func GetBookDetails(ID string) dtos.BookBreadcrumb {
	logger.Sugar().Infof("Retrieving book details for ID: %s", ID)
	body := getPage(fmt.Sprintf("https://www.goodreads.com/book/auto_complete?format=json&q=%s", ID))
	bodyBtytes, err := io.ReadAll(body)
	checkErr(err)

	booksFoundRes := []dtos.GoodReadsSearchBookResult{}
	err = json.Unmarshal(bodyBtytes, &booksFoundRes)
	checkErr(err)

	if len(booksFoundRes) == 0 {
		logger.Sugar().Infof("No books found for ID: %s", ID)
		return dtos.BookBreadcrumb{}
	}
	logger.Sugar().Infof("%d books were found for ID: %s at %s", len(booksFoundRes), ID, booksFoundRes[0].Description.FullContentURL)
	return extractBookInfo(booksFoundRes[0].Description.FullContentURL)
}

func extractBookInfo(bookPage string) dtos.BookBreadcrumb {
	doc, err := goquery.NewDocumentFromReader(getPage(bookPage))
	checkErr(err)

	bookInfo := dtos.BookBreadcrumb{}

	bookInfo.Title = strings.TrimSpace(doc.Find("h1.Text").Text())
	bookInfo.Author = strings.TrimSpace(doc.Find(".ContributorLinksList > span:nth-child(1) > a:nth-child(1) > span:nth-child(1)").Text())
	bookInfo.Series = strings.TrimSpace(doc.Find("h3.Text__italic > a:nth-child(1)").Text())
	bookInfo.MainCover, _ = doc.Find("div.BookCard__clickCardTarget > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > img:nth-child(1)").Attr("src")
	bookInfo.OtherCovers = extractOtherCovers(doc)
	bookInfo.Pages = extractIntPages(strings.TrimSpace(doc.Find(".FeaturedDetails > p:nth-child(1)").Text()))
	bookInfo.Link = bookPage
	bookInfo.Rating = strToFloat(stripOfFormatting(doc.Find("a.RatingStatistics > div:nth-child(1) > div:nth-child(2)").Text()))
	ratingsCountStr := doc.Find("a.RatingStatistics > div:nth-child(2) > div:nth-child(1) > span:nth-child(1)").Text()
	bookInfo.RatingsCount = getRatingsCount(ratingsCountStr)
	bookInfo.Genres = extractGenres(doc)

	logger.Sugar().Infof("Extracted all details for book URL: %s, bookInfo: %+v", bookPage, bookInfo)
	return bookInfo
}

func getPage(pageURL string) io.ReadCloser {
	client := &http.Client{}
	req, err := http.NewRequest("GET", pageURL, nil)
	checkErr(err)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.goodreads.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", getFakeReferrerPage(pageURL))

	res, err := client.Do(req)
	checkErr(err)
	return res.Body
}
