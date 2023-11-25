-include .env

GO := go
DOCKER := docker
APP_ENTRY_POINT := ./cmd/main.go
BINARY_OUT := umf
DOCKER_IMAGE_NAME := umf-service
SERVICE_VERSION := ${shell cat __version__}

run:
	${GO} run ${APP_ENTRY_POINT}

mockrun:
	UMF_DBSTORE_TYPE=mock UMF_KVSTORE_TYPE=mock ${GO} run ${APP_ENTRY_POINT}


test-verbose:
	${GO} test -v

test:
	${GO} test


deploy:
	echo 'Compiling...'
	go build -o ${BINARY_OUT} ${APP_ENTRY_POINT}
	echo 'Deploying...'
	rsync -uazp umf ${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}
	rsync -uazp cmd/config.yaml ${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}cmd/
	echo 'Restarting service...'
	ssh ${SERVER_USER}@${SERVER_HOST} 'sudo service umf restart'
	echo 'Removing build files...'
	rm umf

deploy-kite:
	echo 'Compiling...'
	go build -o ${BINARY_OUT} ${APP_ENTRY_POINT}
	echo 'Deploying to Kite server...'
	rsync -uazp umf ${PK_SERVER_USER}@${PK_SERVER_HOST}:${PK_SERVER_PATH}
	rsync -uazp cmd/config.yaml ${PK_SERVER_USER}@${PK_SERVER_HOST}:${PK_SERVER_PATH}cmd/
	echo 'Restarting service...'
	ssh ${PK_SERVER_USER}@${PK_SERVER_HOST} 'sudo service umf restart'
	echo 'Removing build files...'
	rm umf


build-image:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME} --network=host .


bump-build-version:
	cat __version__ | awk -F. -v OFS=. '{$$3++ ; print}' >  __new_version__
	mv __new_version__ __version__

bump-minor-version:
	cat __version__ | awk -F. -v OFS=. '{$$2++ ; $$3=0; print}' >  __new_version__
	mv __new_version__ __version__

bump-major-version:
	cat __version__ | awk -F. -v OFS=. '{$$1++ ; $$2=0; $$3=0; print}' >  __new_version__
	mv __new_version__ __version__