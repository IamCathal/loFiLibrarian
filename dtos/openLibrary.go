package dtos

type OpenLibraryBook struct {
	Publishers         []string                `json:"publishers"`
	Identifiers        OtherServiceIdentifiers `json:"identifiers"`
	Weight             string                  `json:"weight"`
	Covers             []int                   `json:"covers"`
	PhysicalFormat     string                  `json:"physical_format"`
	Key                string                  `json:"key"`
	Authors            []Authors               `json:"authors"`
	Ocaid              string                  `json:"ocaid"`
	Subjects           []string                `json:"subjects"`
	SourceRecords      []string                `json:"source_records"`
	Title              string                  `json:"title"`
	NumberOfPages      int                     `json:"number_of_pages"`
	Isbn13             []string                `json:"isbn_13"`
	Isbn10             []string                `json:"isbn_10"`
	PublishDate        string                  `json:"publish_date"`
	Works              []Works                 `json:"works"`
	Type               Type                    `json:"type"`
	PhysicalDimensions string                  `json:"physical_dimensions"`
	Lccn               []string                `json:"lccn"`
	LcClassifications  []string                `json:"lc_classifications"`
	LocalID            []string                `json:"local_id"`
	OclcNumbers        []string                `json:"oclc_numbers"`
	LatestRevision     int                     `json:"latest_revision"`
	Revision           int                     `json:"revision"`
	Created            Created                 `json:"created"`
	LastModified       LastModified            `json:"last_modified"`
}

type OtherServiceIdentifiers struct {
	Librarything []string `json:"librarything"`
	Goodreads    []string `json:"goodreads"`
}
type Authors struct {
	Key string `json:"key"`
}
type Works struct {
	Key string `json:"key"`
}
type Type struct {
	Key string `json:"key"`
}
type Created struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type LastModified struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
