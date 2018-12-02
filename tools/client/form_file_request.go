package client

/*
do request for  formData with file
*/

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func NewFileRequest(uri string, params map[string]string, fileHeader *multipart.FileHeader) (*http.Request, error) {
	fileContents, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer fileContents.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(fileHeader.Header)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fileContents)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
