package goodreads

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/util"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func GetBookDetailsWs(ctx context.Context, ID string) (dtos.BookBreadcrumb, error) {
	body := getPage(fmt.Sprintf("https://www.goodreads.com/book/auto_complete?format=json&q=%s", ID))
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return dtos.BookBreadcrumb{}, err
	}

	booksFoundRes := []dtos.GoodReadsSearchBookResult{}
	err = json.Unmarshal(bodyBytes, &booksFoundRes)
	if err != nil {
		return dtos.BookBreadcrumb{}, err
	}

	logger.Sugar().Infof("%d books were found for ID: %s", len(booksFoundRes), ID)
	if len(booksFoundRes) == 0 {
		return dtos.BookBreadcrumb{}, nil
	}

	book := booksFoundRes[0]
	floatRating, err := strToFloat(book.AvgRating)
	if err != nil || len(booksFoundRes) == 0 {
		return dtos.BookBreadcrumb{}, err
	}
	partialBookBreadcrumb := dtos.BookBreadcrumb{
		Title:        book.BookTitleBare,
		Author:       book.Author.Name,
		Series:       "",
		MainCover:    book.ImageURL,
		OtherCovers:  []string{},
		Pages:        book.NumPages,
		Link:         book.Description.FullContentURL,
		Rating:       floatRating,
		RatingsCount: book.RatingsCount,
		Genres:       []string{},
		ISBN:         ID,
	}

	util.WriteBookDetailsBreadcrumb(ctx, partialBookBreadcrumb)

	return lookUpGoodReadsPageForBook(ctx, book.Description.FullContentURL)
}

func lookUpGoodReadsPageForBook(ctx context.Context, bookPageURL string) (dtos.BookBreadcrumb, error) {
	ctx = context.WithValue(ctx, dtos.START_TIME, time.Now().UnixMilli())

	fullBookInfo, err := extractBookInfo(ctx, bookPageURL)
	if err != nil {
		return dtos.BookBreadcrumb{}, err
	}
	fullBookInfo.ISBN = ctx.Value(dtos.BOOK_ID).(string)

	util.WriteBookDetailsBreadcrumb(ctx, fullBookInfo)

	return fullBookInfo, nil
}

// func GetBookDetails(ID string) dtos.BookBreadcrumb {
// 	logger.Sugar().Infof("Retrieving book details for ID: %s", ID)
// 	body := getPage(fmt.Sprintf("https://www.goodreads.com/book/auto_complete?format=json&q=%s", ID))
// 	bodyBtytes, err := io.ReadAll(body)
// 	checkErr(err)

// 	booksFoundRes := []dtos.GoodReadsSearchBookResult{}
// 	err = json.Unmarshal(bodyBtytes, &booksFoundRes)
// 	checkErr(err)

// 	if len(booksFoundRes) == 0 {
// 		logger.Sugar().Infof("No books found for ID: %s", ID)
// 		return dtos.BookBreadcrumb{}
// 	}
// 	logger.Sugar().Infof("%d books were found for ID: %s at %s", len(booksFoundRes), ID, booksFoundRes[0].Description.FullContentURL)
// 	return extractBookInfo(booksFoundRes[0].Description.FullContentURL)
// }

func extractBookInfo(ctx context.Context, bookPage string) (dtos.BookBreadcrumb, error) {
	logger.Sugar().Infof("Retrieving goodreads page for bookId: %s with URL: %s", ctx.Value(dtos.BOOK_ID).(string), bookPage)

	thePage := getPage(bookPage)
	defer thePage.Close()

	doc, err := goquery.NewDocumentFromReader(thePage)
	checkErr(err)

	bookInfo := dtos.BookBreadcrumb{}

	// TODO when an error is hit here (often cant read page numbers) the entire HTML document should be logged
	// for debugging to see if it was just an empty response. Maybe try a retry in a second or two
	bookInfo.Title = strings.TrimSpace(doc.Find("h1.Text").Text())
	bookInfo.Author = strings.TrimSpace(doc.Find(".ContributorLinksList > span:nth-child(1) > a:nth-child(1) > span:nth-child(1)").Text())
	bookInfo.Series = strings.TrimSpace(doc.Find("h3.Text__italic > a:nth-child(1)").Text())
	bookInfo.MainCover, _ = doc.Find("div.BookCard__clickCardTarget > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > img:nth-child(1)").Attr("src")
	ratingsCountStr := doc.Find("a.RatingStatistics > div:nth-child(2) > div:nth-child(1) > span:nth-child(1)").Text()
	bookInfo.OtherCovers = extractOtherCovers(doc)
	bookInfo.Genres = extractGenres(doc)

	bookInfo.Pages, err = extractIntPages(strings.TrimSpace(doc.Find(".FeaturedDetails > p:nth-child(1)").Text()))
	if err != nil {
		return dtos.BookBreadcrumb{}, err
	}

	bookInfo.Rating, err = strToFloat(stripOfFormatting(doc.Find("a.RatingStatistics > div:nth-child(1) > div:nth-child(2)").Text()))
	if err != nil {
		logger.Sugar().Infof("Current bookInfo before rating value extraction failure: %+v", bookInfo)
		return dtos.BookBreadcrumb{}, err
	}
	bookInfo.RatingsCount, err = getRatingsCount(ratingsCountStr)
	if err != nil {
		logger.Sugar().Infof("Current bookInfo before ratings count extraction failure: %+v", bookInfo)
		return dtos.BookBreadcrumb{}, err
	}

	logger.Sugar().Infof("Extracted all details for bookId %s: %+v", ctx.Value(dtos.BOOK_ID).(string), getConciseBookInfoFromBreadCrumb(bookInfo))
	return bookInfo, err
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
