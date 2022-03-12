package main

import "time"

func rpcPerformance(accounts []string) {
	//压力测试
	txAccounts := generateTxAccounts(accounts)
	clients := txAccounts.Clients
	servers := txAccounts.Servers
	chans := make([]chan struct{}, 0)

	for _, client := range clients {
		ch := make(chan struct{})
		chans = append(chans, ch)
		mx := &MultiTransaction{
			client,
			servers,
			10,
		}
		go mx.Start(ch)
	}

	time.Sleep(time.Second * 10)
	for _, ch := range chans {
		ch <- struct{}{}
	}
	time.Sleep(time.Second * 5)
}
