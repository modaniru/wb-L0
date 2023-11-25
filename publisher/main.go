package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/modaniru/wb-L0/internal/entity"
	"github.com/nats-io/stan.go"
)

func main() {
	conn, _ := stan.Connect("prod", "producer")
	wg := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			address := gofakeit.Address()
			track := "WB" + gofakeit.UUID()
			uuid := gofakeit.UUID()
			order := entity.Order{
				OrderUid:          uuid,
				TrackNumber:       track,
				Entry:             "WBIL",
				Locale:            "EN",
				InternalSignature: "",
				CustomerId:        "test",
				DeliveryService:   "meest",
				ShardKey:          "9",
				SmId:              99,
				DateCreated:       gofakeit.FutureDate().String(),
				OffShard:          "1",
				Delivery: entity.Delivery{
					Name:    "Test Testov",
					Phone:   strconv.Itoa(gofakeit.Number(79000000000, 79999999999)),
					Zip:     address.Zip,
					City:    address.City,
					Address: address.Address,
					Region:  address.State,
					Email:   gofakeit.Email(),
				},
				Payment: entity.Payment{
					Transaction:  gofakeit.UUID(),
					RequestId:    "",
					Currency:     gofakeit.CurrencyShort(),
					Provider:     "wbpay",
					Amount:       gofakeit.IntRange(1000, 10000),
					Bank:         gofakeit.RandomString([]string{"sber", "alpha", "tinkoff"}),
					DeliveryCost: gofakeit.IntRange(100, 500),
					GoodsTotal:   317,
					CustomFee:    0,
				},
				Items: []entity.Item{
					{
						ChrtId:      gofakeit.IntRange(9000000, 9999999),
						TrackNumber: track,
						Price:       gofakeit.IntRange(100, 5000),
						Rid:         gofakeit.UUID(),
						Name:        gofakeit.Fruit(),
						Sale:        gofakeit.IntRange(1, 10),
						TotalPrice:  gofakeit.IntRange(100, 5000),
						NmId:        gofakeit.IntRange(100000, 999999),
						Brand:       "test brand",
						Status:      202,
					},
				},
			}
			data, _ := json.Marshal(order)
			err := conn.Publish("test", data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(uuid)
			wg.Done()
		}()
	}
	wg.Wait()
}
