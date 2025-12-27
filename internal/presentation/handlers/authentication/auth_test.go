package authentication

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gophkeep/internal/core/responses"
	mockconfig "gophkeep/internal/mocks/config"
	mocklogger "gophkeep/internal/mocks/logger"
	mockservice "gophkeep/internal/mocks/service_provider"
	testdata "gophkeep/internal/mocks/test_data"
	"gophkeep/internal/presentation/params"
)

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

type AuthHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	handlers    AuthHandler
	mockService *mockservice.MockIService
}

func (s *AuthHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	mockLogger := mocklogger.NewMockILogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	mockService := mockservice.NewMockIService(ctrl)

	s.mockService = mockService

	mockCfg := mockconfig.NewMockConfigProvider(ctrl)
	mockCfg.EXPECT().Logger().AnyTimes().Return(mockLogger).AnyTimes()

	s.ctrl = ctrl
	s.handlers = NewAuthHandler(mockCfg, mockService)

}

func (s *AuthHandlerTestSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *AuthHandlerTestSuite) TestRegisterUser() {
	type args struct {
		request     any
		jsonEncoded bool
	}
	tests := []struct {
		name       string
		args       args
		token      string
		serviceErr error
		wantStatus int
	}{
		{
			name:       "valid input",
			args:       args{request: &params.AuthRequestObject{Username: "username", Password: "password"}, jsonEncoded: true},
			token:      "token",
			serviceErr: nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong json",
			args:       args{request: []byte("secret"), jsonEncoded: false},
			token:      "token",
			serviceErr: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			args:       args{request: &params.AuthRequestObject{Username: "username", Password: "password"}, jsonEncoded: true},
			token:      "token",
			serviceErr: errors.New("service error"),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ctx := testdata.GetTestGinContext(w)
			testdata.MockJsonPost(ctx, nil, tt.args.request, tt.args.jsonEncoded)

			s.mockService.EXPECT().RegisterUser(ctx, gomock.Any()).Return(tt.token, tt.serviceErr).AnyTimes()
			s.handlers.RegisterUser(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var resp responses.AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, tt.token, resp.Token)
			}
		})
	}
}

func (s *AuthHandlerTestSuite) TestLoginUser() {
	type args struct {
		request     any
		jsonEncoded bool
	}
	tests := []struct {
		name       string
		args       args
		token      string
		serviceErr error
		wantStatus int
	}{
		{
			name:       "valid input",
			args:       args{request: &params.AuthRequestObject{Username: "username", Password: "password"}, jsonEncoded: true},
			token:      "token",
			serviceErr: nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong json",
			args:       args{request: []byte("secret"), jsonEncoded: false},
			token:      "token",
			serviceErr: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			args:       args{request: &params.AuthRequestObject{Username: "username", Password: "password"}, jsonEncoded: true},
			token:      "token",
			serviceErr: errors.New("service error"),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ctx := testdata.GetTestGinContext(w)
			testdata.MockJsonPost(ctx, nil, tt.args.request, tt.args.jsonEncoded)

			s.mockService.EXPECT().LoginUser(ctx, gomock.Any()).Return(tt.token, tt.serviceErr).AnyTimes()
			s.handlers.LoginUser(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var resp responses.AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, tt.token, resp.Token)
			}
		})
	}
}
