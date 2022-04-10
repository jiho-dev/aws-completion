package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/tidwall/gjson"
)

func RunCmd(inCmds []string, apiArgs []string, adminVpc bool, flags *flag.FlagSet) (string, error) {
	profile, err := flags.GetString(CMD_PROFILE)
	if err != nil {
		return "", err
	}

	cmd := inCmds[0]

	var cmdOpt []string

	cmdOpt = append(cmdOpt, "ec2")
	cmdOpt = append(cmdOpt, "--profile")
	cmdOpt = append(cmdOpt, profile)

	if adminVpc {
		cmdOpt = append(cmdOpt, "admin-vpc")
		if cmd == CMD_SHOW_HELP {
			cmd = "--h"
		} else {
			cmdOpt = append(cmdOpt, "--admin-action")
		}
	} else if cmd == CMD_SHOW_HELP {
		return ShowEc2Cmd(), nil
	}

	cmdOpt = append(cmdOpt, cmd)

	subShowHelp, _ := flags.GetBool(CMD_SHOW_HELP)
	if subShowHelp {
		cmdOpt = append(cmdOpt, "--h")
	} else {
		var optCnt int
		for _, o := range apiArgs {
			if v, err := flags.GetString(o); v != "" && err == nil {
				if adminVpc {
					if optCnt == 0 {
						cmdOpt = append(cmdOpt, "--parameters")
					}

					cmdOpt = append(cmdOpt, fmt.Sprintf("Name=%s,Values=%v", o, v))
					optCnt++
				} else {
					cmdOpt = append(cmdOpt, fmt.Sprintf("--%s", o))
					cmdOpt = append(cmdOpt, v)
					optCnt += 2
				}
			}
		}
	}

	return ExecuteAwsCli("aws", cmdOpt...)
}

func ExecuteAwsCli(name string, args ...string) (string, error) {
	s := name
	if len(args) > 0 {
		for _, a := range args {
			s += " " + a
		}
	}

	fmt.Println(">", s)
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()

	o := string(out)
	if err != nil {
		fmt.Println("> error:", err.Error())
		fmt.Printf("> output: %s\n", o)
		return "", err
	}

	return o, nil
}

func ParseOutput(output string, outField string) string {
	value := gjson.Get(output, outField)
	output = value.String()

	return output
}

func FormatJson(output string) string {
	var obj map[string]interface{}
	json.Unmarshal([]byte(output), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.KeyColor = color.New(color.FgBlue)
	f.Indent = 3
	f.RawStrings = true

	// Marshall the Colorized JSON
	b, _ := f.Marshal(obj)

	return string(b)
}

func ShowEc2Cmd() string {
	var allEc2Cmds []string

	allEc2Cmds = GetEc2Apis(nil)

	/*
		out, _ := ExecuteAwsCli("aws", "ec2", "help")
		l := strings.Split(out, "\n")
		for _, cmd := range l {
			if strings.Index(cmd, "o ") < 0 {
				continue
			}

			cmd = strings.TrimSpace(cmd)
			if len(cmd) < 1 {
				continue
			}

			cc := strings.Split(cmd, " ")
			if len(cc) < 2 {
				continue
			}

			allEc2Cmds = append(allEc2Cmds, cc[1])
		}
	*/

	//cmds := strings.Join(allEc2Cmds, "\n")
	for _, api := range allEc2Cmds {
		fmt.Printf("  %s \n", api)
	}

	return ""
}

func ShowEc2AdminVpc() {
	fmt.Printf("Admin-Vpc subcommands: \n")

	for api, _ := range AdminVpcCmds {
		fmt.Printf("  %s\n", api)
	}
}

func ShowApiMain(cobraCmd *cobra.Command, args []string) {
	switch cobraCmd.Use {
	case CMD_SHOW_EC2_CMDS:
		ShowEc2Cmd()

	case CMD_SHOW_ADMIN_VPC_CMDS:
		ShowEc2AdminVpc()
	}
}
