package search

type SearchRequest struct {
	Term string `json:"term"`
}

type SearchResult struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}
