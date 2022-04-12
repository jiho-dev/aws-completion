package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jiho-dev/aws-completion/config"
	"github.com/jiho-dev/aws-completion/log"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var awsDir = os.Getenv("HOME") + "/.aws/"
var awscConfigName = config.EXEC_NAME + ".yaml"
var AwscConf *config.AwscConfig
var logger *logrus.Logger

func init() {
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
	//log.InitLogger(logrus.DebugLevel)
	log.InitLogger(logrus.WarnLevel)
	logger = log.GetLogger()

	logger.Infof("Start %s", config.EXEC_NAME)

	if complete() {
		return nil
	}

	var options []string
	fullCmd, optIndex := getFullCommandName()
	if optIndex > 0 {
		options = os.Args[optIndex:]
	}

	logger.Debugf("fullCmd:%s, optIndex: %d, options: %v", fullCmd, optIndex, options)

	switch fullCmd {
	case config.CMD_GENERATE_EC2_CMDS:
		flag := flag.NewFlagSet(config.EXEC_NAME, flag.ExitOnError)
		flag.String(config.CMD_PROFILE, "", "AWS Profile")

		if optIndex > 1 {
			flag.Parse(options)
		}

		profile, err := flag.GetString(config.CMD_PROFILE)
		if err != nil {
			fmt.Printf("Profile error: %v\n", err)
			return nil
		} else if len(profile) < 1 {
			fmt.Printf("%s %s: Need admin profile \n", config.EXEC_NAME, config.CMD_GENERATE_EC2_CMDS)
			return nil
		}

		GenerateApiMain(flag)
		return nil

	case config.CMD_SHOW_EC2_CMDS:
		ShowEc2Cmd()
		return nil

	case config.CMD_SHOW_ADMIN_VPC_CMDS:
		ShowEc2AdminVpc()
		return nil

	case config.CMD_COMPLETION_ZSH:
		generateZshCompletion(false)
		os.Exit(0)
	}

	if optIndex == -1 {
		fmt.Printf("Unknown command: %s\n", os.Args)
		return nil
	}

	apiMainExt(fullCmd, os.Args[optIndex:])
	return nil
}
