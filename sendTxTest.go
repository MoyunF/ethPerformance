package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)

//调用rpc发送交易
func sendTxTest(client *rpc.Client, accounts []string) (txHash string, err error) {
	args := TransactionArgs{
		From:     accounts[0],
		To:       accounts[1],
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
		fmt.Println("交易发送失败from", accounts[0], "to", accounts[1])
		fmt.Println("err is", err)
		return "", err
	} else {
		return txHash, nil
	}
}
