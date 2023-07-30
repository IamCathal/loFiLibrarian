package goodreads

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iamcathal/lofilibrarian/dtos"
)

var (
	// There are five spaces between a books
	// title and its series information if
	// the series information is given
	TITLE_AND_SERIES_INFO_SEPERATOR = regexp.MustCompile("[ ]{3,}")
	// Goodreads returns 30 books per page
	BOOK_COUNT_PER_PAGE = 30
	// Base URL that book links are built on
	GOODREADS_BASE_BOOK_URL = "https://www.goodreads.com"
	// Crude to check if a roughly  valid
	// shelf URL is being queried
	GOODREADS_SHELF_URL_PREFIX = GOODREADS_BASE_BOOK_URL + "/review/list/"
	ONLY_NUMBERS               = regexp.MustCompile(`([0-9]+)`)
)

func checkErr(err error) {
	if err != nil {
		logger.Sugar().Fatal(err)
	}
}

func stripOfFormatting(input string) string {
	formatted := strings.ReplaceAll(input, "\n", "")
	formatted = strings.TrimSpace(formatted)
	return formatted
}

func strToInt(str string) (int, error) {
	intVersion, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("failed parse int from string '%s': %w", str, err)
	}
	return intVersion, nil
}

func strToFloat(floatString string) (float64, error) {
	floatVal, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		return 0, fmt.Errorf("failed parse float from string '%s': %w", floatString, err)
	}
	return floatVal, nil
}

func getRatingsCount(ratingsString string) (int, error) {
	extractedNumbers := ONLY_NUMBERS.FindAllString(ratingsString, 2)
	return strToInt(extractedNumbers[0])
}

func getFakeReferrerPage(URL string) string {
	splitStringByShelfParam := strings.Split(URL, "?")
	return splitStringByShelfParam[0]
}

func extractIntPages(pagesString string) (int, error) {
	extractedNumbers := ONLY_NUMBERS.FindAllString(pagesString, 2)
	if len(extractedNumbers) != 1 {
		return 0, fmt.Errorf("failed to extract one numbers from pagesString '%s'", pagesString)
	}
	return strToInt(extractedNumbers[0])
}

func extractGenres(doc *goquery.Document) []string {
	genres := []string{}

	doc.Find("span.BookPageMetadataSection__genreButton").Each(func(i int, genre *goquery.Selection) {
		if genre.Text() != "Genres" {
			genres = append(genres, genre.Text())
		}
	})

	noDuplicates := removeDuplicateGenres(genres)
	if len(noDuplicates) >= 6 {
		return noDuplicates[:6]
	}
	return noDuplicates
}

func removeDuplicateGenres(genres []string) []string {
	seenGenres := make(map[string]bool)
	noDuplicatedGenres := []string{}

	for _, genre := range genres {
		_, exists := seenGenres[genre]
		if !exists {
			seenGenres[genre] = true
			noDuplicatedGenres = append(noDuplicatedGenres, genre)
		}
	}
	return noDuplicatedGenres
}

func extractOtherCovers(doc *goquery.Document) []string {
	otherEditionCovers := []string{}

	doc.Find("div[class='otherEditionCovers']").Each(func(i int, otherEditionDiv *goquery.Selection) {
		otherEditionDiv.Find("img").Each(func(i int, otherEditionImg *goquery.Selection) {
			otherEditionCoverImg, exists := otherEditionImg.Attr("src")
			if exists {
				otherEditionCovers = append(otherEditionCovers, otherEditionCoverImg)
			}
		})
	})

	logger.Sugar().Infof("Found %d other covers %v", len(otherEditionCovers), otherEditionCovers)
	return otherEditionCovers
}

func getConciseBookInfoFromBreadCrumb(breadCrumb dtos.BookBreadcrumb) string {
	return fmt.Sprintf("%s by %s (%s) %.2f stars and has genres %+v",
		breadCrumb.Title, breadCrumb.Author, breadCrumb.Series,
		breadCrumb.Rating, breadCrumb.Genres)
}
