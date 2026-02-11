package linkding

import (
	"net/url"
	"time"
)

type Bookmark struct {
	Id       int
	Url      string
	Title    string
	TagNames []string `json:"tag_names"`
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
	Tags          []string
	BundleId      int
	ModifiedSince time.Time
}

type UserProfile struct {
	Version string
}
