package linkding

import (
	"net/url"
	"time"
)

type Bookmark struct {
	Id    int
	Url   string
	Title string
}

type Asset struct {
	Id          int
	AssetType   string `json:"asset_type"`
	ContentType string `json:"content_type"`
	DisplayName string `json:"display_name"`
}

type PagedResponse[T any] struct {
	Count    int
	Next     *string
	Previous *string
	Results  []T
}

type Client struct {
	BaseUrl url.URL
	Token   string
}

type BookmarksQuery struct {
	Tag           string
	ModifiedSince time.Time
}
