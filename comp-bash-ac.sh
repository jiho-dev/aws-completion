#!/bin/bash

PATH=$PATH:$(pwd)
#source <(awsc completion bash)

complete -C `pwd`/ac ./ac

