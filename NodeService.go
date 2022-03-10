package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"time"
)

//创建节点
func createAccount(account_num int) {
	for i := 0; i < account_num; i++ {
		var password string = "123456"
		newAccount, err := creatNewAccount(client, password)
		if err != nil {
			fmt.Println("err=", err)
		}
		fmt.Println("新账户为：", newAccount)
	}
}

//为指定节点挖矿 作废
func minerForAccount(client *rpc.Client, thread_num int, address string) {
	miner_setEtherbase(client, address)
	result, err := minerStart(client, thread_num)
	if result == false && err != nil {
		fmt.Println("挖矿失败")
	}
}

//结束指定节点挖矿 作废
func minerStopForAccount(client *rpc.Client, address string) {
	stop, err := minerStop(client)
	if err != nil && stop != true {
		fmt.Println("err=", err)
	} else {
		fmt.Println("结束", address, "挖矿")
	}
}

//查询账户余额当账户全部由前后停止
func monitorBalance(address []string) {
	for {
		for i, v := range address {
			balance, err := getBalance(client, v)
			if err != nil {
				fmt.Println("err=", err)
			} else {
				fmt.Printf("账户%d的账号为：%s，余额为：%d\n", i, v, balance)
			}
		}
		time.Sleep(time.Second * 10)
	}
}
