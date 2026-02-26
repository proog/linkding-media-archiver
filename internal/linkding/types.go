package linkding

import (
	"net/url"
	"time"
)

type Bookmark struct {
	Id       int      `json:"id"`
	Url      string   `json:"url"`
	Title    string   `json:"title"`
	TagNames []string `json:"tag_names"`
}

type Asset struct {
	Id          int    `json:"id"`
	AssetType   string `json:"asset_type"`
	ContentType string `json:"content_type"`
	DisplayName string `json:"display_name"`
}

type PagedResponse[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type Client struct {
	BaseUrl url.URL
	Token   string
}

type BookmarksQuery struct {
	Tags          []string
	BundleId      int
	ModifiedSince time.Time
}

type UserProfile struct {
	Version string `json:"version"`
}
