#compdef awscomp
autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit

_bosh_comp() {
    local output="$(/Users/jiho.jung/bin/awscomp complete  -- ${COMP_WORDS[@]:0:$COMP_CWORD} "${COMP_WORDS[$COMP_CWORD]}")"
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

complete -o nospace -F _bosh_comp awscomp

