package cmd

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path"
	"strings"

	"github.com/jiho-dev/aws-completion/config"
	flag "github.com/spf13/pflag"
)

var awsDir = os.Getenv("HOME") + "/.aws/"
var awscConfigName = "ac.yaml"
var AwscConf *config.AwscConfig

func init() {
	syslogger, err := syslog.New(syslog.LOG_INFO, "aws-completion")
	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(syslogger)

	confFile := path.Join(awsDir, awscConfigName)
	conf, err := config.ParseConfig(confFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ParseConfig Err: %v\n", err)
		os.Exit(1)
	}

	AwscConf = conf
}

func getFullCommandName() (string, int) {
	var fullCmd string

	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			return fullCmd, i - 1
		}

		if len(fullCmd) > 0 {
			fullCmd += "-"
		}

		fullCmd += os.Args[i]
	}

	return fullCmd, -1
}

// Execute executes cmd
func Execute() error {

	log.Printf("")
	log.Printf("Start ...")

	if complete() {
		return nil
	}

	fullCmd, optIndex := getFullCommandName()
	log.Printf("fullCmd:%s, optIndex: %d", fullCmd, optIndex)

	switch fullCmd {
	case CMD_GENERATE_EC2_CMDS:
		flag := flag.NewFlagSet("ac", flag.ExitOnError)
		flag.String(CMD_PROFILE, "", "AWS Profile")
		if optIndex > 1 {
			flag.Parse(os.Args[optIndex:])
		}

		GenerateApiMain(flag)
		return nil

	case CMD_SHOW_EC2_CMDS:
		ShowEc2Cmd()
		return nil

	case CMD_SHOW_ADMIN_VPC_CMDS:
		ShowEc2AdminVpc()
		return nil
	}

	if optIndex == -1 {
		fmt.Printf("Unknown command: %s\n", os.Args)
		return nil
	}
	apiMainExt(fullCmd, os.Args[optIndex:])
	return nil
}
