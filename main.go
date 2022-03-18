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
	unlockAllAccounts(accounts)
	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		//txPool_nums:交易池中的交易数，qps：每秒向交易池发送的交易，thread_num：挖矿时使用的线程数
		rpcPerformance(accounts, client, 10000, 10, 1)
	}()
	wg.Wait()

}
