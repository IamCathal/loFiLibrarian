package goodreads

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
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
	body, err := getPage(fmt.Sprintf("https://www.goodreads.com/book/auto_complete?format=json&q=%s", ID))
	if err != nil {
		return dtos.BookBreadcrumb{}, err
	}
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

	possibleBookUrl := ""
	if book.Description.FullContentURL == "" {
		// try to use the bookUrl instead
		possibleBookUrl = GOODREADS_BASE_BOOK_URL + book.BookURL
	} else {
		possibleBookUrl = book.Description.FullContentURL
	}

	logger.Sugar().Infof("Possible book URL is %s", possibleBookUrl)

	partialBookBreadcrumb := dtos.BookBreadcrumb{
		Title:        book.BookTitleBare,
		Author:       book.Author.Name,
		Series:       "",
		MainCover:    book.ImageURL,
		OtherCovers:  []string{},
		Pages:        book.NumPages,
		Link:         possibleBookUrl,
		Rating:       floatRating,
		RatingsCount: book.RatingsCount,
		Genres:       []string{},
		ISBN:         ID,
	}
	util.WriteBookDetailsBreadcrumb(ctx, partialBookBreadcrumb)

	return lookUpGoodReadsPageForBook(ctx, possibleBookUrl)
}

func lookUpGoodReadsPageForBook(ctx context.Context, bookPageURL string) (dtos.BookBreadcrumb, error) {
	ctx = context.WithValue(ctx, dtos.START_TIME, time.Now().UnixMilli())

	fullBookInfo, err := getBookInfoFromURL(ctx, bookPageURL)
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

func getBookInfoFromURL(ctx context.Context, bookPageURL string) (dtos.BookBreadcrumb, error) {
	logger.Sugar().Infof("Retrieving goodreads page for bookId: %s with URL: %s", ctx.Value(dtos.BOOK_ID).(string), bookPageURL)
	attemptsMade := 0
	maxRetryCount := 3

	errorsEncountered := []error{}

	for {
		if attemptsMade >= maxRetryCount {
			break
		}
		attemptsMade++
		logger.Sugar().Infof("Attempt %d get book info from %s", attemptsMade, bookPageURL)

		thePage, getPageErr := getPage(bookPageURL)
		if getPageErr != nil {
			errorsEncountered = append(errorsEncountered, getPageErr)
			thePage.Close()
			continue
		}
		defer thePage.Close()

		doc, makeDocumentFromBodyErr := goquery.NewDocumentFromReader(thePage)
		if makeDocumentFromBodyErr != nil {
			errorsEncountered = append(errorsEncountered, makeDocumentFromBodyErr)
			thePage.Close()
		}

		bookInfo, extractDetailsErr := extractBookInfoFromResponse(ctx, doc)
		if extractDetailsErr != nil {
			errorsEncountered = append(errorsEncountered, extractDetailsErr)
			thePage.Close()
			continue
		} else {
			thePage.Close()
			logger.Sugar().Infof("Successfully extracted book info after %d attempts", attemptsMade)
			return bookInfo, nil
		}
	}

	logger.Sugar().Warnf("Failed to get book info from URL after %d retries: %+v", attemptsMade, errorsEncountered)
	return dtos.BookBreadcrumb{}, fmt.Errorf("failed to get book info from URL after %d retries: %+v", attemptsMade, errorsEncountered)
}

func extractBookInfoFromResponse(ctx context.Context, doc *goquery.Document) (dtos.BookBreadcrumb, error) {
	bookInfo := dtos.BookBreadcrumb{}
	var err error

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

func getPage(pageURL string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		myErr := fmt.Errorf("failed to create GET request for URL ' %s ': %w", pageURL, errWithTrace(err))
		return nil, myErr
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.goodreads.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", getFakeReferrerPage(pageURL))

	res, err := client.Do(req)
	if err != nil {
		myErr := fmt.Errorf("failed to get URL ' %s ': %w", pageURL, errWithTrace(err))
		return nil, myErr
	}

	return res.Body, nil
}

func errWithTrace(err error) error {
	return fmt.Errorf(err.Error(), string(debug.Stack()))
}
