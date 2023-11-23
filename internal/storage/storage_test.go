package storage

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

//todo Mock cache
func TestCraeteStorageAndFillCache(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Error(err.Error())
	}

	orderUid := "2"
	orderJson := []byte("test2")

	mock.ExpectPrepare("select order_uid, order_json from orders;")
	mock.ExpectQuery("select order_uid, order_json from orders;").
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "order_json"}).AddRow(orderUid, orderJson))

	_ = NewOrderStorage(db, NewInmemoryCache())
}


func TestSaveOrder(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Error(err.Error())
	}

	uid := "1"
	data := []byte("test")

	uid2 := "2"
	data2 := []byte("test2")

	mock.ExpectPrepare("select order_uid, order_json from orders;")
	mock.ExpectQuery("select order_uid, order_json from orders;").
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "order_json"}).AddRow(uid2, data2))
	mock.ExpectPrepare(regexp.QuoteMeta("insert into orders (order_uid, order_json) values ($1, $2) returning order_json;"))
	mock.ExpectQuery("insert into orders").
		WithArgs(uid, data).
		WillReturnRows(sqlmock.NewRows([]string{"order_json"}).AddRow(data))

	storage := NewOrderStorage(db, NewInmemoryCache())
	err = storage.SaveOrder(context.Background(), uid, data)
	assert.NoError(t, err)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveOrderAlredyInCacheError(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Error(err.Error())
	}

	uid := "1"
	data := []byte("test")

	mock.ExpectPrepare("select order_uid, order_json from orders;")
	mock.ExpectQuery("select order_uid, order_json from orders;").
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "order_json"}).AddRow(uid, data))

	storage := NewOrderStorage(db, NewInmemoryCache())
	err = storage.SaveOrder(context.Background(), uid, data)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetById(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Error(err.Error())
	}

	uid := "1"
	data := []byte("test")

	mock.ExpectPrepare("select order_uid, order_json from orders;")
	mock.ExpectQuery("select order_uid, order_json from orders;").
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "order_json"}).AddRow(uid, data))

	storage := NewOrderStorage(db, NewInmemoryCache())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	actual, err := storage.GetByUid(context.Background(), uid)
	assert.NoError(t, err)
	assert.Equal(t, data, actual)
}

func TestGetByIdNotFound(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Error(err.Error())
	}

	uid := "1"
	data := []byte("test")

	mock.ExpectPrepare("select order_uid, order_json from orders;")
	mock.ExpectQuery("select order_uid, order_json from orders;").
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "order_json"}).AddRow(uid, data))

	storage := NewOrderStorage(db, NewInmemoryCache())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	actual, err := storage.GetByUid(context.Background(), "2")
	assert.Error(t, err)
	assert.Nil(t, actual)
}