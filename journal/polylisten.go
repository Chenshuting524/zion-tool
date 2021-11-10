package journal

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

func PolyChainListen(ctx *cli.Context) error {
	fmt.Println("start to listen", "start", true)
	//获取config
	c, err := getConfig(ctx)
	if err != nil {
		return err
	}
	client, err := ethclient.Dial(c.NodeList[0])
	if err != nil {
		return err
	}
	startBlockNo, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(fmt.Sprintf("try to get start block number failed, err: %v", err))
	} else {
		fmt.Println("start from block", startBlockNo)
	}

	period, txn, err := getPeriodAndTxn(ctx)
	if err != nil {
		return err
	}
	fmt.Println("get period and txn", "period", period, "txn", txn)

	cnt := 0
	curBlockNum := startBlockNo + 1
	startTime, currentTime, preTime := uint64(0), uint64(0), uint64(0)
	fmt.Println("Cycle listen")

	//先查询前一个块的时间，是开始时间和tx总数
	header0, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(startBlockNo))
	preTime = header0.Time
	if err != nil {
		time.Sleep(500 * time.Millisecond)
	}
	txn0, err := client.TransactionCount(context.Background(), header0.Hash())
	if err != nil {
		time.Sleep(500 * time.Millisecond)
	}
	totalTx, accumulativeTx := txn0, txn0
	fmt.Println("startBlock", startBlockNo,startTime)

	//进入循环
	for cnt < period {
	retryHeader:
		header, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(curBlockNum))
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryHeader
		}

		currentTime = header.Time

	retryTxCnt:
		txn, err := client.TransactionCount(context.Background(), header.Hash())
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryTxCnt
		}

	retryPendingTX:
		pendingTxNum, err := client.PendingTransactionCount(context.Background())
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryPendingTX
		}

		if currentTime > startTime {
			tps := totalTx / uint((currentTime - preTime))
			fmt.Println("calculate tps", "currentBlock", curBlockNum-1, "Header time",
				preTime, "endTime", currentTime, "pendingTx NUM", pendingTxNum, "total tx", totalTx, "tps", tps, "accumulative", accumulativeTx)
		}

		preTime = currentTime
		totalTx = txn
		accumulativeTx += txn
		curBlockNum += 1
		cnt += 1
	}

	fmt.Println("finish listen")
	return nil
}
