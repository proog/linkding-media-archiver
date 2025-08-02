package linkding

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func NewClient(baseUrl string, token string) (*Client, error) {
	if baseUrl == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	parsedUrl, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	if !parsedUrl.IsAbs() {
		return nil, fmt.Errorf("base URL is not absolute: %s", baseUrl)
	}

	return &Client{BaseUrl: *parsedUrl, Token: token}, nil
}

func (client *Client) GetBookmarks(query BookmarksQuery) ([]Bookmark, error) {
	logger := slog.With("tag", query.Tag, "modifiedSince", query.ModifiedSince)
	logger.Debug("Fetching bookmarks")

	endpointUrl := client.url("bookmarks/")
	queryParams := endpointUrl.Query()

	if query.Tag != "" {
		queryParams.Set("q", "#"+query.Tag)
	}

	if !query.ModifiedSince.IsZero() {
		queryParams.Set("modified_since", query.ModifiedSince.UTC().Format(time.RFC3339))
	}

	endpointUrl.RawQuery = queryParams.Encode()

	results, err := getAllItems[Bookmark](endpointUrl, func(u url.URL) (*http.Response, error) {
		return client.get(u)
	})

	if err == nil {
		logger.Debug("Fetched bookmarks", "count", len(results))
	}

	return results, err
}

func (client *Client) GetBookmarkAssets(bookmarkId int) ([]Asset, error) {
	logger := slog.With("bookmarkId", bookmarkId)
	logger.Debug("Fetching assets for bookmark")

	endpointUrl := client.url("bookmarks", strconv.Itoa(bookmarkId), "assets/")

	results, err := getAllItems[Asset](endpointUrl, func(u url.URL) (*http.Response, error) {
		return client.get(u)
	})

	if err == nil {
		slog.Debug("Fetched assets for bookmark", "count", len(results))
	}

	return results, err
}

func (client *Client) DownloadBookmarkAsset(bookmarkId int, assetId int) (io.ReadCloser, error) {
	logger := slog.With("bookmarkId", bookmarkId, "assetId", assetId)
	logger.Debug("Downloading asset content")

	endpointUrl := client.url("bookmarks", strconv.Itoa(bookmarkId), "assets", strconv.Itoa(assetId), "download/")
	resp, err := client.get(endpointUrl)

	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

func (client *Client) AddBookmarkAsset(bookmarkId int, file *os.File) (*Asset, error) {
	logger := slog.With("bookmarkId", bookmarkId)
	logger.Debug("Adding asset for bookmark")

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileName := stat.Name()
	fileSize := stat.Size()
	mimeType, err := GetMimeType(fileName)
	if err != nil {
		return nil, err
	}

	const fieldName = "file"

	// Pipe the file contents directly to the http request without loading the entire file into memory
	readBody, writeBody := io.Pipe()
	formData := multipart.NewWriter(writeBody)

	partErr := make(chan error, 1)
	go func() {
		defer writeBody.Close()
		defer close(partErr)

		part, err := createMultipartPart(formData, fieldName, fileName, mimeType)
		if err != nil {
			partErr <- err
			return
		}

		// This blocks until the http request reads from the pipe
		if _, err := io.CopyN(part, file, fileSize); err != nil {
			writeBody.CloseWithError(err)
			partErr <- err
			return
		}

		// Important! Write the closing boundary to the part
		partErr <- formData.Close()
	}()

	url := client.url("bookmarks", strconv.Itoa(bookmarkId), "assets/upload/")
	headers := map[string]string{"Content-Type": formData.FormDataContentType()}
	contentLength := emptyMultipartPartLength(fieldName, fileName, mimeType) + fileSize

	resp, err := client.send(http.MethodPost, url, headers, readBody, contentLength)

	if err := errors.Join(err, <-partErr); err != nil {
		return nil, err
	}

	logger.Debug("Added asset")

	return deserialize[Asset](resp)
}

func (client *Client) url(path ...string) url.URL {
	path = append([]string{"api"}, path...)
	return *client.BaseUrl.JoinPath(path...)
}

func (client *Client) get(url url.URL) (*http.Response, error) {
	return client.send(http.MethodGet, url, nil, nil, 0)
}

func (client *Client) send(method string, url url.URL, headers map[string]string, body io.Reader, contentLength int64) (*http.Response, error) {
	req, err := http.NewRequest(method, url.String(), body)

	if err != nil {
		return nil, err
	}

	if contentLength != 0 {
		req.ContentLength = contentLength
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	req.Header.Set("Authorization", fmt.Sprint("Token ", client.Token))

	logger := slog.With("method", method, "url", url.String())
	logger.Debug("Sending HTTP request")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return resp, err
	}

	logger = logger.With("statusCode", resp.StatusCode)
	logger.Debug("Received HTTP response")

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("expected success status code, was %s", resp.Status)
	}

	return resp, nil
}
