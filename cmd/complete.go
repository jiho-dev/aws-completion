package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	//"github.com/ericaro/compgen"
	"github.com/jiho-dev/aws-completion/compgen"
	"github.com/jiho-dev/aws-completion/config"
)

type CompInfo struct {
	Args   []string
	Inword bool
	Prefix string
	Pos    int
}

func (ci *CompInfo) LogPrint() {
	log.Printf("CompInfo: args: %v, Pos:%d, Inword: %v, Prefix: %s  \n",
		ci.Args, ci.Pos, ci.Inword, ci.Prefix)
}

func (ci *CompInfo) ShowParam() {
	var names = []string{"COMP_LINE", "COMP_POINT", "COMP_TYPE", "COMP_KEY", "COMP_WORDS", "COMP_CWORD"}

	for _, name := range names {
		v, ok := os.LookupEnv(name)
		if ok {
			log.Printf("env %s: %s", name, v)
		} else {
			log.Printf("no env %s", name)
		}
	}
}

func (ci *CompInfo) GetLastArg() (string, int) {
	var cmd string
	var idx int

	for i, c := range ci.Args {
		if strings.HasPrefix(c, "-") {
			break
		}

		idx = i
		cmd = c
	}

	return cmd, idx
}

func (ci *CompInfo) ArgsHasOptions() bool {
	for _, c := range ci.Args {
		if strings.HasPrefix(c, "-") {
			return true
		}
	}

	return false
}

func (ci *CompInfo) ParseArgs() error {
	args, inword, err := compgen.Args()
	if err != nil {
		log.Printf("ERR: %v ", err)

		return err
	}

	ci.Args = args
	ci.Inword = inword
	ci.Pos, ci.Prefix = compgen.Prefix(args, inword)

	return nil
}

func (ci *CompInfo) GetFullCmd() string {
	var fullCmd string

	for i := 1; i < len(ci.Args); i++ {
		// is the first option ?
		if strings.HasPrefix(ci.Args[i], "-") {
			return fullCmd
		}

		if len(fullCmd) > 0 {
			fullCmd += "-"
		}

		fullCmd += ci.Args[i]
	}

	return fullCmd
}

func (ci *CompInfo) GetSubCommand() ([]string, *config.ApiOption) {
	conf := AwscConf

	var cmds []string
	keys := map[string]bool{}

	fullCmd := ci.GetFullCmd()
	c := strings.Count(fullCmd, "-")

	log.Printf("fullCmd: %s, dashcnt: %d", fullCmd, c)

	var matchedOpt *config.ApiOption

	// find the exact command
	opt, ok := conf.ApiOptions[fullCmd]
	if ok {
		//log.Printf("Exact Matched cmd: %+v", opt)
		matchedOpt = &opt
	}

	// add dash at the end to search more command
	if len(fullCmd) > 0 && len(ci.Prefix) < 1 {
		fullCmd += "-"
	}

	c = strings.Count(fullCmd, "-")

	// find more command
	for key, _ := range conf.ApiOptions {
		log.Printf("apiopt key:%s ", key)

		if fullCmd == key {
			// exact matched
			continue
		} else if len(fullCmd) < 1 {
			s := strings.Split(key, "-")
			keys[s[0]] = true
			continue
		}

		if !strings.HasPrefix(key, fullCmd) {
			continue
		}

		key1 := strings.Split(key, "-")
		keys[key1[c]] = true
	}

	for key, _ := range keys {
		cmds = append(cmds, key)
	}

	return cmds, matchedOpt
}

// complete performs bash command line completion for defined flags
func complete() bool {

	var cmds []string
	var ci CompInfo
	if !compgen.IsCompletionMode() {
		return false
	}

	//ci.ShowParam()

	err := ci.ParseArgs()
	if err != nil {
		log.Printf("ERR: %v ", err)
		return true
	}

	ci.LogPrint()

	subcmds, opt := ci.GetSubCommand()
	lastCmd, lastIdx := ci.GetLastArg()

	if !ci.ArgsHasOptions() {
		cmds = append(cmds, subcmds...)
	}

	/*
		log.Printf("subcmd list: %v", subcmds)
		log.Printf("Matched Opt: %+v", opt)
		log.Printf("lastCmd: %d, %s", lastIdx, lastCmd)
	*/

	if opt != nil {
		fs := flag.NewFlagSet(lastCmd, flag.ContinueOnError)
		for _, c := range opt.Required {
			fs.String(c, c, "")
		}

		for _, c := range opt.Args {
			fs.String(c, c, "")
		}

		fs.String(CMD_PROFILE, "", "AWS Profile")
		fs.SetOutput(ioutil.Discard)

		t := compgen.NewTerminator(fs)
		// register the callback
		t.Flag(CMD_PROFILE, CompgenProfile)
		//t.Arg(1, CompgenSubCmd)

		args := ci.Args[lastIdx:]
		inword := ci.Inword
		lastArg := args[len(args)-1]

		//log.Printf("LastArg: %s", lastArg)
		//log.Printf("Terminator Org Args: %v", args)

		// add dash to automatically show options
		if !ci.Inword && len(ci.Prefix) == 0 && !strings.HasPrefix(lastArg, "-") {
			args = append(args, "--")
			inword = true
		}

		pred, err := t.Compgen(args, inword)
		if err != nil {
			log.Printf("err:%v", err)
			return true
		}

		cmds = append(cmds, pred...)
	}

	//log.Printf("final cmd: %v", cmds)

	// completion output
	for _, c := range cmds {
		fmt.Printf("%s\n", c)
	}

	return true
}

func CompgenSubCmd(prefix string) []string {

	return nil
}
