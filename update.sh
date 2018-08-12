#!/bin/sh
git pull
go run update.go
git add .
git commit -m "update framework.md" -a
git push
