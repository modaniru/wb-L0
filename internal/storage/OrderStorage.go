package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var(
	NotFoundOrderError = errors.New("no found order")
)

type OrderStorage struct{
	db *sql.DB
	cache *Cache
}

func NewOrderStorage(db *sql.DB, cache *Cache) *OrderStorage{
	storage := &OrderStorage{db: db, cache: cache}
	storage.fillCache()
	return storage
}

func (o *OrderStorage) SaveOrder(ctx context.Context, orderId string, order []byte) error{
	query := "insert into orders (order_uid, order_json) values ($1, $2);"
	stmt, err := o.db.Prepare(query)
	if err != nil{
		return fmt.Errorf("prepare sql error: %w", err)
	}
	_, err = stmt.Exec(orderId, order)
	if err != nil{
		return fmt.Errorf("exec sql error: %w", err)
	}
	return nil
}

func (o *OrderStorage) GetByUid(ctx context.Context, orderId string) ([]byte, error){
	res := o.cache.Get(orderId)
	if res == nil{
		return nil, NotFoundOrderError
	}
	return res, nil
}

func (o *OrderStorage) fillCache(){
	query := "select order_uid, order_json from orders;"
	stmt, err := o.db.Prepare(query)
	if err != nil{
		log.Fatal(fmt.Errorf("prepare sql error: %w", err).Error())
	}
	rows, err := stmt.Query()
	if err != nil{
		log.Fatal(fmt.Errorf("exec sql query error: %w", err).Error())
	}
	defer rows.Close()

	for rows.Next(){
		key := ""
		data := []byte{}
		err = rows.Scan(&key, &data)
		if err != nil{
			log.Fatal(fmt.Errorf("scan query result error: %w", err).Error())
		}
	
		o.cache.Put(key, data)
	}
}