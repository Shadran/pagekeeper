package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

type GzipBase64ImageEncoder struct {
}

func (e *GzipBase64ImageEncoder) Encode(reader io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	_, err := io.Copy(gzip.NewWriter(base64.NewEncoder(base64.StdEncoding, buf)), reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func (e *GzipBase64ImageEncoder) Decode(source string) (io.Reader, error) {
	return gzip.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(source)))
}
