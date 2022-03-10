package main

import "time"

func main() {
	InitNet()
	InitRpc()
	//创建账户
	createAccount(5)
	//获取账户列表
	accounts, _ := getAccounts(client)
	//初始化账户余额
	setAccountBalance(accounts)
	//监控账户信息
	go monitorBalance(accounts)
	time.Sleep(600 * time.Second)

}
