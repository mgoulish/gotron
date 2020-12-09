#! /bin/bash


go run ./gotron.go  A  127.0.0.1   9091    9090 &

go run ./gotron.go  B  127.0.0.1   9090    9091 &



ps -aef | grep gotron | grep -v grep




