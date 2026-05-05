package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

type ICompressor interface {
	Compress(data []byte) ([]byte, error)
}
type Compressor struct {
}

func NewCompressor() ICompressor {
	return &Compressor{}
}

func (c *Compressor) Compress(data []byte) ([]byte, error) {
	var compressed bytes.Buffer

	gw := gzip.NewWriter(&compressed)
	if _, err := gw.Write(data); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}
	return compressed.Bytes(), nil
}
