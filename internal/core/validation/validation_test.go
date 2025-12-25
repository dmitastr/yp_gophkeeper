package validation

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
	mockconfig "gophkeep/internal/mocks/config"
	mock_logger "gophkeep/internal/mocks/logger"
)

var passwordTestData = &models.Password{
	Login:    "login",
	Password: "password",
	Comment:  "comment",
}

var bankcardTestData = &models.BankCard{
	Number:         5120350100064537,
	ExpiryData:     "03/30",
	CardholderName: "Ivan",
	CVC:            "000",
	Comment:        "comment",
}

var textTestData = []byte("text secret")

var binaryTestData = []byte("binary secret")

type ValidatorTestSuite struct {
	suite.Suite
	cfg       config.ConfigProvider
	validator *Validator
	ctx       context.Context
	ctrl      *gomock.Controller
	testData  map[models.SecretType][]byte
}

func (s *ValidatorTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.ctrl = ctrl

	mockLogger := mock_logger.NewMockILogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockCfg := mockconfig.NewMockConfigProvider(ctrl)
	mockCfg.EXPECT().Logger().AnyTimes().Return(mockLogger).AnyTimes()

	s.testData = make(map[models.SecretType][]byte)
	passwordData, _ := json.Marshal(passwordTestData)
	s.testData[models.PASSWORD] = passwordData
	bankcardData, _ := json.Marshal(bankcardTestData)
	s.testData[models.BANK_CARD] = bankcardData

	s.testData[models.TEXT] = textTestData
	s.testData[models.BINARY] = binaryTestData

	s.validator = NewValidator(mockCfg)

}

func (s *ValidatorTestSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *ValidatorTestSuite) TestValidator_decode() {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{message: base64.StdEncoding.EncodeToString([]byte("hello world"))},
			want:    []byte("hello world"),
			wantErr: false,
		},
		{
			name:    "message is not encoded input",
			args:    args{message: "hello world"},
			want:    []byte("hello world"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.validator.decode(tt.args.message)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *ValidatorTestSuite) TestValidator_parseBankCardSecret() {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{content: s.testData[models.BANK_CARD]},
			want:    s.testData[models.BANK_CARD],
			wantErr: false,
		},
		{
			name:    "wrong secret type",
			args:    args{content: s.testData[models.TEXT]},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.validator.parseBankCardSecret(tt.args.content)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *ValidatorTestSuite) TestValidator_parseBinarySecret() {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{content: s.testData[models.BINARY]},
			want:    s.testData[models.BINARY],
			wantErr: false,
		},
		{
			name:    "empty input",
			args:    args{content: make([]byte, 0)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil input",
			args:    args{content: nil},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.validator.parseBinarySecret(tt.args.content)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *ValidatorTestSuite) TestValidator_parsePasswordSecret() {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{content: s.testData[models.PASSWORD]},
			want:    s.testData[models.PASSWORD],
			wantErr: false,
		},
		{
			name:    "wrong secret type",
			args:    args{content: s.testData[models.TEXT]},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.validator.parsePasswordSecret(tt.args.content)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *ValidatorTestSuite) TestValidator_parseTextSecret() {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{content: s.testData[models.TEXT]},
			want:    textTestData,
			wantErr: false,
		},
		{
			name:    "empty text",
			args:    args{content: []byte{0xff, 0xfe, 0xfd}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.validator.parseTextSecret(tt.args.content)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}
