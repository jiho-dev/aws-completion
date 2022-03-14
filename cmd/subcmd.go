package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jiho-dev/aws-completion/config"
	"github.com/spf13/cobra"
)

func InitApiParams(fullApi string, cmd *cobra.Command) {
	opt, ok := AwscConf.ApiOptions[fullApi]
	if !ok {
		return
	}

	for _, o := range opt.Required {
		cmd.Flags().String(o, "", "")
		cmd.MarkFlagRequired(o)
	}

	for _, o := range opt.Args {
		cmd.Flags().String(o, "", "")
	}

	addProfileCmd(cmd)
	cmd.Flags().Bool(CMD_SHOW_HELP, false, "")
}

func InitApiCmd(api string, sub map[string][]string) *cobra.Command {
	cmd1 := &cobra.Command{
		Use:               api,
		Run:               apiMain,
		CompletionOptions: CompOpt,
	}

	for api2, sub2 := range sub {
		if api2 == config.API_TERMINATED {
			InitApiParams(api, cmd1)
			continue
		}

		cmd2 := &cobra.Command{
			Use:               api2,
			Run:               cmd1.Run,
			ValidArgsFunction: getApiArgs,
			CompletionOptions: CompOpt,
		}

		for _, api3 := range sub2 {
			if api3 == config.API_TERMINATED {
				fullApi := api + "-" + api2
				InitApiParams(fullApi, cmd2)
			} else {
				cmd3 := &cobra.Command{
					Use:               api3,
					Run:               cmd1.Run,
					ValidArgsFunction: getApiArgs,
					CompletionOptions: CompOpt,
				}

				fullApi := api + "-" + api2 + "-" + api3
				InitApiParams(fullApi, cmd3)

				cmd2.AddCommand(cmd3)
			}
		}

		cmd1.AddCommand(cmd2)
	}

	return cmd1
}

func addProfileCmd(cmd *cobra.Command) {
	cmd.Flags().String(CMD_PROFILE, "", "")
	cmd.MarkFlagRequired(CMD_PROFILE)
	cmd.RegisterFlagCompletionFunc(CMD_PROFILE, getProfile)
}

func getProfile(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return listProfiles(), cobra.ShellCompDirectiveNoFileComp
}

func getApiArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// XXX: disable showing contents of current directory

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func ReverseStringSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func apiMain(cobraCmd *cobra.Command, args []string) {
	var apis []string

	flags := cobraCmd.Flags()
	apis = append(apis, cobraCmd.Use)

	var c1 = cobraCmd
	for c1.HasParent() {
		c1 = c1.Parent()
		apis = append(apis, c1.Use)
	}

	ReverseStringSlice(apis)
	// cut the first cmd out
	apis = apis[1:]
	fullApi := strings.Join(apis, "-")

	opts, ok := AwscConf.ApiOptions[fullApi]
	if !ok {
		fmt.Printf("Unsupported API: %s \n", fullApi)
		cobraCmd.Help()
		os.Exit(0)
	}

	var apiArgs []string
	apiArgs = append(apiArgs, opts.Args...)
	apiArgs = append(apiArgs, opts.Required...)

	_, isAdminVpc := AdminVpcCmds[fullApi]

	output, err := RunCmd([]string{fullApi}, apiArgs, isAdminVpc, flags)
	if err != nil {
		if output != "" {
			fmt.Printf("Output: %s \n", output)
		}

		fmt.Printf("ERR: %s \n", err)
		return
	}

	if output == "" {
		fmt.Printf("No Output\n")
		return
	}

	output1 := ParseOutput(output, opts.OutputField)
	if output1 == "" {
		output1 = output
	}

	output2 := FormatJson(output1)
	if output2 == "" || output2 == "{}" {
		output2 = output1
	}

	fmt.Printf("%s\n", output2)
}
