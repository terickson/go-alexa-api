#!/bin/bash
rm -f go-alexa-api
env GOOS=linux GOARCH=arm GOARM=7 go build
scp go-alexa-api pi@media-server.home:
ssh pi@media-server.home "sudo systemctl stop go-alexa-api"
ssh pi@media-server.home "sudo mv /home/pi/go-alexa-api /usr/apps//go-alexa-api/alexa"
ssh pi@media-server.home "sudo systemctl start go-alexa-api"
