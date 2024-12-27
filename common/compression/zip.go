package compression

import (
	"bytes"
	"compress/zlib"
	"io"
)

func CompressString(str string) (string, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err = w.Write([]byte(str)); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

func DecompressString(str string) (string, error) {
	r, err := zlib.NewReader(bytes.NewBuffer([]byte(str)))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}
