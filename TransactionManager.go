package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
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

type TransactionArgs struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
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
	fmt.Println("交易发送成功.....", string(body))
	return nil
}

//调用rpc发送交易
func sendTransaction(client *rpc.Client, from string, to string) (txHash string, err error) {
	args := TransactionArgs{
		From:     from,
		To:       to,
		Gas:      "0x76c0",
		GasPrice: "0x9184e72a000",
		Value:    "0x1",
	}
	//args := map[string]string{
	//	"from":     from,
	//	"to":       to,
	//	"gas":      "0x76c0",
	//	"gasPrice": "0x9184e72a000",
	//	"value":    "0x1",
	//}
	//bytes, e := json.Marshal(args)
	//if e != nil {
	//	fmt.Printf("参数序列化失败")
	//	return "", nil
	//}
	//jsonStr := string(bytes)
	fmt.Println("构造的参数为", args)

	err = client.Call(&txHash, "eth_sendTransaction", args)
	if err != nil {
		fmt.Println("交易发送失败from", from, "to", to)
		return "", err
	} else {
		return txHash, nil
	}
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
					//txHash, err := sendTransaction(client, body["from"], body["to"])
					//if err != nil {
					//	fmt.Println(err)
					//} else {
					//	fmt.Println("交易发送成功from", body["from"], "to", body["to"], "txHash:", txHash)
					//}
					SendTx(byteBody)
					wg.Done()
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
