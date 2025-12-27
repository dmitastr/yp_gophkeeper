package secrets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	mockconfig "gophkeep/internal/mocks/config"
	mocklogger "gophkeep/internal/mocks/logger"
	mockservice "gophkeep/internal/mocks/service_provider"
	testdata "gophkeep/internal/mocks/test_data"
)

func TestSecretHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SecretHandlerTestSuite))
}

type SecretHandlerTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	handlers SecretsHandler
	testData testdata.TestData
}

func (s *SecretHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	mockLogger := mocklogger.NewMockILogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	mockService := mockservice.NewMockIService(ctrl)
	mockService.EXPECT().AddSecret(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockService.EXPECT().GetAllSecrets(gomock.Any()).Return(nil, nil).AnyTimes()
	mockService.EXPECT().GetSecret(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockService.EXPECT().UpdateSecret(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockService.EXPECT().DeleteSecret(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockCfg := mockconfig.NewMockConfigProvider(ctrl)
	mockCfg.EXPECT().Logger().AnyTimes().Return(mockLogger).AnyTimes()

	s.testData = testdata.GetTestData()
	s.ctrl = ctrl
	s.handlers = NewSecretsHandler(mockCfg, mockService)

}

func (s *SecretHandlerTestSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *SecretHandlerTestSuite) TestAddSecret() {
	type args struct {
		request     any
		jsonEncoded bool
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "valid input",
			args:       args{request: &requests.SecretRequest{BodyString: "secret", SecretType: models.TEXT}, jsonEncoded: true},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong json",
			args:       args{request: []byte("secret"), jsonEncoded: false},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ctx := testdata.GetTestGinContext(w)
			testdata.MockJsonPost(ctx, nil, tt.args.request, tt.args.jsonEncoded)

			s.handlers.AddSecret(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func (s *SecretHandlerTestSuite) TestUpdateSecret() {
	type args struct {
		request     any
		jsonEncoded bool
		params      []gin.Param
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid input",
			args: args{
				request:     &requests.SecretRequest{BodyString: "secret", SecretType: models.TEXT},
				jsonEncoded: true,
				params:      []gin.Param{{Key: "secretID", Value: "1"}},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "missing secretID",
			args: args{
				request:     &requests.SecretRequest{BodyString: "secret", SecretType: models.TEXT},
				jsonEncoded: true,
				params:      []gin.Param{{Key: "secret", Value: "1"}},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			args: args{
				request:     []byte("secret"),
				jsonEncoded: false,
				params:      []gin.Param{{Key: "secretID", Value: "1"}},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ctx := testdata.GetTestGinContext(w)
			testdata.MockJsonPost(ctx, tt.args.params, tt.args.request, tt.args.jsonEncoded)

			s.handlers.UpdateSecret(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func (s *SecretHandlerTestSuite) TestGetSecret() {
	type args struct {
		params []gin.Param
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid input",
			args: args{
				params: []gin.Param{{Key: "secretID", Value: "1"}},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "missing secretID",
			args: args{
				params: []gin.Param{{Key: "secret", Value: "1"}},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ctx := testdata.GetTestGinContext(w)
			testdata.MockJsonGet(ctx, tt.args.params, nil)

			s.handlers.GetSecret(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func (s *SecretHandlerTestSuite) TestDeleteSecret() {
	type args struct {
		params []gin.Param
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid input",
			args: args{
				params: []gin.Param{{Key: "secretID", Value: "1"}},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "missing secretID",
			args: args{
				params: []gin.Param{{Key: "secret", Value: "1"}},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ctx := testdata.GetTestGinContext(w)
			testdata.MockJsonDelete(ctx, tt.args.params, nil)

			s.handlers.DeleteSecret(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
