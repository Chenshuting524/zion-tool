package journal

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

func QueryTransaction(ctx *cli.Context) error {
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
	blocknum, transactionhash := getBlockandTransaction(ctx)
	//查询区块
	if blocknum != 0 {
		block, err := client.BlockByHash(context.Background(), common.HexToHash(transactionhash))
		if err != nil {
			return err
		}
		fmt.Println("blockheader", block.Header())
	}
	//查询交易
	if transactionhash != "" {
		tx, pending, err := client.TransactionByHash(context.Background(), common.HexToHash(transactionhash))
		if err != nil {
			return err
		}
		fmt.Println("input", tx.Data(), "pengdingstatus", pending, "err", err)

		if !pending {

			transactionReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(transactionhash))
			fmt.Println("receipt", transactionReceipt)
			if err != nil {
				return fmt.Errorf("faild to get receipt %s", common.HexToHash(transactionhash).Hex())
			}

			if transactionReceipt.Status == 0 {
				return fmt.Errorf("receipt failed %s", common.HexToHash(transactionhash).Hex())
			}
		}

	}
	fmt.Println("finish query")
	return nil

}
