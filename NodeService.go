package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"os"
	"time"
)

//创建节点
func createAccount(account_num int, client *rpc.Client) {
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

//查询账户余额
func monitorBalance(address []string, client *rpc.Client) {
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

//解锁全部账户
func unlockAllAccounts(accounts []string, client *rpc.Client) {
	for _, v := range accounts {
		_, err := unlockAccount(client, v, "123456", 300)
		if err != nil {
			fmt.Println("解锁失败", v, " err is", err)
		} else {
			fmt.Println(v, " 解锁成功")
		}
	}
}

//设置初始节点足够的钱
func setAccountBalance(accounts []string) {
	generateJson(accounts)
	if pids != nil {
		for pid := range pids {
			pro, err := os.FindProcess(pid)
			if err == nil {
				err1 := pro.Kill()
				if err1 != nil {
					fmt.Println("kill failure", err)
				} else {
					fmt.Println("kill successful", pid)
				}
			} else {
				fmt.Println("find failure", err)
			}
		}
		pids = pids[0:0]
	}
	err := os.RemoveAll("./data0/geth")
	if err != nil {
		log.Fatal(err)
	}
	command := "geth"
	args := []string{"init", "--datadir", "data0", "genesis.json"}
	execCommand(command, args)
	args = []string{"--datadir", "data0", "--nodiscover", "--networkid", "10", "--http", "--http.api", "personal,eth,net,web3,admin,miner,txpool"}
	execCommand(command, args)
	fmt.Println("金额设置成功，尝试重新监听rpc端口")
	InitRpc()
}
