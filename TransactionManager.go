package main

import (
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/time/rate"
	"math/rand"
	"sync"
	"time"
)

type MultiTransaction struct {
	Client  string
	Servers []string
	Qps     int
}

type TransactionAccounts struct {
	Clients []string
	Servers []string
}

func generateTxAccounts(accounts []string) TransactionAccounts {
	//生成交易发送双方的账户，偶数为发送方，奇数为接收方
	var result TransactionAccounts
	for i, k := range accounts {
		if i%2 == 0 {
			result.Clients = append(result.Clients, k)
		} else {
			result.Servers = append(result.Servers, k)
		}
	}
	return result
}

func SendTx(body []byte) error {
	fmt.Println(string(body))
	return nil
}

func (mx *MultiTransaction) Start(ch chan struct{}) {
	var wg sync.WaitGroup
	pool, err := ants.NewPool(100)

	if err != nil {
		fmt.Println(err)
	}
	limiter := rate.NewLimiter(rate.Every(time.Duration(int(time.Second)/mx.Qps)), mx.Qps)

	for {
		select {
		case <-ch:
			wg.Wait()
			pool.Release()
			fmt.Printf("%s closed\n", mx.Client)
			break
		default:
			if limiter.Allow() {
				wg.Add(1)

				body := make(map[string]string)
				body["from"] = mx.Client
				body["to"] = mx.Servers[rand.Intn(len(mx.Servers))]

				byteBody, _ := json.Marshal(body)

				err := pool.Submit(func() {
					err := SendTx(byteBody)
					if err != nil {
						fmt.Println(err)
					}
					wg.Done()
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
