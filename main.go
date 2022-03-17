package main

import (
	"sync"
)

func main() {
	//InitNet()
	wg := sync.WaitGroup{}
	InitRpc()
	//创建账户
	//createAccount(10)
	//获取账户列表
	accounts, _ := getAccounts(client)
	//初始化账户余额
	//setAccountBalance(accounts)
	//监控账户信息
	//go monitorBalance(accounts)
	//unlockAllAccounts(accounts)
	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		rpcPerformance(accounts, client, 100)
	}()
	wg.Wait()
	//交易总数/花费时间
	//sendTxTest(client, accounts)

	//go monitorTxpool(client)
	//go minerStart(client, 1)
}
