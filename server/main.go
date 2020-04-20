/*
 * Copyright (C) 2020 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"github.com/ontio/layer2/server/cmd"
	"github.com/ontio/layer2/server/config"
	"github.com/ontio/layer2/server/core"
	"github.com/ontio/layer2/server/log"
	"github.com/ontio/layer2/server/restful"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "layer2 server cli"
	app.Action = startLayer2Server
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		//common setting
		cmd.LogLevelFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func startLayer2Server(ctx *cli.Context) {
	//
	logLevel := ctx.GlobalInt(cmd.GetFlagName(cmd.LogLevelFlag))
	logName := fmt.Sprintf("%s%s", "./logout", string(os.PathSeparator))
	log.InitLog(logLevel, logName, log.Stdout)

	//
	err := config.InitConfig()
	if err != nil {
		log.Errorf("initConfig error: %s", err)
		return
	}
	log.Info("Config init success")
	//
	explorer := core.NewExplorer()
	code, err := explorer.Start()
	if code != core.SUCCESS || err != nil {
		log.Errorf("start explorer error: %s", err)
		return
	}
	log.Info("start explorer success")
	//
	rt := restful.InitRestServer()
	go rt.Start()
	log.Infof("restful start success")
	//
	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("cross chain explorer server received exit signal: %s.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
