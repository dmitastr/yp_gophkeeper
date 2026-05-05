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
	basePath    string
}

func (s *SecretsPGDatasourceTestSuite) SetupSuite() {
	basePath := "C:/Users/dastr/Programming/go_practicum_course/yp_gophkeeper"
	s.basePath = basePath

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

	s.repository = db
}

func (s *SecretsPGDatasourceTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		s.T().Log("Error while terminating PG container", err)
	}
}

func (s *SecretsPGDatasourceTestSuite) TearDownTest() {
	query := "DELETE FROM users; DELETE FROM secrets;"
	_, err := s.repository.pool.Exec(s.ctx, query)
	s.Require().NoError(err)
}

func (s *SecretsPGDatasourceTestSuite) SetupTest() {
	sqlBytes, _ := os.ReadFile(path.Join(s.basePath, "database/fixtures/100001_add_tests_data.up.sql"))
	_, err := s.repository.pool.Exec(s.ctx, string(sqlBytes))
	s.Require().NoError(err)
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

func (s *SecretsPGDatasourceTestSuite) TestAddSecret() {
	type args struct {
		userID     models.UserID
		secretType models.SecretType
		secret     []byte
		comment    string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid input",
			args: args{
				userID:     models.UserID(0),
				secretType: models.TEXT,
				secret:     []byte(`secret`),
				comment:    "comment",
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			_, err := s.repository.AddSecret(s.T().Context(), tt.args.userID, tt.args.secretType, tt.args.secret, tt.args.comment, time.Now())
			assert.NoError(t, err)

			secret, err := s.repository.GetSecret(s.ctx, 1, tt.args.userID)

			assert.Equal(t, tt.args.secret, secret.Content)
			assert.Equal(t, tt.args.secretType, secret.Type)
			assert.Equal(t, tt.args.comment, secret.Comment)
		})
	}
}

func (s *SecretsPGDatasourceTestSuite) TestGetSecret() {
	type args struct {
		userID   models.UserID
		secretID int
	}

	secret := &models.Secret{
		Type:    models.TEXT,
		Content: []byte("secret"),
		Comment: "comment",
	}

	tests := []struct {
		name    string
		args    args
		want    *models.Secret
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				userID: models.UserID(0),
			},
			want:    secret,
			wantErr: false,
		},
		{
			name: "secret not found",
			args: args{
				userID: models.UserID(0),
			},
			want:    secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			secretID, err := s.repository.AddSecret(s.T().Context(), tt.args.userID, secret.Type, secret.Content, secret.Comment, time.Now())
			assert.NoError(t, err)
			assert.NotNil(t, secretID)

			if tt.wantErr {
				*secretID += 1000
			}

			got, err := s.repository.GetSecret(s.ctx, *secretID, tt.args.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, got.Content, secret.Content)
			assert.Equal(t, got.Type, secret.Type)
			assert.Equal(t, got.Comment, secret.Comment)
		})
	}
}

func (s *SecretsPGDatasourceTestSuite) TestGetAllSecrets() {
	userID := models.UserID(10)
	secrets := []*models.Secret{
		{
			Type:      models.TEXT,
			Content:   []byte("secret"),
			Comment:   "comment",
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		},
	}

	secretsInfo := make([]models.SecretInfo, len(secrets))

	for i, secret := range secrets {
		secretID, err := s.repository.AddSecret(s.ctx, userID, secret.Type, secret.Content, secret.Comment, secret.CreatedAt)
		s.Require().NoError(err)

		secretsInfo[i] = models.SecretInfo{Type: secret.Type, ID: *secretID, CreatedAt: secret.CreatedAt}
	}

	type args struct {
		userID models.UserID
	}

	tests := []struct {
		name    string
		args    args
		want    []models.SecretInfo
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				userID: userID,
			},
			want:    secretsInfo,
			wantErr: false,
		},
		{
			name: "secrets not found",
			args: args{
				userID: userID + 1000,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			secretsGot, err := s.repository.GetAllSecrets(s.T().Context(), tt.args.userID)

			if tt.wantErr {
				assert.Empty(t, secretsGot)
				return
			}

			for i, secret := range secretsGot {
				secret.CreatedAt = secret.CreatedAt.UTC()
				secretsGot[i] = secret
			}

			assert.NoError(t, err)
			assert.EqualValues(t, secretsInfo, secretsGot)
		})
	}
}

func (s *SecretsPGDatasourceTestSuite) TestUpdateSecret() {
	type args struct {
		userID     models.UserID
		secret     *models.Secret
		secretNew  []byte
		commentNew string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid input",
			args: args{
				userID:     models.UserID(0),
				secret:     &models.Secret{Type: models.TEXT, Content: []byte(`secret`), Comment: "comment", UpdatedAt: time.Now().UTC().Truncate(time.Microsecond)},
				secretNew:  []byte("secretNew"),
				commentNew: "commentNew",
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			secretID, err := s.repository.AddSecret(
				s.T().Context(),
				tt.args.userID,
				tt.args.secret.Type,
				tt.args.secret.Content,
				tt.args.secret.Comment,
				tt.args.secret.UpdatedAt)
			assert.NoError(t, err)
			tt.args.secret.ID = *secretID

			tt.args.secret.Comment = tt.args.commentNew
			tt.args.secret.Content = tt.args.secretNew
			err = s.repository.UpdateSecret(s.ctx, tt.args.secret, tt.args.userID)
			assert.NoError(t, err)

			secretGot, err := s.repository.GetSecret(s.ctx, *secretID, tt.args.userID)
			assert.NoError(t, err)
			secretGot.UpdatedAt = secretGot.UpdatedAt.UTC()

			assert.Equal(t, tt.args.commentNew, secretGot.Comment)
			assert.Equal(t, tt.args.secretNew, secretGot.Content)
			assert.Equal(t, tt.args.secret.UpdatedAt, secretGot.UpdatedAt)
		})
	}
}

func (s *SecretsPGDatasourceTestSuite) TestDeleteSecret() {
	type args struct {
		userID models.UserID
		secret *models.Secret
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid input",
			args: args{
				userID: models.UserID(0),
				secret: &models.Secret{Type: models.TEXT, Content: []byte(`secret`), Comment: "comment", UpdatedAt: time.Now().UTC().Truncate(time.Microsecond)},
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			secretID, err := s.repository.AddSecret(
				s.T().Context(),
				tt.args.userID,
				tt.args.secret.Type,
				tt.args.secret.Content,
				tt.args.secret.Comment,
				tt.args.secret.UpdatedAt)
			assert.NoError(t, err)

			err = s.repository.DeleteSecret(s.ctx, *secretID, tt.args.userID)
			assert.NoError(t, err)

			_, err = s.repository.GetSecret(s.ctx, *secretID, tt.args.userID)
			assert.Error(t, err)
		})
	}
}
