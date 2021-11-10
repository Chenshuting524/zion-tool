/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package flag

import (
	"github.com/urfave/cli"
)

var (
	ConfigPathFlag = cli.StringFlag{
		Name:  "config",
		Usage: "server config file `<path>`",
		Value: "./config.json",
	}

	NumberFlag = cli.Uint64Flag{
		Name:  "num",
		Usage: "test instance number, e.g: user number in tps/stable testing",
		Value: 1,
	}

	PeriodFlag = cli.StringFlag{
		Name:  "period",
		Usage: "set lasting time, e.g: 1d, 1d2h, 1d2h30m40s",
	}

	TxPerPeriod = cli.Uint64Flag{
		Name:  "txn",
		Usage: "set tx number per period",
		Value: 10,
	}

	IncrGasPrice = cli.Uint64Flag{
		Name:  "inc",
		Usage: "gas price increase n wei",
		Value: 0,
	}

	BlockNumberFlag = cli.Uint64Flag{
		Name:  "blocknum",
		Usage: "use for Blocknum",
		Value: 0,
	}
	BlockHashFlag = cli.StringFlag{
		Name:  "blockhash",
		Usage: "query block",
		Value: "",
	}
	TransactionHashFlag = cli.StringFlag{
		Name:  "transactionhash",
		Usage: "query transaction hash",
		Value: "",
	}
)
