package repository_test

import (
	"context"
	"database/sql"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/phuslu/log"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	driver  *sql.DB
	mock    sqlmock.Sqlmock
	errMain error
)

func TestMain(m *testing.M) {
	driver, mock, errMain = sqlmock.New()
	defer func() {
		if err := driver.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close mock driver")
		}
	}()
	DB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn:       driver,
		DriverName: "postgres",
	}), &gorm.Config{SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Info)})
	m.Run()
}

// Initializes UsersRepository with valid DB and appLog
func TestNewUsersRepositoryImpl_ValidInputs(t *testing.T) {
	db := &gorm.DB{}
	appLog := &log.DefaultLogger

	repo, err := repository.NewUsersRepositoryImpl(db, appLog)

	require.NoError(t, err)

	require.NotNil(t, repo)

	require.Equal(t, db, repo.DB)

	require.Equal(t, appLog, repo.Logger)
}

// Returns error when DB is nil and Log is nil
func TestNewUsersRepositoryImpl_DBIsNil_LogIsNil(t *testing.T) {
	repo, err := repository.NewUsersRepositoryImpl(nil, nil)

	require.Error(t, err)

	require.Nil(t, repo)

	expectedErr := "DB or Logger is nil"
	require.Equal(t, expectedErr, err.Error())
}

// Retrieves all users when no pagination parameters are provided
func TestGetAllNoPagination(t *testing.T) {
	id := ksuid.New().String()

	ctx := context.Background()
	logger := &log.DefaultLogger
	repo, err := repository.NewUsersRepositoryImpl(DB, logger)
	if err != nil {
		t.Fatalf("failed to create UsersRepository: %v", err)
	}

	testCases := []struct {
		name         string
		mockQuery    string
		queryParams  *request.Page
		expectedRows int
		expectError  bool
		expectMatch  bool
	}{
		{
			name:         "Successfully Mock Get All With No Pagination",
			mockQuery:    regexp.QuoteMeta(`SELECT * FROM "users" ORDER BY id ASC LIMIT $1`),
			queryParams:  &request.Page{Size: 10},
			expectedRows: 3,
			expectError:  false,
			expectMatch:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "username", "email"}).
				AddRow(id, "user", "user@gmail.com").
				AddRow(ksuid.New(), "user1", "user1@example.com").
				AddRow(ksuid.New(), "user2", "user2@example.com")

			mockQuery := mock.ExpectQuery(testCase.mockQuery)

			args := []interface{}{testCase.queryParams.Size}
			if testCase.queryParams.Before != "" {
				args = append(args, testCase.queryParams.Before)
			}
			if testCase.queryParams.After != "" {
				args = append(args, testCase.queryParams.After)
			}

			for _, arg := range args {
				mockQuery = mockQuery.WithArgs(arg)
			}

			mockQuery.WillReturnRows(rows)
			users, err := repo.GetAll(ctx, testCase.queryParams)

			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedRows, len(users))
			}

			// Validate whether the expectations matched
			match := mock.ExpectationsWereMet() == nil
			require.Equal(t, testCase.expectMatch, match)
		})
	}
}

func TestUsersRepositoryMethods(t *testing.T) {
	logger := &log.DefaultLogger
	repo, err := repository.NewUsersRepositoryImpl(DB, logger)
	if err != nil {
		t.Fatalf("failed to create UsersRepository: %v", err)
	}

	// Example entity for testing
	user := &entity.Users{ID: ksuid.New().String(), Username: "user", Email: "user@example.com", Password: "examplepassword", CreatedAt: time.Now().UnixNano()}
	ctx := context.Background()
	t.Run("Create Users Case", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "users" (.+) VALUES (.+)`).
			WithArgs(user.ID, user.Username, user.Email, user.Password, user.CreatedAt, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure Create Users Case Because Already Exist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "users" (.+) VALUES (.+)`).
			WithArgs(user.ID, user.Username, user.Email, user.Password, user.CreatedAt, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrDuplicatedKey)
		mock.ExpectRollback()

		err := repo.Create(ctx, user)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Users Case", func(t *testing.T) {
		mock.ExpectBegin()

		// Update the regular expression to match the actual SQL query format

		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+`).
			WithArgs(user.Username, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		id := user.ID
		user.ID = ""
		err := repo.Update(ctx, user, id)
		require.NoError(t, err)
		// Verify that all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure Update Users Case Because Record Not Found", func(t *testing.T) {
		mock.ExpectBegin()

		// Update the regular expression to match the actual SQL query format

		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+`).
			WithArgs(user.Username, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		id := user.ID
		user.ID = ""
		err := repo.Update(ctx, user, id)
		require.Error(t, err)
		// Verify that all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Users Case", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+`).WithArgs(sqlmock.AnyArg(), user.ID, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, user.ID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure Delete Users Case Because Record Not Found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+`).WithArgs(sqlmock.AnyArg(), user.ID, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Delete(ctx, user.ID)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountById Users Case", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+`).WithArgs(user.ID, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		count, err := repo.CountById(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure CountById Users Case Because Record Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+`).WithArgs(user.ID, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)

		count, err := repo.CountById(ctx, user.ID)
		require.Error(t, err)
		require.Equal(t, int64(0), count)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Count Users Case", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		count, err := repo.Count(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure Count Users Case Because Invalid Value", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "users"`).WillReturnError(gorm.ErrInvalidValue)

		count, err := repo.Count(ctx)
		require.Error(t, err)
		require.Equal(t, int64(0), count)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetById Users Case", func(t *testing.T) {
		var result entity.Users
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+ LIMIT .+`).WithArgs(user.ID, sqlmock.AnyArg(), 1).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(user.ID, user.Username, user.Email))

		err := repo.GetById(ctx, &result, user.ID)
		require.NoError(t, err)
		require.Equal(t, user.Username, user.Username)
		require.Equal(t, user.Email, user.Email)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure GetById Users Case Because Record Not Found", func(t *testing.T) {
		var result entity.Users
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+ LIMIT .+`).WithArgs(user.ID, sqlmock.AnyArg(), 1).WillReturnError(gorm.ErrRecordNotFound)

		err := repo.GetById(ctx, &result, user.ID)
		require.Error(t, err)
		require.Equal(t, user.Username, user.Username)
		require.Equal(t, user.Email, user.Email)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetByEmail Users Case", func(t *testing.T) {
		var result entity.Users
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+ LIMIT .+`).WithArgs(user.Email, sqlmock.AnyArg(), 1).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(user.ID, user.Username, user.Email))

		err := repo.GetByEmail(ctx, &result, user.Email)
		require.NoError(t, err)
		require.Equal(t, user.Username, user.Username)
		require.Equal(t, user.Email, user.Email)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure GetByEmail Users Case Because Record Not Found", func(t *testing.T) {
		var result entity.Users
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+ LIMIT .+`).WithArgs(user.Email, sqlmock.AnyArg(), 1).WillReturnError(gorm.ErrRecordNotFound)

		err := repo.GetByEmail(ctx, &result, user.Email)
		require.Error(t, err)
		require.Equal(t, user.Username, user.Username)
		require.Equal(t, user.Email, user.Email)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ExistByKeyValue Users Case", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+`).WithArgs("example@gmail.com", sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		exist, err := repo.ExistByKeyValue(ctx, map[string]any{"email": "example@gmail.com"})
		require.NoError(t, err)
		require.Equal(t, true, exist)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure ExistByKeyValue Users Case Because Record Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "users" WHERE .+`).WithArgs("example@gmail.com", sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)

		exist, err := repo.ExistByKeyValue(ctx, map[string]any{"email": "example@gmail.com"})
		require.Error(t, err)
		require.Equal(t, false, exist)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RestoreById Users Case", func(t *testing.T) {
		mock.ExpectBegin()
		id := ksuid.New().String()
		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+ AND .+`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), id, 0).WillReturnResult(sqlmock.NewResult(1, 1)) // Menggunakan AnyArg() untuk fleksibilitas
		mock.ExpectCommit()
		err := repo.Restore(ctx, id)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure RestoreById Users Case Because Record Not Found", func(t *testing.T) {
		mock.ExpectBegin()
		id := ksuid.New().String()
		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+ AND .+`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), id, 0).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		err := repo.Restore(ctx, id)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
