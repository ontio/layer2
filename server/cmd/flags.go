package cmd

import (
	"github.com/ontio/layer2/server/config"
	"github.com/urfave/cli"
	"strings"
)

var (
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: uint(config.DefConfig.LogLevel),
	}
	RestPortFlag = cli.UintFlag{
		Name:  "restport",
		Usage: "restful server listening port `<number>`",
		Value: 0,
	}
	ProjectIdFlag = cli.UintFlag{
		Name:  "projectid",
		Usage: "airdrop project id `<number>`",
		Value: 0,
	}
)

func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
