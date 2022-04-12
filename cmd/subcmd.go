package cmd

import (
	"fmt"
	"os"
)

func apiMainExt(fullApi string, args []string) {
	opts, ok := AwscConf.ApiOptions[fullApi]
	if !ok {
		fmt.Printf("Unsupported API: %s \n", fullApi)
		os.Exit(0)
	}

	flags := opts.GetFlags(fullApi)
	flags.String("profile", "", "AWS Profile")

	if len(args) > 1 {
		flags.Parse(args)
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
