package openlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/iamcathal/lofilibrarian/dtos"
	"go.uber.org/zap"
)

const (
	ISBN_SEARCH_BASE_URL = "https://openlibrary.org/isbn"
)

var (
	logger *zap.Logger
)

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func IsbnSearch(ctx context.Context, isbn string) (dtos.OpenLibraryBook, error) {
	logger.Sugar().Infof("Retrieving openLibrary isbn search for bookId: %s", ctx.Value(dtos.BOOK_ID).(string))
	attemptsMade := 0
	maxRetryCount := 3

	errorsEncountered := []error{}

	bookInfo := dtos.OpenLibraryBook{}

	for {
		if attemptsMade >= maxRetryCount {
			break
		}
		attemptsMade++
		logger.Sugar().Infof("Attempt %d to search openLibrary for bookId: %s", attemptsMade, isbn)

		thePage, getPageErr := makeIsbnSearch(ctx, isbn)
		if getPageErr != nil {
			errorsEncountered = append(errorsEncountered, fmt.Errorf("failed to search openlibrary for bookId: %s: %w", isbn, getPageErr))
			thePage.Close()
			time.Sleep(1 * time.Second)
			continue
		}
		defer thePage.Close()

		decodeErr := json.NewDecoder(thePage).Decode(&bookInfo)
		if decodeErr == nil {
			logger.Sugar().Infof("Retrieved details for bookId: %s (%s)", isbn, bookInfo.Title)
			return bookInfo, nil
		}

		errorsEncountered = append(errorsEncountered, fmt.Errorf("failed to unmarshal openLibrary json response: %w", decodeErr))
	}

	logger.Sugar().Warnf("Failed to search openLibrary for bookId: %s after %d retries: %+v", isbn, attemptsMade, errorsEncountered)
	return dtos.OpenLibraryBook{}, fmt.Errorf("failed to search openLibrary for bookId: %s after %d retries: %+v", isbn, attemptsMade, errorsEncountered)
}

func makeIsbnSearch(ctx context.Context, isbn string) (io.ReadCloser, error) {
	fullSearchUrl := fmt.Sprintf("%s/%s.json", ISBN_SEARCH_BASE_URL, isbn)
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullSearchUrl, nil)
	if err != nil {
		myErr := fmt.Errorf("failed to create GET request for URL ' %s ': %w", fullSearchUrl, errWithTrace(err))
		return nil, myErr
	}
	req.Header.Set("User-Agent", "LofiLibrarian :)")

	res, err := client.Do(req)
	if err != nil {
		myErr := fmt.Errorf("failed to get URL ' %s ': %w", fullSearchUrl, errWithTrace(err))
		return nil, myErr
	}
	return res.Body, nil
}

func errWithTrace(err error) error {
	return fmt.Errorf(err.Error(), string(debug.Stack()))
}
