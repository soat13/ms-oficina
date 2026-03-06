package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/soat13/fase-1-oficina/internal/auth/infra/out/db"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestGetByCPFWithError(t *testing.T) {
	t.Run("should return error when db query fails", func(t *testing.T) {
		ctx := context.Background()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		bunDB := bun.NewDB(sqlDB, pgdialect.New())
		defer bunDB.Close()

		reader := db.NewBunUserReader(bunDB)

		wantErr := errors.New("db error")
		email := "00063958466"

		mock.ExpectQuery(`users.*document.*00063958466`).
			WillReturnError(wantErr)

		user, queryErr := reader.GetByCPF(ctx, email)

		require.Nil(t, user)
		require.ErrorIs(t, queryErr, wantErr)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
