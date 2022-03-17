package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	client *rpc.Client
)

var pids []int

type Config struct {
	ChainID int `json:"chainId"`
}

type Alloc struct {
}

type genesis struct {
	Config     Config `json:"config"`
	Nonce      string `json:"nonce"`
	Timestamp  string `json:"timestamp"`
	ParentHash string `json:"parentHash"`
	ExtraData  string `json:"extraData"`
	GasLimit   string `json:"gasLimit"`
	Difficulty string `json:"difficulty"`
	Mixhash    string `json:"mixhash"`
	Coinbase   string `json:"coinbase"`
	Alloc      Alloc  `json:"alloc"`
}

//创建新账户
func creatNewAccount(client *rpc.Client, password string) (newAccount string, err error) {
	err = client.Call(&newAccount, "personal_newAccount", password)
	if err != nil {
		return "", err
	}
	return newAccount, nil
}

//解锁账户基础操作
func unlockAccount(client *rpc.Client, address string, passphrase string, duration int) (unlock bool, err error) {
	err = client.Call(&unlock, "personal_unlockAccount", address, passphrase, duration)
	if err == nil {
		return unlock, nil
	} else {
		return false, err
	}
}

//获取账户列表
func getAccounts(client *rpc.Client) (accounts []string, err error) {
	err = client.Call(&accounts, "eth_accounts")
	if err == nil {
		return accounts, nil
	} else {
		return nil, errors.New("账户列表获取错误")
	}

}

//获取挖矿账户
func getCoinbase(client *rpc.Client) (coinbase string, err error) {
	err = client.Call(&coinbase, "eth_coinbase")
	if err == nil {
		return coinbase, nil
	} else {
		return "", errors.New("挖矿账户获取错误")
	}
}

//获取余额
func getBalance(client *rpc.Client, account string) (Balance int64, err error) {

	var balance string
	err = client.Call(&balance, "eth_getBalance", account, "latest")
	if err != nil {
		return -1, err
	}
	Balance, _ = strconv.ParseInt(balance, 0, 64)
	return Balance, nil

}

//节点挖矿
func minerStart(client *rpc.Client, thread_num int) (start bool, err error) {
	err = client.Call(&start, "miner_start", thread_num)
	if err == nil {
		return true, nil
	} else {
		return false, errors.New("挖矿启动失败")
	}
}

//指定挖矿收益账户
func miner_setEtherbase(client *rpc.Client, address string) {
	var result bool
	err := client.Call(&result, "miner_setEtherbase", address)
	if err == nil {
		fmt.Println("设置收益账户成功，设置为", address)
	} else {
		fmt.Println("设置失败")
	}
}

//挖矿停止
func minerStop(client *rpc.Client) (stop bool, err error) {
	err = client.Call(&stop, "miner_stop")
	if err == nil {
		return false, nil
	} else {
		return true, errors.New("终止挖矿失败")
	}
}

//挖矿账户选择

func execCommand(commandName string, params []string) bool {

	cmd := exec.Command(commandName, params...)
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

	//显示运行的命令
	fmt.Println(cmd.Args)

	//stdout, err := cmd.StdoutPipe()
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return false
	//}

	cmd.Start()
	var pid = cmd.Process.Pid
	fmt.Println(pid)
	pids = append(pids, pid)
	//
	//reader := bufio.NewReader(stdout)
	//
	////实时循环读取输出流中的一行内容
	//for {
	//	line, err2 := reader.ReadString('\n')
	//	if err2 != nil || io.EOF == err2 {
	//		break
	//	}
	//	fmt.Println(line)
	//}

	//cmd.Wait()
	time.Sleep(10 * time.Second)
	return true
}

//func main() {
//
//	//创建新账户
//	var password string = "123456"
//	newAccount, err := creatNewAccount(client, password)
//	if err != nil {
//		fmt.Println("err=", err)
//	}
//	fmt.Println("新账户为：", newAccount)
//
//	//获取账户列表
//	accounts, err := getAccounts(client)
//	if err != nil {
//		fmt.Println("err=", err)
//	}
//	for i, v := range accounts {
//		balance, err := getBalance(client, v)
//		if err != nil {
//			fmt.Println("err=", err)
//		} else {
//			fmt.Printf("账户%d的账号为：%s，余额为：%d\n", i, v, balance)
//		}
//
//	}
//
//	//获取挖矿账户
//	coinbase, err := getCoinbase(client)
//	if err != nil {
//		fmt.Println("err=", err)
//	}
//	fmt.Println("挖矿账户为：", coinbase)
//
//	//延迟关闭
//	defer client.Close()
//
//}

//修改账户余额
func generateJson(accounts []string) {
	balance := map[string]string{
		"balance": "0x1000000000000000000",
	}
	config := map[string]interface{}{
		"chainId": 10,
	}
	var alloc = make(map[string]interface{})
	for _, v := range accounts {
		alloc[v] = balance
	}

	info := map[string]interface{}{
		"config":     config,
		"nonce":      "0x0000000000000042",
		"timestamp":  "0x0",
		"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"extraData":  "",
		"gasLimit":   "0xffffffff",
		"difficulty": "0xfffff",
		"mixhash":    "0x0000000000000000000000000000000000000000000000000000000000000000",
		"coinbase":   "0x3333333333333333333333333333333333333333",
		"alloc":      alloc,
	}
	bytes, e := json.Marshal(info)
	if e != nil {
		fmt.Printf("序列化失败")
		return
	} else {
		jsonStr := string(bytes)
		fmt.Println(jsonStr)
	}

	filePtr, err := os.Create("genesis.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	//带JSON缩进格式写文件
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println("Generate failed", err.Error())
	} else {
		fmt.Println("Generate success")
	}
	filePtr.Write(data)

}

//func writeFile() {
//	personInfo := []genesis{{"David", 30, true, []string{"跑步", "读书", "看电影"}}, {"Lee", 27, false, []string{"工作", "读书", "看电影"}}}
//
//	// 创建文件
//	filePtr, err := os.Create("person_info.json")
//	if err != nil {
//		fmt.Println("Create file failed", err.Error())
//		return
//	}
//	defer filePtr.Close()
//
//	// 创建Json编码器
//	encoder := json.NewEncoder(filePtr)
//
//	err = encoder.Encode(personInfo)
//	if err != nil {
//		fmt.Println("Encoder failed", err.Error())
//
//	} else {
//		fmt.Println("Encoder success")
//	}
//	// 带JSON缩进格式写文件　　//data, err := json.MarshalIndent(personInfo, "", "  ")   //if err != nil {   // fmt.Println("Encoder failed", err.Error())   //   //} else {   // fmt.Println("Encoder success")   //}   //   //filePtr.Write(data)
//}
