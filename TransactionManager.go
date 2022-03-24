package main

import (
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

type txpoolStatus struct {
	Pending string `json:"pending"`
	Queued  string `json:"queued"`
}

//交易负载生成
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
	//fmt.Println("构造的参数为", args)

	err = client.Call(&txHash, "eth_sendTransaction", args)
	if err != nil {
		fmt.Println("交易发送失败from", from, "to", to)
		return "", err
	} else {
		return txHash, nil
	}
}

//交易池创建
func (mx *MultiTransaction) Start(ch chan struct{}, sum chan int, client *rpc.Client) {
	var wg sync.WaitGroup
	pool, err := ants.NewPool(100)

	if err != nil {
		fmt.Println(err)
	}
	limiter := rate.NewLimiter(rate.Every(time.Duration(int(time.Second)/mx.Qps)), mx.Qps)
FOR:
	for {
		select {
		case <-ch:
			wg.Wait()
			pool.Release()
			break FOR
		default:
			if limiter.Allow() {
				wg.Add(1)
				//10 5发送 5接收
				body := make(map[string]string)
				body["from"] = mx.Client
				body["to"] = mx.Servers[rand.Intn(len(mx.Servers))]

				//byteBody, _ := json.Marshal(body)

				err := pool.Submit(func() {
					txHash, err := sendTransaction(client, body["from"], body["to"])
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("交易发送成功from", body["from"], "to", body["to"], "txHash:", txHash)
						sum <- 1
					}
					//SendTx(byteBody)
					wg.Done()
				})
				if err != nil {
					fmt.Println("交易发送失败", err)
				}
			}
		}
	}
}

//交易池监听
func txpool_status(client *rpc.Client) (txPoolStatus txpoolStatus, err error) {
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
	//fmt.Println("构造的参数为", args)
	err = client.Call(&txPoolStatus, "txpool_status")
	if err != nil {
		return txPoolStatus, err
	} else {
		return txPoolStatus, nil
	}
}
