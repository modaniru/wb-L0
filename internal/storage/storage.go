package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	NotFoundOrderError  = errors.New("no found order")
	OrderAlreadyInCache = errors.New("order already in cache, so skip")
)

type OrderStorage struct {
	db    *sql.DB
	cache CacheI
}

func NewOrderStorage(db *sql.DB, cache CacheI) *OrderStorage {
	_, err := db.Exec(`create table if not exists orders (
		order_uid varchar primary key,
		order_json JSONB
	);`)
	if err != nil{
		log.Fatal(err.Error())
	}
	storage := &OrderStorage{db: db, cache: cache}
	storage.fillCache()
	return storage
}

func (o *OrderStorage) SaveOrder(ctx context.Context, orderId string, order []byte) error {
	if _, ok := o.cache.Get(orderId); ok {
		return OrderAlreadyInCache
	}
	query := "insert into orders (order_uid, order_json) values ($1, $2) returning order_json;"
	stmt, err := o.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare sql error: %w", err)
	}
	row := stmt.QueryRow(orderId, order)

	res := []byte{}
	err = row.Scan(&res)
	if err != nil {
		return fmt.Errorf("scan query error: %w", err)
	}
	o.cache.Put(orderId, res)
	return nil
}

func (o *OrderStorage) GetByUid(ctx context.Context, orderId string) ([]byte, error) {
	res, ok := o.cache.Get(orderId)
	if !ok {
		return nil, NotFoundOrderError
	}
	return res, nil
}

type rowsCount struct{
	Cache int
	DB int
}

func (o *OrderStorage) GetRowsCount(ctx context.Context) (*rowsCount, error){
	rc := rowsCount{}
	rc.Cache = o.cache.Len()


	query := "select count(*) as count from orders;"
	stmt, err := o.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare sql error: %w", err)
	}

	row := stmt.QueryRow()
	err = row.Scan(&rc.DB)
	if err != nil {
		return nil, fmt.Errorf("prepare sql error: %w", err)
	}

	return &rc, nil
}

func (o *OrderStorage) fillCache() {
	query := "select order_uid, order_json from orders;"
	stmt, err := o.db.Prepare(query)
	if err != nil {
		log.Fatal(fmt.Errorf("prepare sql error: %w", err).Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(fmt.Errorf("exec sql query error: %w", err).Error())
	}
	defer rows.Close()

	for rows.Next() {
		key := ""
		data := []byte{}
		err = rows.Scan(&key, &data)
		if err != nil {
			log.Fatal(fmt.Errorf("scan query result error: %w", err).Error())
		}

		o.cache.Put(key, data)
	}
}
