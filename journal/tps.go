package journal

import (
	"fmt"
	"sync"
	"time"

	"github.com/Chenshuting524/zion-tool/sdk"
	"github.com/Chenshuting524/zion-tool/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"
)

//var orlogger log.Logger

func Init() {
	math.Init(18)
	//orlogger = log.New("Handle TPS", ": ")
}

// HandleTPS try to test hotstuff tps, params nodeList represents multiple ethereum rpc url addresses,
// and num denote that this test will use multi account to send simple transaction
func HandleTPS(ctx *cli.Context) error {
	fmt.Println("start to handle tps", "start", true)

	// load config instance
	c, err := getConfig(ctx)
	if err != nil {
		return err
	}

	// load and try to increase gas price
	setGasPriceIncr(ctx)

	// load period and tx number per period
	fmt.Println("try to get period and txn...")
	period, txn, err := getPeriodAndTxn(ctx)
	if err != nil {
		return err
	}
	fmt.Println("get period and txn", "period", period, "txn", txn)

	// generate master account
	fmt.Println("try to generate master account...")
	master, err := generateMasterAccount(c)
	if err != nil {
		return err
	}
	fmt.Println("generate master account", "period", period, "txn", txn)

	// create account
	fmt.Println("try to generate multi test accounts...")
	instanceNo := getInstanceNumber(ctx)
	accounts, err := generateMultiTestingAccounts(c, instanceNo)
	if err != nil {
		return err
	}
	fmt.Println("generated multi test accounts")

	// prepare balance
	fmt.Println("try to prepare test accounts balance...")
	if err := prepareTestingAccountsBalance(master, accounts, instanceNo, period, txn); err != nil {
		return err
	}
	fmt.Println("prepare test accounts balance success")

	// send transactions continuously
	to := master.Address()
	//while
	end := time.Now().Add(time.Duration(period))
	for {
		var wg sync.WaitGroup
		for _, acc := range accounts {
			wg.Add(1)
			println("start multi process")
			go func(acc *sdk.Account, to common.Address, txn int) {
				hashlist := sendTransfer(acc, to, txn)
				//发完交易之后,开始遍历hash,查询交易是否全部落账
				if hashlist != nil {
					defer wg.Done()
				}
				for i := range hashlist {
					log.Info("query transaction status")
				retryHash:
					_, pending, err := acc.TransactionByHash(*hashlist[i])
					if err != nil {
						log.Info("failed to call TransactionByHash: %v", err)
						goto retryHash
					}
					if !pending {
						continue
					} else {
						goto retryHash
					}
				}

				//等待所有线程结束后开启新一轮
			}(acc, to, txn)
		}
		wg.Wait()
		if time.Now().Before(end) {
			continue
		} else {
			break
		}
	}
	return nil
}

func sendTransfer(acc *sdk.Account, to common.Address, txn int) []*common.Hash {
	hashlist := make([]*common.Hash, txn)
	for i := 0; i < txn; i++ {
		txhash, err := acc.Transfer(to, amountPerTx)
		if err != nil {
			fmt.Println("transfer failed", "err", err)
		} else {
			//发送成功，将hash保存下来
			hashlist = append(hashlist, &txhash)
			//fmt.Println("transfer success", "hash", hash)
		}
	}
	return hashlist
}

func WaitTxConfirm() {

}
