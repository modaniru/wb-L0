package natssub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "log/slog"

	"github.com/modaniru/wb-L0/internal/entity"
	"github.com/modaniru/wb-L0/internal/storage"
	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	storage *storage.OrderStorage
}

func NewSubscriber(storage *storage.OrderStorage) *Subscriber {
	return &Subscriber{storage: storage}
}

//TODO test
func (s *Subscriber) GetMsgHandler() stan.MsgHandler {
	return func(msg *stan.Msg) {
		data := msg.Data
		order := entity.Order{}
		err := json.Unmarshal(data, &order)
		if err != nil {
			log.Error(fmt.Sprintf("unmarshal order error"))
			return
		}
		err = s.storage.SaveOrder(context.Background(), order.OrderUid, data)
		if err != nil {
			if errors.Is(err, storage.OrderAlreadyInCache) {
				log.Warn(err.Error())
				msg.Ack()
				return
			}
			log.Error(err.Error())
			return
		}
		log.Info("order saved")
		msg.Ack()
	}
}
