package cmd

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/jiho-dev/aws-completion/config"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var AdminVpcCmds = map[string]bool{
	"allocate-random-ip-pool":        true,
	"blackpearl-health":              true,
	"create-public-ipv4-pool":        true,
	"delete-public-ipv4-pool":        true,
	"deregister-public-ipv4-pool":    true,
	"disable-public-ipv4-pool":       true,
	"disassociate-public-ip":         true,
	"enable-public-ipv4-pool":        true,
	"list-address-associations":      true,
	"list-blackpearl":                true,
	"list-network-acl":               true,
	"list-network-interface":         true,
	"list-public-ips":                true,
	"list-public-ipv4-pool":          true,
	"list-route-table":               true,
	"list-security-group":            true,
	"list-vrouters":                  true,
	"register-public-ipv4-pool":      true,
	"release-ip-pool":                true,
	"release-public-ip":              true,
	"request-ip-pool":                true,
	"show-dataversion":               true,
	"show-flowlog":                   true,
	"show-network-interface":         true,
	"show-papyrus-flowlog":           true,
	"show-papyrus-summary":           true,
	"show-revision":                  true,
	"show-snat":                      true,
	"show-summary":                   true,
	"show-vrevision":                 true,
	"show-vrouter-flowlog":           true,
	"show-vrouter-flow":              true,
	"show-vrouter-network-acl":       true,
	"show-vrouter-network-interface": true,
	"show-vrouter-port":              true,
	"show-vrouter-route":             true,
	"show-vrouter-security-group":    true,
	"show-vrouter-subnet":            true,
	"show-vrouter-summary":           true,
	"show-vrouter-table":             true,
	"update-network-interface":       true,
}

func InitGenerateCmd() *cobra.Command {
	var genCmd = &cobra.Command{
		Use:   CMD_GENERATE_SUB_CMDS,
		Short: "Generate sub-commands",
		//Hidden:                true,
		//DisableFlagsInUseLine: true,
		//ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		//Args: cobra.ExactValidArgs(1),
		Run:               generateApiMain,
		ValidArgsFunction: getApiArgs,
	}

	addProfileCmd(genCmd)

	return genCmd
}

func Contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func ContainPrefixs(prefix []string, searchterm string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(searchterm, p) {
			return true
		}
	}

	return false
}

func GetEc2Apis(prefixFilters []string) []string {
	var ec2Cmds []string

	//out := aws ec2 help | grep -E "o "
	out, err := ExecuteAwsCli("aws", "ec2", "help")
	if err != nil {
		fmt.Printf("Err: %s \n", err)
		return nil
	}

	//fmt.Printf("out: [%s] \n", out)
	cmds := strings.Split(out, "\n")
	for _, tmp := range cmds {
		tmp = strings.TrimSpace(tmp)
		if tmp == "" {
			continue
		}

		tmps := strings.Split(tmp, " ")
		if len(tmps) < 2 {
			continue
		}

		cmd := tmps[1]

		if len(prefixFilters) > 0 &&
			!ContainPrefixs(prefixFilters, cmd) {
			continue
		} else if len(cmd) < 1 {
			continue
		} else if cmd[0] < 'a' || 'z' < cmd[0] {
			continue
		}

		if cmd == "describe-local-gateway-route-table-virtual-interface-group-associa-" {
			cmd = "describe-local-gateway-route-table-virtual-interface-group-associations"
		}

		ec2Cmds = append(ec2Cmds, cmd)
	}

	sort.Strings(ec2Cmds)
	return ec2Cmds
}

////////////////////////////

func generateApiMain(cobraCmd *cobra.Command, args []string) {
	prefixFilter := AwscConf.ApiPrefixFilter
	allApis := GetEc2Apis(prefixFilter)

	// append AdminVpc APIs
	for k, _ := range AdminVpcCmds {
		allApis = append(allApis, k)
	}

	sort.Strings(allApis)

	AwscConf.ApiPrefix = generateApiPrefix(allApis)

	apiOptions := map[string]config.ApiOption{}
	flags := cobraCmd.Flags()
	generateApiParameters(apiOptions, allApis, flags)

	AwscConf.ApiOptions = apiOptions

	confFile := path.Join(awsDir, awscConfigName)
	config.WriteConfig(AwscConf, confFile)
}

func generateApiPrefix(allApis []string) map[string]map[string][]string {
	apiPrefix := map[string]map[string][]string{}

	for _, cmd := range allApis {
		//fmt.Printf("Ec2Cmd: %s \n", cmd)

		c := strings.Split(cmd, "-")
		first := c[0]
		firstMap, ok := apiPrefix[first]
		if !ok {
			firstMap = map[string][]string{}
			apiPrefix[first] = firstMap
		}

		if len(c) < 2 {
			firstMap[config.API_TERMINATED] = nil
			continue
		}

		second := c[1]
		rest := strings.Join(c[2:], "-")

		secondList, ok1 := firstMap[second]
		if !ok1 {
			secondList = []string{}
		}

		if len(rest) < 1 {
			rest = config.API_TERMINATED
		}

		secondList = append(secondList, rest)

		firstMap[second] = secondList
	}

	return apiPrefix
}

func generateApiParameters(apiOptions map[string]config.ApiOption, allApis []string, flags *flag.FlagSet) {
	flags.Bool(CMD_SHOW_HELP, true, "")

	for _, api := range allApis {
		var opt *config.ApiOption

		if _, ok := AdminVpcCmds[api]; ok {
			opt = generateAdminVpcParameters(api, flags)
		} else {
			opt = generateEc2ApiParameters(api)
		}

		if opt != nil {
			apiOptions[api] = *opt
		}
	}
}

func generateEc2ApiParameters(api string) *config.ApiOption {
	out, err := ExecuteAwsCli("aws", "ec2", api, "help")
	if err != nil {
		fmt.Printf("Err: %s \n", err)
		return nil
	}

	newOpts := config.ApiOption{}
	newOpts.OutputField = "Output"

	oldOpts, ok := AwscConf.ApiOptions[api]
	if ok {
		if oldOpts.OutputField != "" {
			newOpts.OutputField = oldOpts.OutputField
		}

		newOpts.Required = oldOpts.Required
	}

	args := strings.Split(out, "\n")
	var seeOpts, seeSyn bool
	for _, arg := range args {
		if strings.Contains(arg, "SYNOPSIS") {
			seeSyn = true
			continue
		} else if strings.Contains(arg, "OPTIONS") {
			seeOpts = true
		}

		if seeSyn && seeOpts {
			break
		}

		arg = strings.TrimSpace(arg)

		if strings.Contains(arg, "--dry-run") {
			continue
		} else if !strings.Contains(arg, "[--") {
			continue
		} else if !strings.Contains(arg, "<value>") {
			// XXX
			continue
		}

		if seeSyn {
			tmp := strings.Split(arg, " ")
			key := tmp[0][3:]

			if !Contains(newOpts.Required, key) {
				newOpts.Args = append(newOpts.Args, key)
			}
		}
	}

	return &newOpts
}

func generateAdminVpcParameters(api string, flags *flag.FlagSet) *config.ApiOption {
	inCmds := []string{api}
	output, err := RunCmd(inCmds, nil, true, flags)
	if err != nil {
		if output != "" {
			fmt.Printf("Output: %s \n", output)
		}

		fmt.Printf("ERR: %s \n", err)
		return nil
	}

	if output == "" {
		fmt.Printf("No Output\n")
		return nil
	}

	output1 := ParseOutput(output, "Result")
	if output1 == "" {
		output1 = output
	}

	output2 := FormatJson(output1)
	if output2 == "" || output2 == "{}" {
		output2 = output1
	}

	//fmt.Printf("%s\n", output2)

	newOpts := config.ApiOption{}
	newOpts.OutputField = "Result"

	oldOpts, ok := AwscConf.ApiOptions[api]
	if ok {
		if oldOpts.OutputField != "" {
			newOpts.OutputField = oldOpts.OutputField
		}

		newOpts.Required = oldOpts.Required
	}

	args := strings.Split(output2, "\n")
	var seeParams bool
	for _, arg := range args {
		if strings.HasPrefix(arg, "Parameters:") {
			seeParams = true
			continue
		}

		if seeParams {
			required := strings.Contains(arg, "(required)")

			arg = strings.TrimSpace(arg)
			tmp := strings.Split(arg, " ")
			key := tmp[0]
			key = strings.TrimSpace(key)

			if required {
				if !Contains(newOpts.Required, key) {
					newOpts.Required = append(newOpts.Required, key)
				}
			} else if !Contains(newOpts.Required, key) {
				newOpts.Args = append(newOpts.Args, key)
			}
		}
	}

	return &newOpts
}
