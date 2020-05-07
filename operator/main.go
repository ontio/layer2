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
	"flag"
	"fmt"
	"github.com/ontio/layer2/operator/cmd"
	"github.com/ontio/layer2/operator/config"
	"github.com/ontio/layer2/operator/core"
	"github.com/ontio/layer2/operator/log"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var ConfigPath string

func init() {
	flag.StringVar(&ConfigPath, "cfg", config.DEFAULT_CONFIG_FILE_NAME, "config file of server")
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "Ontology Layer2 Service"
	app.Action = startServer
	app.Version = config.Version
	app.Copyright = "Copyright in 2019 The Ontology Authors"
	app.Flags = []cli.Flag{
		cmd.LogLevelFlag,
		cmd.ConfigPathFlag,
	}
	app.Commands = []cli.Command{
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func startServer(ctx *cli.Context) {
	// get all cmd flag
	logLevel := ctx.GlobalInt(cmd.GetFlagName(cmd.LogLevelFlag))
	log.InitLog(logLevel, log.PATH, log.Stdout)

	configPath := ctx.GlobalString(cmd.GetFlagName(cmd.ConfigPathFlag))
	if configPath != "" {
		ConfigPath = configPath
	}

	// read config
	servConfig := config.NewServiceConfig(ConfigPath)
	if servConfig == nil {
		log.Errorf("startServer - create config failed!")
		return
	}

	initOperatorServer(servConfig)
	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("waitToExit - Layer2 Operator received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}

func initOperatorServer(servConfig *config.ServiceConfig) {
	mgr, err := core.NewLayer2Operator(servConfig)
	if err != nil {
		log.Error("initOperatorServer failed!")
		return
	}
	mgr.Start()
}

func main() {
	log.Infof("main - Layer2 Operator Starting...")
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
