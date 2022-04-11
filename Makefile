all:
#	go build -o awsc
#	/bin/cp ./awsc ~/bin

	go build -o awscomp
	/bin/mv ./awscomp ~/bin


zsh:
	./awsc completion zsh > _awsc
	cp _awsc ~/.oh-my-zsh/functions/
