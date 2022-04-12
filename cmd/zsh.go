package cmd

import (
	"fmt"
	"os"

	"text/template"

	"github.com/jiho-dev/aws-completion/config"
)

var zshSource = fmt.Sprintf(`
#compdef {{.ExecName}}
autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit

_{{.ExecName}}() {
	local output="$({{.Executable}} complete {{.Debug}} -- ${COMP_WORDS[@]:0:$COMP_CWORD} "${COMP_WORDS[$COMP_CWORD]}")"
	COMPREPLY=()
	local TMPIFS="$IFS"
	IFS=''
  while read -r line; do
		if [[ -n "$line" ]]; then
      COMPREPLY+=("$line")
    fi
	done <<< "$output"
	IFS="$TMPIFS"
}

complete -o nospace -F _{{.ExecName}} {{.ExecName}}
`)

func generateZshCompletion(enableDebug bool) {
	tmpl := template.Must(template.New("bash_source").Parse(zshSource))
	me, err := os.Executable()
	debug := ""

	if enableDebug {
		debug = "--debug"
	}
	if err != nil {
		logger.Errorf("Could not determine executable location: %v", err)
	}

	err = tmpl.Execute(os.Stdout, struct {
		Executable string
		ExecName   string
		Debug      string
	}{
		Executable: me,
		ExecName:   config.EXEC_NAME,
		Debug:      debug,
	})

	if err != nil {
		logger.Errorf("Could not render source template for zsh: %v", err)
	}
}
