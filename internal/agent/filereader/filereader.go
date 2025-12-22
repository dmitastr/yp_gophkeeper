package filereader

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
)

type inputParserFunc func(input, fileName, comment string) (*requests.SecretRequest, error)

type FileInputReader interface {
	FileRead(file string) ([]byte, error)
	FileWrite(content []byte, file string) error
	ParseInput(input, fileName, comment string, secretType models.SecretType) (*requests.SecretRequest, error)
}

type FileInputReaderImpl struct {
	inputFuncs map[models.SecretType]inputParserFunc
}

func NewFileInputReader() FileInputReader {
	fr := &FileInputReaderImpl{}
	fr.inputFuncs = map[models.SecretType]inputParserFunc{
		models.PASSWORD:  fr.parsePasswordJSON,
		models.BANK_CARD: fr.parseBankCardJSON,
		models.TEXT:      fr.parseText,
		models.BINARY:    fr.parseBinary,
	}
	return fr
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

func (fr *FileInputReaderImpl) ParseInput(input, fileName, comment string, secretType models.SecretType) (*requests.SecretRequest, error) {
	parser, ok := fr.inputFuncs[secretType]
	if !ok {
		return nil, fmt.Errorf("invalid secret type: %s", secretType)
	}
	return parser(input, fileName, comment)
}

func (fr *FileInputReaderImpl) parsePasswordJSON(_, fileName, comment string) (*requests.SecretRequest, error) {
	if fileName == "" {
		return nil, errors.New("no file name provided")
	}

	fileContent, err := fr.FileRead(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var s models.Secret
	if err := json.Unmarshal(fileContent, &s); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	secretRequest := &requests.SecretRequest{
		BodyString: base64.StdEncoding.EncodeToString(fileContent),
		SecretType: models.PASSWORD,
		Comment:    comment,
	}

	if s.Comment != "" {
		secretRequest.Comment = s.Comment
	}

	return secretRequest, nil
}

func (fr *FileInputReaderImpl) parseBankCardJSON(_, fileName, comment string) (*requests.SecretRequest, error) {
	if fileName == "" {
		return nil, errors.New("no file name provided")
	}

	fileContent, err := fr.FileRead(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var s models.BankCard
	if err := json.Unmarshal(fileContent, &s); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	secretRequest := &requests.SecretRequest{
		BodyString: base64.StdEncoding.EncodeToString(fileContent),
		SecretType: models.BANK_CARD,
		Comment:    comment,
	}

	if s.Comment != "" {
		secretRequest.Comment = s.Comment
	}

	return secretRequest, nil
}

func (fr *FileInputReaderImpl) parseText(input, fileName, comment string) (secretRequest *requests.SecretRequest, err error) {
	var body string
	if fileName != "" {
		fileContent, err := fr.FileRead(fileName)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		body = base64.StdEncoding.EncodeToString(fileContent)

	} else if input != "" {
		body = input

	} else {
		err = errors.New("file and input arguments are both empty")
	}
	secretRequest = &requests.SecretRequest{BodyString: body, SecretType: models.TEXT, Comment: comment}

	return secretRequest, err
}

func (fr *FileInputReaderImpl) parseBinary(input, fileName, comment string) (secretRequest *requests.SecretRequest, err error) {
	if fileName != "" {
		fileContent, err := fr.FileRead(fileName)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(fileContent)

		secretRequest = &requests.SecretRequest{BodyString: encoded, SecretType: models.BINARY, Comment: comment}

	} else if input != "" {
		err = errors.New("binary input from argument is not supported")

	} else {
		err = errors.New("file and input arguments are both empty")
	}

	return secretRequest, err
}
