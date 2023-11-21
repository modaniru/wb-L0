package natssub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"wb-l0/internal/entity"
	"wb-l0/internal/storage"

	log "log/slog"

	"github.com/nats-io/stan.go"
)

type Subscriber struct{
	storage *storage.OrderStorage
}

func NewSubscriber(storage *storage.OrderStorage) *Subscriber{
	return &Subscriber{storage: storage}
}
func (s *Subscriber) GetMsgHandler() stan.MsgHandler{
	return func(msg *stan.Msg) {
		data := msg.Data
		order := entity.Order{}
		err := json.Unmarshal(data, &order)
		if err != nil{
			log.Error(fmt.Sprintf("unmarshal order error"))
			return
		}
		err = s.storage.SaveOrder(context.Background(), order.OrderUid, data)
		if err != nil{
			if errors.Is(err, storage.OrderAlreadyInCache){
				log.Warn(err.Error())
				return
			}
			log.Error(err.Error())
			return
		}
		log.Info("order saved")
	}
}