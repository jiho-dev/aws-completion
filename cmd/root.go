package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/jiho-dev/aws-completion/config"
	"github.com/spf13/cobra"
)

var awsDir = os.Getenv("HOME") + "/.aws/"
var awscConfigName = "awsc.yaml"
var AwscConf *config.AwscConfig

var CompOpt = cobra.CompletionOptions{
	DisableDefaultCmd:   true,
	DisableNoDescFlag:   true,
	HiddenDefaultCmd:    true,
	DisableDescriptions: true,
}

/////////////////////////////////

var rootCmd = &cobra.Command{
	Use:               "awsc",
	Short:             "awsc <api> <sub-api> [flags]",
	Long:              "aws-completion to support shell completion of AWS APIs",
	CompletionOptions: CompOpt,
}

var CompletionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  "To load completions",
	Hidden:                true,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletionV2(os.Stdout, false)
		case "zsh":
			cmd.Root().GenZshCompletionNoDesc(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, false)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

var showEc2ApiCmd = &cobra.Command{
	Use:               "show-ec2-api",
	Short:             "show ec2 api",
	Run:               ShowApiMain,
	ValidArgsFunction: getApiArgs,
}

var showAdminVpcApiCmd = &cobra.Command{
	Use:               "show-admin-vpc-api",
	Short:             "show admin-vpc api",
	Run:               ShowApiMain,
	ValidArgsFunction: getApiArgs,
}

func init() {
	confFile := path.Join(awsDir, awscConfigName)
	conf, err := config.ParseConfig(confFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ParseConfig Err: %v\n", err)
		os.Exit(1)
	}

	AwscConf = conf

	for api, subApi := range conf.ApiPrefix {
		subCmd := InitApiCmd(api, subApi)
		rootCmd.AddCommand(subCmd)
	}

	rootCmd.AddCommand(CompletionCmd)

	genCmd := InitGenerateCmd()
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(showEc2ApiCmd)
	rootCmd.AddCommand(showAdminVpcApiCmd)
}

// Execute executes cmd
func Execute() error {
	return rootCmd.Execute()
}

func Help(cmd *cobra.Command, s []string) {
	fmt.Printf("%s: aws completion \n\n", cmd.Use)
	fmt.Printf("Usage: %s <api> <sub-api> [flags] \n", cmd.Use)
}
