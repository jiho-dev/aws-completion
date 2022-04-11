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
var awscConfigName = EXEC_NAME + ".yaml"
var AwscConf *config.AwscConfig

func init() {
	//syslogger, err := syslog.New(syslog.LOG_INFO, EXEC_NAME)
	syslogger, err := syslog.New(syslog.LOG_WARNING, EXEC_NAME)
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
	log.Printf("Start %s", EXEC_NAME)

	if complete() {
		return nil
	}

	var options []string
	fullCmd, optIndex := getFullCommandName()
	if optIndex > 0 {
		options = os.Args[optIndex:]
	}

	//log.Printf("fullCmd:%s, optIndex: %d, options: %v", fullCmd, optIndex, options)

	switch fullCmd {
	case CMD_GENERATE_EC2_CMDS:
		flag := flag.NewFlagSet(EXEC_NAME, flag.ExitOnError)
		flag.String(CMD_PROFILE, "", "AWS Profile")

		if optIndex > 1 {
			flag.Parse(options)
		}

		profile, err := flag.GetString(CMD_PROFILE)
		if err != nil {
			fmt.Printf("Profile error: %v\n", err)
			return nil
		} else if len(profile) < 1 {
			fmt.Printf("%s %s: Need admin profile \n", EXEC_NAME, CMD_GENERATE_EC2_CMDS)
			return nil
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
