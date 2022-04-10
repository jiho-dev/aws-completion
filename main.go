package main

import (
	"github.com/jiho-dev/aws-completion/cmd"
	"github.com/jiho-dev/aws-completion/config"
)

func main() {
	if true {
		_ = cmd.Execute()
		//cmd.ShowEc2Cmd()
	} else {
		// for test
		config.YamlTest()
	}

}
