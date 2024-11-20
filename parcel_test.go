package main

import (
	"database/sql"
	"errors"
	"regexp"
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

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		number    int    = 101
		client    string = "Test Client"
		address   string = "Test Address"
		status    string = "Registered"
		createdAt string = "2023-11-20T10:00:00Z"
	)

	tests := []struct {
		name       string
		mocks      func(dbMock sqlmock.Sqlmock)
		number     int
		wantParcel require.ValueAssertionFunc
		wantErr    require.ErrorAssertionFunc
	}{
		{
			name: "success",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"number", "client", "status", "address", "created_at"}).
					AddRow(number, client, status, address, createdAt)
				dbMock.ExpectQuery("SELECT number, client, status, address, created_at FROM parcel WHERE id = ?").
					WithArgs(number).
					WillReturnRows(rows)
			},
			number: number,
			wantParcel: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcel, ok := got.(Parcel)
				require.True(t, ok)
				require.Equal(t, number, parcel.Number)
				require.Equal(t, client, parcel.Client)
				require.Equal(t, address, parcel.Address)
				require.Equal(t, status, parcel.Status)
				require.Equal(t, createdAt, parcel.CreatedAt)
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectQuery("SELECT number, client, status, address, created_at FROM parcel WHERE id = ?").
					WithArgs(number).
					WillReturnError(sql.ErrNoRows)
			},
			number: number,
			wantParcel: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcel, ok := got.(Parcel)
				require.True(t, ok)
				require.Equal(t, Parcel{}, parcel)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectQuery("SELECT number, client, status, address, created_at FROM parcel WHERE id = ?").
					WithArgs(number).
					WillReturnError(errors.New("database error"))
			},
			number: number,
			wantParcel: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcel, ok := got.(Parcel)
				require.True(t, ok)
				require.Equal(t, Parcel{}, parcel)
			},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.EqualError(t, err, "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			store := ParcelStore{db: db}
			tt.mocks(dbMock)

			parcel, err := store.Get(tt.number)
			tt.wantErr(t, err)
			tt.wantParcel(t, parcel)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestGetByClient(t *testing.T) {
	t.Parallel()

	type args struct {
		client int
	}

	tests := []struct {
		name        string
		mocks       func(dbMock sqlmock.Sqlmock, client int)
		args        args
		wantParcels require.ValueAssertionFunc
		wantErr     require.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				client: 102,
			},
			mocks: func(dbMock sqlmock.Sqlmock, client int) {
				rows := sqlmock.NewRows([]string{"number", "client", "status", "address", "created_at"}).
					AddRow(101, 102, "Registered", "Address 1", "2023-11-20T10:00:00Z").
					AddRow(102, 102, "Delivered", "Address 2", "2023-11-21T11:00:00Z")
				dbMock.ExpectQuery("SELECT number, client, status, address, created_at FROM percel WHERE client = ?").
					WithArgs(client).
					WillReturnRows(rows)
			},
			wantParcels: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcels, ok := got.([]Parcel)
				require.True(tt, ok)
				require.Len(tt, parcels, 2)
				require.Equal(tt, 101, parcels[0].Number)
				require.Equal(tt, 102, parcels[0].Client)
				require.Equal(tt, "Registered", parcels[0].Status)
				require.Equal(tt, "Address 1", parcels[0].Address)
				require.Equal(tt, "2023-11-20T10:00:00Z", parcels[0].CreatedAt)

				require.Equal(tt, 102, parcels[1].Number)
				require.Equal(tt, 102, parcels[1].Client)
				require.Equal(tt, "Delivered", parcels[1].Status)
				require.Equal(tt, "Address 2", parcels[1].Address)
				require.Equal(tt, "2023-11-21T11:00:00Z", parcels[1].CreatedAt)
			},
			wantErr: require.NoError,
		},
		{
			name: "no records",
			args: args{
				client: 103,
			},
			mocks: func(dbMock sqlmock.Sqlmock, client int) {
				rows := sqlmock.NewRows([]string{"number", "client", "status", "address", "created_at"})
				dbMock.ExpectQuery("SELECT number, client, status, address, created_at FROM percel WHERE client = ?").
					WithArgs(client).
					WillReturnRows(rows)
			},
			wantParcels: func(tt require.TestingT, got interface{}, i ...interface{}) {
				parcels, ok := got.([]Parcel)
				require.True(tt, ok)
				require.Empty(tt, parcels)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			args: args{
				client: 104,
			},
			mocks: func(dbMock sqlmock.Sqlmock, client int) {
				dbMock.ExpectQuery("SELECT number, client, status, address, created_at FROM percel WHERE client = ?").
					WithArgs(client).
					WillReturnError(errors.New("database error"))
			},
			wantParcels: func(tt require.TestingT, got interface{}, i ...interface{}) {
				require.Nil(tt, got)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			store := ParcelStore{db: db}
			tt.mocks(dbMock, tt.args.client)

			parcels, err := store.GetByClient(tt.args.client)
			tt.wantErr(t, err)
			tt.wantParcels(t, parcels)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestSetStatus(t *testing.T) {
	t.Parallel()

	type args struct {
		number int
		status string
	}

	tests := []struct {
		name    string
		mocks   func(dbMock sqlmock.Sqlmock, number int, status string)
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				number: 101,
				status: "Delivered",
			},
			mocks: func(dbMock sqlmock.Sqlmock, number int, status string) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("UPDATE parcel SET status = ? WHERE number = ?")).
					WithArgs(status, number).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows affected",
			args: args{
				number: 999,
				status: "Delivered",
			},
			mocks: func(dbMock sqlmock.Sqlmock, number int, status string) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("UPDATE parcel SET status = ? WHERE number = ?")).
					WithArgs(status, number).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.NoError(tt, err, i...)
			},
		},
		{
			name: "database error",
			args: args{
				number: 101,
				status: "Delivered",
			},
			mocks: func(dbMock sqlmock.Sqlmock, number int, status string) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("UPDATE parcel SET status = ? WHERE number = ?")).
					WithArgs(status, number).
					WillReturnError(errors.New("database error"))
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "database error", i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			store := ParcelStore{db: db}
			tt.mocks(dbMock, tt.args.number, tt.args.status)

			err = store.SetStatus(tt.args.number, tt.args.status)
			tt.wantErr(t, err)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestSetAddress(t *testing.T) {
	t.Parallel()

	type args struct {
		number  int
		address string
	}

	tests := []struct {
		name    string
		mocks   func(dbMock sqlmock.Sqlmock)
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("UPDATE parcel SET address = ? WHERE number = ?")).
					WithArgs("new address", 101).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				number:  101,
				address: "new address",
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("UPDATE parcel SET address = ? WHERE number = ?")).
					WithArgs("new address", 101).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				number:  101,
				address: "new address",
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "database error", i...)
			},
		},
		{
			name: "no rows affected",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("UPDATE parcel SET address = ? WHERE number = ?")).
					WithArgs("new address", 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args: args{
				number:  999,
				address: "new address",
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.NoError(tt, err, i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			store := NewParcelStore(db)
			tt.mocks(dbMock)

			err = store.SetAddress(tt.args.number, tt.args.address)
			tt.wantErr(t, err)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		number int
	}

	tests := []struct {
		name    string
		mocks   func(dbMock sqlmock.Sqlmock)
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM parcel WHERE number = ? AND status = registered")).
					WithArgs(101).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				number: 101,
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM parcel WHERE number = ? AND status = registered")).
					WithArgs(101).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				number: 101,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "database error", i...)
			},
		},
		{
			name: "no rows affected",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM parcel WHERE number = ? AND status = registered")).
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args: args{
				number: 999,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.NoError(tt, err, i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			store := NewParcelStore(db)
			tt.mocks(dbMock)

			err = store.Delete(tt.args.number)
			tt.wantErr(t, err)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}
