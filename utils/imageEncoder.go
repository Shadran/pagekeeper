package utils

import "io"

type ImageEncoder interface {
	Encode(reader io.Reader) string
	Decode(source string) io.Reader
}
