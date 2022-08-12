
# Overview
 
aws-completion(awscomp) is a wrapper allowing shell completion by [TAB]  
awscomp shows candidate words seperated by space from the origianal API,   
like `describe-network-interfaces` will be seperated into `describe` `network` `interfaces` 

# Install

copy `bin/$(OS)/awscomp` to `$HOME/bin/`  
copy `config/awscomp.yaml` to `$HOME/.aws/`  

If you want to update aws commands for completion , run generate-ec2-cmds as following:   

`$ awscomp generate ec2 cmds --profile <your_admin_profile>`

# Completion  

$ awscomp [TAB][TAB]

## bash  
add the line in .bashrc  
`complete -C ~/bin/awscomp awscomp`

## zsh

add the line in .zshrc   

`eval "$(awscomp completion-zsh)"`  


