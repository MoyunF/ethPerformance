package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"strconv"
	"sync"
	"time"
)

func rpcPerformance(accounts []string, client *rpc.Client, txPool_nums int) {
	//压力测试
	txAccounts := generateTxAccounts(accounts) //负载生成模式   10个账户 5个发送交易 5个
	clients := txAccounts.Clients
	servers := txAccounts.Servers
	chans := make([]chan struct{}, 0)
	chan_sum := make(chan int, 100)
	tx_sum := 0 //发送交易计数器
	var wg sync.WaitGroup
	for _, client := range clients {
		ch := make(chan struct{})
		chans = append(chans, ch)
		mx := &MultiTransaction{
			client,
			servers,
			10,
		}
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			mx.Start(ch, chan_sum)
		}()
	}

	for {
		select {
		case tx := <-chan_sum:
			tx_sum += tx
			if tx_sum >= txPool_nums {
				//终止交易发送
				for _, ch := range chans {
					ch <- struct{}{}
				}
				wg.Wait()
				close(chan_sum)
				goto END
			}
		default:
		}
	}
END:
	fmt.Println("共发送了", tx_sum, "笔交易")
}

//监控交易池，交易池中交易数目达到txNumber后开始挖矿
func monitorTxpool(client *rpc.Client, txNumber int64) {
	ticker1 := time.NewTicker(5 * time.Second) //打印日志计时器
	isMiner := false
	for {
		txpoolInformation, err := txpool_status(client)
		if err != nil {
			fmt.Println("交易池监听失败 err is", err)
		} else {
			pending, _ := strconv.ParseInt(txpoolInformation.Pending, 0, 64)
			queued, _ := strconv.ParseInt(txpoolInformation.Queued, 0, 64)
			select {
			case <-ticker1.C:
				fmt.Println("监听成功:Pending is ", pending, "......Queued is ", queued)
			default:
				if pending >= txNumber && isMiner == false {
					//交易池蓄满，开始打包上链发送交易
					isMiner, _ = minerStart(client, 1)
				}
				if pending == 0 && isMiner == true {
					//交易执行完成，结束挖矿
					isMiner, _ = minerStop(client) //isMiner = fale
				}
			}
		}
	}
}
