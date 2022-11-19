package goodreads

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func strToInt(str string) int {
	intVersion, err := strconv.Atoi(str)
	checkErr(err)
	return intVersion
}

func strToFloat(floatString string) float64 {
	floatVal, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		logger.Sugar().Fatalf("failed to parse floatString: %s", floatString)
	}
	return floatVal
}

func getFakeReferrerPage(URL string) string {
	splitStringByShelfParam := strings.Split(URL, "?")
	return splitStringByShelfParam[0]
}

func extractIntPages(pagesString string) int {
	numbers := strings.Split(pagesString, " ")
	if len(numbers) != 2 {
		logger.Sugar().Fatalf("failed to parse pagesString: %s", pagesString)
	}
	return strToInt(numbers[0])
}

func extractGenres(doc *goquery.Document) []string {
	genres := []string{}

	doc.Find(".actionLinkLite.bookPageGenreLink").Each(func(i int, genre *goquery.Selection) {
		genres = append(genres, genre.Text())
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

	return otherEditionCovers
}
