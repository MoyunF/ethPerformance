package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)

//创建网络
func InitNet() {
	command := "geth"
	args := []string{"init", "--datadir", "data0", "genesis.json"}
	execCommand(command, args)

	//geth --datadir data0 --nodiscover --networkid 10 --http --http.api personal,eth,net,web3,admin,miner,txpool
	args = []string{"--datadir", "data0", "--nodiscover", "--networkid", "10", "--http", "--http.api", "personal,eth,net,web3,admin,miner,txpool"}
	execCommand(command, args)
}

func InitRpc() {
	//获取连接与eth客户端
	client, _ = rpc.Dial("http://localhost:8545")
	if client == nil {
		fmt.Println("rpc.Dial err")
		//panic("连接错误")
		return
	} else {
		fmt.Println("connect sucessuful")
	}
}

//计算吞吐量
func executionSummary() {

}
