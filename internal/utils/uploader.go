package utils

import (
	"io"
)

type Uploader interface {
	Upload(file io.Reader, filename string) (string, error)
}
