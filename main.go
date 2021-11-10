package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Chenshuting524/zion-tool/flag"
	"github.com/Chenshuting524/zion-tool/journal"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"
)

var (
	CmdTPS = cli.Command{
		Name:   "tps",
		Usage:  "try to test zion consensus tps.",
		Action: journal.HandleTPS,
		Flags: []cli.Flag{
			flag.ConfigPathFlag,
			flag.NumberFlag,
			flag.TxPerPeriod,
			flag.PeriodFlag,
			flag.IncrGasPrice,
		},
	}

	CmdNativeCall = cli.Command{
		Name:   "native",
		Usage:  "try to test native call",
		Action: journal.HandleNative,
		Flags: []cli.Flag{
			flag.ConfigPathFlag,
		},
	}

	CmdCalculate = cli.Command{
		Name: "listen",
		Usage: "Try to listen and the calculate tps",
		Action: journal.PolyChainListen,
		Flags:[]cli.Flag{
			flag.ConfigPathFlag,
			flag.BlockNumberFlag,
			flag.PeriodFlag,
		},
	}
	CmdQuery = cli.Command{
		Name: "query",
		Usage: "Try to query transaction or block",
		Action: journal.QueryTransaction,
		Flags:[]cli.Flag{
			flag.ConfigPathFlag,
			flag.BlockNumberFlag,
			flag.BlockHashFlag,
			flag.TransactionHashFlag,
		},
	}
)

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "hotstuff test tool"
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2021 The Ontology Authors"
	app.Flags = []cli.Flag{
		flag.ConfigPathFlag,
		flag.NumberFlag,
		flag.TxPerPeriod,
		flag.PeriodFlag,
		flag.IncrGasPrice,
	}
	app.Commands = []cli.Command{
		CmdTPS,
		CmdNativeCall,
		CmdCalculate,
	}
	app.Before = beforeCommands
	app.After = afterCommands
	return app
}

func main() {
	app := setupApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// action execute after commands
func beforeCommands(ctx *cli.Context) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	journal.Init()
	return nil
}

func afterCommands(ctx *cli.Context) error {
	log.Info("\r\n" +
		"\r\n" +
		"\r\n")
	return nil
}
