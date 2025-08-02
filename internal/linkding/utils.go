package linkding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
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

func getAllItems[T any](initialUrl url.URL, getPage func(url.URL) (*http.Response, error)) ([]T, error) {
	results := make([]T, 0)
	nextUrl := &initialUrl

	for nextUrl != nil {
		resp, err := getPage(*nextUrl)

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

func createMultipartPart(formData *multipart.Writer, fieldName, fileName, mimeType string) (io.Writer, error) {
	// Can't use formData.CreateFormFile() as it forces the application/octet-stream content type
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	header.Set("Content-Type", mimeType)

	return formData.CreatePart(header)
}

func emptyMultipartPartLength(fieldName, fileName, mimeType string) int64 {
	body := &bytes.Buffer{}
	formData := multipart.NewWriter(body)

	createMultipartPart(formData, fieldName, fileName, mimeType)
	formData.Close()

	return int64(body.Len())
}
