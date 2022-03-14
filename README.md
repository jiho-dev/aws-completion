
# Overview
 
aws-completion is a wrapper allowing shell completion by [TAB]  

# install

copy `bin/$(OS)/awsc` to `$HOME/bin/`  
copy `config/awsc.yaml` to `$HOME/.aws/`  

if you want to add more aws cli under completion, add commands `ApiPrefixFilter` section  
the list item starts aws command and has parameter name  

# Completion  

## bash  
`$ source <(~/bin/awsc completion bash)` when shell starts  

## zsh
`$ ./awsc completion zsh > "${fpath[1]}/_awsc"`  

or 

`make zsh`  
