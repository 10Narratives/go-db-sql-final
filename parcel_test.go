package main

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		parcel *Parcel
	}

	var (
		number    int64  = 101
		client    int64  = 102
		address   string = "test address"
		status    string = ParcelStatusRegistered
		createdAt string = "test time"
	)

	tests := []struct {
		name       string
		mocks      func(dbMock sqlmock.Sqlmock)
		args       args
		wantParcel require.ValueAssertionFunc
		wantErr    require.ErrorAssertionFunc
	}{
		{
			name: "success",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec("INSERT INTO parcel").
					WithArgs(client, status, address, createdAt).
					WillReturnResult(sqlmock.NewResult(number, 1))
			},
			args: args{
				parcel: &Parcel{
					Client:    client,
					Address:   address,
					Status:    status,
					CreatedAt: createdAt,
				},
			},
			wantParcel: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcel, ok := got.(*Parcel)
				require.True(t, ok)
				require.NotNil(t, parcel, i...)
				require.Equal(t, number, parcel.Number, i...)
				require.Equal(t, client, parcel.Client, i...)
				require.Equal(t, address, parcel.Address, i...)
				require.Equal(t, status, parcel.Status, i...)
				require.Equal(t, createdAt, parcel.CreatedAt, i...)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec("INSERT INTO parcel").
					WithArgs(client, status, address, createdAt).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				parcel: &Parcel{
					Client:    client,
					Address:   address,
					Status:    status,
					CreatedAt: createdAt,
				},
			},
			wantParcel: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcel, ok := got.(*Parcel)
				require.True(t, ok)
				require.NotNil(t, parcel, i...)
				require.Equal(t, int64(0), parcel.Number, i...)
				require.Equal(t, client, parcel.Client, i...)
				require.Equal(t, address, parcel.Address, i...)
				require.Equal(t, status, parcel.Status, i...)
				require.Equal(t, createdAt, parcel.CreatedAt, i...)
			},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.EqualError(t, err, "database error", i...)
			},
		},
		{
			name:  "no parcel",
			mocks: func(dbMock sqlmock.Sqlmock) {},
			args: args{
				parcel: nil,
			},
			wantParcel: require.Nil,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.EqualError(t, err, "gotten pointer is equal to nil", i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			store := NewParcelStore(db)
			tt.mocks(dbMock)

			err = store.Add(tt.args.parcel)
			tt.wantErr(t, err)
			tt.wantParcel(t, tt.args.parcel)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}
