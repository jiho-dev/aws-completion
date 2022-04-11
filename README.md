
# Overview
 
aws-completion is a wrapper allowing shell completion by [TAB]  

# Install

copy `bin/$(OS)/awscomp` to `$HOME/bin/`  
copy `config/awscomp.yaml` to `$HOME/.aws/`  

If you want to update aws commands under completion, run generate-ec2-cmds
And run below:   

`$ awscomp generate ec2 cmds --profile <your_admin_profile>`

# Completion  

$ awscomp [TAB][TAB]

## bash  
add the line in .bashrc
`complete -C ~/bin/awscomp awscomp`

## zsh
`$ ./awsc completion zsh > "${fpath[1]}/_awsc"`  

or 

`make zsh`  
