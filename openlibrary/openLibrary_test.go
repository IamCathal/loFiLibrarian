package openlibrary

import (
	"fmt"
	"os"
	"testing"

	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// setup

	code := m.Run()

	os.Exit(code)
}

func TestErrorWithTrace(t *testing.T) {
	baseErrorMessageText := "Run out of cans and its past 22:00"
	testErr := fmt.Errorf(baseErrorMessageText)

	errWithTrace := errWithTrace(testErr)

	assert.Contains(t, errWithTrace.Error(), testErr.Error())

	assert.Contains(t, errWithTrace.Error(), "github.com/iamcathal/lofilibrarian/openlibrary.errWithTrace")
	assert.Contains(t, errWithTrace.Error(), "goroutine")
}

func TestOpenLibBookToBreadcrumb(t *testing.T) {
	expectedBookBreadCrumb := dtos.BookBreadcrumb{
		Title:        "Boulevard Wren and Other Stories",
		Author:       "Blindboy Boatclub",
		Series:       "",
		MainCover:    "",
		Pages:        304,
		Rating:       0,
		RatingsCount: 0,
		Genres: []string{
			"Short Stories",
			"Irish Literature",
			"Craic agus Ceoil",
		},
		ISBN: "bookIsbn",
	}

	openLibraryBook := dtos.OpenLibraryBook{
		Title: expectedBookBreadCrumb.Title,
		Authors: []dtos.Authors{{
			Key: expectedBookBreadCrumb.Author,
		}},
		NumberOfPages: expectedBookBreadCrumb.Pages,
		Subjects:      expectedBookBreadCrumb.Genres,
	}

	actualBreadCrumb := openLibBookToBreadcrumb(expectedBookBreadCrumb.ISBN, openLibraryBook)

	assert.Equal(t, expectedBookBreadCrumb, actualBreadCrumb)
}
