package dtos

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// setup

	code := m.Run()

	os.Exit(code)
}

func TestNewMorpheusEvent(t *testing.T) {
	expectedMessage := "event messages"
	expectedLevel := "INFO"

	actualMorpheusEvent := NewMorpheusEvent(expectedMessage, expectedLevel)

	assert.Equal(t, expectedMessage, actualMorpheusEvent.Msg)
	assert.Equal(t, expectedLevel, actualMorpheusEvent.Level)
	assert.Equal(t, "lofilibrarian", actualMorpheusEvent.Type)
	assert.NotNil(t, actualMorpheusEvent.Timestamp)
}

func TestNewWsLiveStatus(t *testing.T) {
	expectedAppStartTime := time.Now().Add(-time.Minute * 5)

	actualWsLiveStatus := NewWsLiveStatus(expectedAppStartTime)

	assert.Equal(t, WS_LIVE_STATUS_WS_MESSAGE_TYPE, actualWsLiveStatus.Type)
	assert.Equal(t, expectedAppStartTime.UnixMilli(), actualWsLiveStatus.ServerStartupTime)
	assert.Less(t, expectedAppStartTime.UnixMilli(), actualWsLiveStatus.ServerSentTime)
}

func TestWsNewError(t *testing.T) {
	expectedBookId := "bookId"
	expectedTimeTaken := int64(25)
	expectedRequestId := "requestId"
	expectedMessage := "unspeakable things have occured"

	expectedContext := createBookDetailsContext(expectedRequestId, expectedBookId, expectedTimeTaken)

	actualWsError := NewWsError(expectedContext, expectedMessage)

	assert.Equal(t, ERROR_WS_MESSAGE_TYPE, actualWsError.Type)
	assert.Equal(t, expectedMessage, actualWsError.ErrorMessage)
	assert.NotNil(t, actualWsError.Timestamp)
	assert.Equal(t, expectedRequestId, actualWsError.ID)
	assert.Equal(t, expectedBookId, actualWsError.BookId)
	assert.LessOrEqual(t, expectedTimeTaken, actualWsError.TimeTaken)
}

func TestNewWsBookInfo_fromGoodreads(t *testing.T) {
	expectedBookId := "bookId"
	expectedTimeTaken := int64(25)
	expectedRequestId := "requestId"
	expectedContext := createBookDetailsContext(expectedRequestId, expectedBookId, expectedTimeTaken)

	expectedBookInfo := BookBreadcrumb{
		Title:  "Citizen Quinn",
		Author: "Gavin Daly & Ian Kahoe",
	}

	actualWsBookInfo := NewWsBookInfo(expectedContext, expectedBookInfo, false)

	assert.Equal(t, BOOK_INFO_WS_MESSAGE_TYPE, actualWsBookInfo.Type)
	assert.Equal(t, expectedBookInfo, actualWsBookInfo.BookInfo)
	assert.False(t, actualWsBookInfo.IsFromOpenLibrary)

	assert.NotNil(t, actualWsBookInfo.Timestamp)
	assert.Equal(t, expectedRequestId, actualWsBookInfo.ID)
	assert.Equal(t, expectedBookId, actualWsBookInfo.BookId)
	assert.LessOrEqual(t, expectedTimeTaken, actualWsBookInfo.TimeTaken)
}

func TestNewWsBookInfo_fromOpenLibrary(t *testing.T) {
	expectedBookId := "bookId"
	expectedTimeTaken := int64(25)
	expectedRequestId := "requestId"
	expectedContext := createBookDetailsContext(expectedRequestId, expectedBookId, expectedTimeTaken)

	expectedBookInfo := BookBreadcrumb{
		Title:  "Oblivion: Stories",
		Author: "David Foster Wallace",
	}

	actualWsBookInfo := NewWsBookInfo(expectedContext, expectedBookInfo, true)

	assert.Equal(t, BOOK_INFO_WS_MESSAGE_TYPE, actualWsBookInfo.Type)
	assert.Equal(t, expectedBookInfo, actualWsBookInfo.BookInfo)
	assert.True(t, actualWsBookInfo.IsFromOpenLibrary)

	assert.NotNil(t, actualWsBookInfo.Timestamp)
	assert.Equal(t, expectedRequestId, actualWsBookInfo.ID)
	assert.Equal(t, expectedBookId, actualWsBookInfo.BookId)
	assert.LessOrEqual(t, expectedTimeTaken, actualWsBookInfo.TimeTaken)
}

func TestNewWsMessage(t *testing.T) {
	expectedBookId := "bookId"
	expectedTimeTaken := int64(25)
	expectedMessage := "Few pints of plain going there"
	expectedContext := createBookDetailsContext("", expectedBookId, expectedTimeTaken)

	actualWsMessage := NewWsMessage(expectedContext, expectedMessage)

	assert.Equal(t, MGS_WS_MESSAGE_TYPE, actualWsMessage.Type)
	assert.Equal(t, expectedMessage, actualWsMessage.Msg)
	assert.NotNil(t, actualWsMessage.Timestamp)
	assert.Equal(t, expectedBookId, actualWsMessage.BookId)
}

func createBookDetailsContext(requestId, bookId string, timeTaken int64) context.Context {
	ctx := context.Background()

	ctx = context.WithValue(ctx, START_TIME, time.Now().Add(-time.Second*5).UnixMilli())
	ctx = context.WithValue(ctx, BOOK_ID, bookId)
	ctx = context.WithValue(ctx, REQUEST_ID, requestId)
	return ctx
}
