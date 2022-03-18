package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"strconv"
	"sync"
	"time"
)

func rpcPerformance(accounts []string, client *rpc.Client, txPool_nums int, qps int, thread_num int) {
	//压力测试
	//txPool_nums:交易池中的交易数，qps：每秒向交易池发送的交易，thread_num：挖矿时使用的线程数
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
			qps,
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
	fmt.Println("共发送了", tx_sum, "笔交易......开始进行压力测试")

	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Add(-1)
		monitorTxpool(client, tx_sum, thread_num)
	}()
	wg1.Wait()

	fmt.Println("压力测试结束")
}

//监控交易池，交易池中交易数目达到txNumber后开始挖矿
func monitorTxpool(client *rpc.Client, txNumber int, thread_num int) {
	ticker1 := time.NewTicker(5 * time.Second) //打印日志计时器
	isMiner := false
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		executionSummary(ch, txNumber)
	}()
	for {
		txpoolInformation, err := txpool_status(client)
		if err != nil {
			fmt.Println("交易池监听失败 err is", err)
		} else {
			pending, _ := strconv.ParseInt(txpoolInformation.Pending, 0, 0)
			queued, _ := strconv.ParseInt(txpoolInformation.Queued, 0, 0)
			select {
			case <-ticker1.C:
				fmt.Println("监听成功:Pending is ", pending, "......Queued is ", queued)
			default:
				if int(pending) >= txNumber && isMiner == false {
					//交易池蓄满，开始打包上链发送交易
					ch <- 0
					isMiner, _ = minerStart(client, thread_num)
					fmt.Println("开始打包上链,isMiner", isMiner)
				}
				if pending == 0 && isMiner == true {
					//交易执行完成，结束挖矿
					isMiner, _ = minerStop(client) //isMiner = fale
					fmt.Println("交易执行完成停止挖矿,isMiner", isMiner)
					ch <- 1
					goto END
				}
			}
		}
	}
END:
	wg.Wait()
}
