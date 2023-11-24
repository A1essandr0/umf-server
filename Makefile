-include .env

GO := go
DOCKER := docker

run:
	${GO} run ./cmd/main.go

mockrun:
	UMF_DBSTORE_TYPE=mock UMF_KVSTORE_TYPE=mock ${GO} run ./cmd/main.go