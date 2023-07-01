-include .env

GO := go
DOCKER := docker

run:
	${GO} run ./cmd/main.go
