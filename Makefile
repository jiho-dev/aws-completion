all:
	go build -o awsc
	/bin/cp ./awsc ~/bin


zsh:
	./awsc completion zsh > _awsc
	cp _awsc ~/.oh-my-zsh/functions/
