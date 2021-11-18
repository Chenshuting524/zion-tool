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
	roundNo := 0
	timeout := time.After(time.Second * time.Duration(period))
	for loop := true; loop; {
		select {
		case <-timeout:
			fmt.Println("timeout 1")
			loop = false
		default:
			var wg sync.WaitGroup
			for _, acc := range accounts {
				wg.Add(1)
				go func(acc *sdk.Account, to common.Address, txn int, period int) {
					hashlist := sendTransfer(acc, to, txn)
					//发完交易之后,开始遍历hash,查询交易是否全部落账
					/*
						for i := range hashlist {
							log.Info("query transaction status")
							fmt.Println("query transaction status")
							//fmt.Println(hashlist[i])
							err = WaitTxConfirm(acc, hashlist[i], period)
							if err != nil {
								fmt.Println("error", err)
								continue
							}
						}*/

					//只查看最后一笔落账没
					for {
						err = WaitTxConfirm(acc, hashlist[len(hashlist)-1], period)
						if err != nil {
							fmt.Println("error", err)
							continue
						} else {
							break
						}
					}
					//fmt.Println("round1111")
					defer wg.Done()
					//等待所有线程结束后开启新一轮
				}(acc, to, txn, period)
			}
			wg.Wait()
		}
		roundNo += 1
		fmt.Println("round", roundNo, " finish", time.Now())
	}

	/*
		for !flag {
			println("start multi process")
			var wg sync.WaitGroup
			for _, acc := range accounts {
				wg.Add(1)
				go func(acc *sdk.Account, to common.Address, txn int, period int) {
					hashlist := sendTransfer(acc, to, txn)
					//发完交易之后,开始遍历hash,查询交易是否全部落账
					for i := range hashlist {
						log.Info("query transaction status")
						fmt.Println("query transaction status")
						//fmt.Println(hashlist[i])
						err = WaitTxConfirm(acc, hashlist[i], period)
						if err != nil {
							fmt.Println("error")
							continue
						}
					}
					//fmt.Println("round1111")
					defer wg.Done()
					//等待所有线程结束后开启新一轮
				}(acc, to, txn, period)
			}
			wg.Wait()
			fmt.Println("round",roundNo+1," finish")

			select {
			case <-time.After(time.Second * time.Duration(period)):
				fmt.Println("timeout 1")
				flag = true
			default:
				fmt.Println("continue")
				flag = false
			}

		}*/
	fmt.Println("finish", time.Now())
	return nil
}

func sendTransfer(acc *sdk.Account, to common.Address, txn int) []common.Hash {
	hashlist := make([]common.Hash, 0)
	for i := 0; i < txn; i++ {
		txhash, err := acc.Transfer(to, amountPerTx)
		if err != nil {
			fmt.Println("transfer failed", "err", err)
		} else {
			//发送成功，将hash保存下来
			hashlist = append(hashlist, txhash)
			//fmt.Println("transfer success")
		}
	}
	return hashlist
}

func WaitTxConfirm(acc *sdk.Account, hash common.Hash, period int) error {
	ticker := time.NewTicker(time.Millisecond * 1)
	end := time.Now().Add(time.Duration(period))
	for now := range ticker.C {
		//fmt.Println("START")
		_, pending, err := acc.TransactionByHash(hash)
		if err != nil {
			log.Info("failed to call TransactionByHash: %v", err)
			if now.After(end) {
				break
			}
			continue
		}
		if !pending {
			//fmt.Println("to get receipt", hash)
			break
		} else {
			fmt.Println("pending", hash)
			if now.After(end) {
				break
			}
			continue
		}
	}
	for now_2 := range ticker.C {
		tx, err := acc.TransactionReceipt(hash)
		end_wait := time.Now().Add(time.Duration(period))
		if err != nil {
			fmt.Println("failed to get receipt", hash.Hex())
			if now_2.After(end_wait) {
				return fmt.Errorf("failed to get receipt %s", hash.Hex())
			} else {
				continue
			}
		} else {
			if tx.Status == 0 {
				return fmt.Errorf("receipt failed %s", hash.Hex())
			} else {
				fmt.Println("confirm", hash)
				break
			}

		}
	}

	return nil
}
