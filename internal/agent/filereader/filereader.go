package filereader

import (
	"fmt"
	"os"
)

type FileInputReader interface {
	FileRead(file string) ([]byte, error)
	FileWrite(content []byte, file string) error
}

type FileInputReaderImpl struct {
}

func NewFileInputReader() FileInputReader {
	return &FileInputReaderImpl{}
}

func (fr *FileInputReaderImpl) FileRead(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func (fr *FileInputReaderImpl) FileWrite(content []byte, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
