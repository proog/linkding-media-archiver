package linkding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
)

func deserialize[T any](resp *http.Response) (*T, error) {
	var result T

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return &result, err
}

func getAllItems[T any](url *url.URL, getPage func(*url.URL) (*http.Response, error)) ([]T, error) {
	results := make([]T, 0)
	nextUrl := url

	for nextUrl != nil {
		resp, err := getPage(nextUrl)

		if err != nil {
			return nil, err
		}

		result, err := deserialize[PagedResponse[T]](resp)

		if err != nil {
			return nil, err
		}

		results = append(results, result.Results...)

		if result.Next == nil {
			break
		}

		nextUrl, err = url.Parse(*result.Next)

		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func createMultipartBody(buffer *bytes.Buffer, file *os.File) (*multipart.Writer, error) {
	formData := multipart.NewWriter(buffer)
	defer formData.Close()

	fileName := filepath.Base(file.Name())
	mimeType := getMimeType(fileName)

	slog.Info("Creating multipart with MIME type", "mimeType", mimeType, "fileName", fileName)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", fileName))
	header.Set("Content-Type", mimeType)

	part, err := formData.CreatePart(header)

	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	return formData, nil
}

func getMimeType(fileName string) string {
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))

	if mimeType == "" {
		return "application/octet-stream"
	}

	return mimeType
}
