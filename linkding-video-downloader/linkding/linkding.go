package linkding

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func NewClient(baseUrl string, token string) (*Client, error) {
	parsedUrl, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	return &Client{BaseUrl: parsedUrl, Token: token}, nil
}

func (client *Client) GetBookmarks(tag string) ([]Bookmark, error) {
	logger := slog.With("tag", tag)
	logger.Info("Fetching bookmarks")

	endpointUrl := client.url("bookmarks/")
	query := endpointUrl.Query()
	query.Set("q", "#"+tag)
	endpointUrl.RawQuery = query.Encode()

	results, err := getAllItems[Bookmark](endpointUrl, func(url *url.URL) (*http.Response, error) {
		return client.get(url, nil)
	})

	if err == nil {
		logger.Info("Fetched bookmarks", "count", len(results))
	}

	return results, err
}

func (client *Client) GetBookmarkAssets(bookmarkId int) ([]Asset, error) {
	logger := slog.With("bookmarkId", bookmarkId)
	logger.Info("Fetching assets for bookmark")

	endpointUrl := client.url("bookmarks", strconv.Itoa(bookmarkId), "assets/")

	results, err := getAllItems[Asset](endpointUrl, func(url *url.URL) (*http.Response, error) {
		return client.get(url, nil)
	})

	if err == nil {
		slog.Info("Fetched assets for bookmark", "count", len(results))
	}

	return results, err
}

func (client *Client) DownloadBookmarkAsset(bookmarkId int, assetId int) (io.ReadCloser, error) {
	logger := slog.With("bookmarkId", bookmarkId, "assetId", assetId)
	logger.Info("Downloading asset content")

	endpointUrl := client.url("bookmarks", strconv.Itoa(bookmarkId), "assets", strconv.Itoa(assetId), "download/")
	resp, err := client.get(endpointUrl, nil)

	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

func (client *Client) AddBookmarkAsset(bookmarkId int, file *os.File) (*Asset, error) {
	logger := slog.With("bookmarkId", bookmarkId)
	logger.Info("Adding asset for bookmark")

	var multipartBody bytes.Buffer
	formData, err := createMultipartBody(&multipartBody, file)

	if err != nil {
		return nil, err
	}

	url := client.url("bookmarks", strconv.Itoa(bookmarkId), "assets/upload/")
	headers := map[string]string{"Content-Type": formData.FormDataContentType()}

	resp, err := client.send(http.MethodPost, url, &multipartBody, headers)

	if err != nil {
		return nil, err
	}

	logger.Info("Added asset")

	return deserialize[Asset](resp)
}

func (client *Client) url(path ...string) *url.URL {
	path = append([]string{"api"}, path...)
	return client.BaseUrl.JoinPath(path...)
}

func (client *Client) get(url *url.URL, headers map[string]string) (*http.Response, error) {
	return client.send(http.MethodGet, url, nil, headers)
}

func (client *Client) send(method string, url *url.URL, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url.String(), body)

	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	req.Header.Set("Authorization", fmt.Sprint("Token ", client.Token))

	logger := slog.With("method", method, "url", url)
	logger.Info("Sending HTTP request")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return resp, err
	}

	logger = logger.With("statusCode", resp.StatusCode)
	logger.Info("Received HTTP response")

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("expected success status code, was %s", resp.Status)
	}

	return resp, nil
}
