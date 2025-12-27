package postgresdb

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
	mockconfig "gophkeep/internal/mocks/config"
	testhelpers "gophkeep/internal/mocks/helpers"
	mocklogger "gophkeep/internal/mocks/logger"
)

func TestNewPostgresStorage(t *testing.T) {
	suite.Run(t, new(SecretsPGDatasourceTestSuite))
}

type SecretsPGDatasourceTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *PostgresStorage
	ctx         context.Context
	ctrl        *gomock.Controller
}

func (s *SecretsPGDatasourceTestSuite) SetupTest() {
	basePath := "C:/Users/dastr/Programming/go_practicum_course/yp_gophkeeper"
	ctx := context.Background()
	s.ctx = ctx
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	s.Require().NoError(err)
	s.pgContainer = pgContainer

	ctrl := gomock.NewController(s.T())
	s.ctrl = ctrl

	mockLogger := mocklogger.NewMockILogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	mockConfig := mockconfig.NewMockConfigProvider(ctrl)
	mockConfig.EXPECT().Logger().AnyTimes().Return(mockLogger).AnyTimes()

	cfg := config.Config{
		DBConnStr:     pgContainer.ConnectionString,
		MigrationsDir: path.Join(basePath, "database/migrations"),
	}

	mockConfig.EXPECT().GetConfig().AnyTimes().Return(cfg)
	db, err := NewPostgresStorage(ctx, mockConfig)
	s.Require().NoError(err)

	sqlBytes, _ := os.ReadFile(path.Join(basePath, "database/fixtures/100001_add_tests_data.up.sql"))
	_, err = db.pool.Exec(ctx, string(sqlBytes))
	s.Require().NoError(err)

	s.repository = db
}

func (s *SecretsPGDatasourceTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		s.T().Log("Error while terminating PG container", err)
	}
}

func (s *SecretsPGDatasourceTestSuite) TestGetUser() {
	type args struct {
		username string
	}

	ts, _ := time.Parse(time.RFC3339, `2024-12-25T10:30:00Z`)
	user := &models.User{
		Username:  "username",
		Hash:      "hash",
		CreatedAt: ts,
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{username: "username"},
			want:    user,
			wantErr: false,
		},
		{
			name:    "no such user",
			args:    args{username: "not_found"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u, err := s.repository.GetUser(s.T().Context(), tt.args.username)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, u)
		})
	}
}

func (s *SecretsPGDatasourceTestSuite) TestAddUser() {
	type args struct {
		users []*models.User
	}

	ts, _ := time.Parse(time.RFC3339, `2024-12-25T10:30:00Z`)
	users := []*models.User{
		{ID: 1, Username: "username_1", Hash: "hash", CreatedAt: ts},
		{ID: 2, Username: "username_2", Hash: "hash", CreatedAt: ts},
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.User
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    args{users: users},
			want:    users,
			wantErr: false,
		},
		{
			name:    "violation unique constraint",
			args:    args{users: []*models.User{{ID: 1, Username: "username", Hash: "hash", CreatedAt: ts}}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			result := make([]*models.User, 0)
			for _, user := range tt.args.users {
				u, err := s.repository.AddUser(s.T().Context(), user)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				result = append(result, u)
			}
			assert.EqualValues(t, users, result)
		})
	}
}
