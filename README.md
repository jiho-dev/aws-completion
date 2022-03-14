
# Overview
 
aws-completion is a wrapper allowing shell completion by [TAB]  

# Install

copy `bin/$(OS)/awsc` to `$HOME/bin/`  
copy `config/awsc.yaml` to `$HOME/.aws/`  

If you want to add more aws commands under completion, add commands in `ApiPrefixFilter` section.  
And run below:   

`$ awsc generate-sub-cmds --profile <your_profile>`

Item starts the prefix of aws command  
## Default Prefix Filters 
```
ApiPrefixFilter:
    - describe-dhcp
    - describe-comp
    - describe-flow
    - describe-host
    - describe-i
    - describe-k
    - describe-n
    - describe-network
    - describe-p
    - describe-route
    - describe-secu
    - describe-sub
    - describe-tags
    - describe-v
    - get-console
```

# Completion  

## bash  
`$ source <(~/bin/awsc completion bash)` when shell starts  

## zsh
`$ ./awsc completion zsh > "${fpath[1]}/_awsc"`  

or 

`make zsh`  
