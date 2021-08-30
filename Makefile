# Использовать bash с опцией pipefail
# pipefail - фейлит выполнение пайпа, если команда выполнилась с ошибкой
#SHELL=/bin/bash -o pipefail
# SHELL = /bin/bash

CURRENT_DIR = $$(pwd)

# Подготовка Makefile
# https://habr.com/ru/post/449910/#makefile_preparation

UNAME := $(shell uname)
BUILD_DATE := $(shell date +%Y%m%d-%H%M)

# envrioments from .env file
ifeq (,$(wildcard .env))
  $(shell test ! -f example.env || cp example.env .env || sleep 5; echo "create .env")
endif
ifeq (,$(wildcard .env))
	$(shell exit 1;)
else
	include .env
	export $(shell sed 's/=.*//' .env)
	# export
endif

os ?= $(shell uname|tr A-Z a-z)
ifeq ($(shell uname -m),x86_64)
  arch   ?= "amd64"
endif
ifeq ($(shell uname -m),i686)
  arch   ?= "386"
endif
ifeq ($(shell uname -m),aarch64)
  arch   ?= "arm"
endif
ifeq ($(shell uname -s),Darwin)
  arch   ?= "darwin"
endif

AUTO_APPROVE :=

PROJECT_NAME ?= "noname"

# Если переменная CI_JOB_ID не определена
ifeq ($(CI_JOB_ID),)
	# присваиваем значение local
	CI_JOB_ID := local
else
	AUTO_APPROVE := "-auto-approve"
endif

ifeq ($(TAG),)
  TAG := latest
endif

ifeq ($(CI_PROJECT_DIR),)
  CI_PROJECT_DIR := $(PWD)
endif

ifeq ($(MODE),)
  MODE := prod
endif

ifeq ($(APP_NAME),)
  APP_NAME := noname
endif

ifeq ($(AWS_ECR_NAME),)
  AWS_ECR_NAME := "docker.io"
endif

ifeq ($(AWS_REPO_NAME),)
  AWS_REPO_NAME := "user"
endif

ifeq ($(CI_JWERF_IMAGES_REPOOB_ID),)
  WERF_IMAGES_REPO := "${AWS_ECR_NAME}/${AWS_REPO_NAME}/${APP_NAME}"
endif

ifeq ($(K8S_NAMESPACE),)
  K8S_NAMESPACE := "default"
endif

ifeq ($(DEPLOY_MODE),)
  DEPLOY_MODE := "none"
endif

ifeq ($(CLUSTER_NAME),)
  CLUSTER_NAME := "my-claster"
endif

ifeq ($(AWS_PROFILE),)
  AWS_PROFILE := "default"
endif

ifeq ($(REGION),)
  REGION := eu-west-2
endif

# export BUILD_DATE
# export K8S_NAMESPACE
# export CI_JOB_ID
# export TAG
# export CI_PROJECT_DIR
# export MODE
# export APP_NAME
# export WERF_IMAGES_REPO
# export DEPLOY_MODE
# export CLUSTER_NAME
# # export AWS_PROFILE
# export AWS_REGION
export COMPOSE_HTTP_TIMEOUT=120

# Read all subsquent tasks as arguments of the first task
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(args) $(RUN_ARGS):;@:)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
landscape   := $(shell command -v landscape 2> /dev/null)
# terraform   := $(shell command -v terraform 2> /dev/null)
debug       :=
# Defaulting to level: TRACE. Valid levels are: [TRACE DEBUG INFO WARN ERROR]
export TF_LOG=
# export TF_LOG=ERROR

.ONESHELL:
.SHELL := /bin/bash

BOLD=$(shell tput bold)
RED=$(shell tput setaf 1)
GREEN=$(shell tput setaf 2)
YELLOW=$(shell tput setaf 3)
RESET=$(shell tput sgr0)

.PHONY: help
.DEFAULT_GOAL := help
help:
	@echo "\n$(GREEN)Available commands$(RESET)"
	@echo "---------------------------------------------------------------------"
	@grep -h -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo "---------------------------------------------------------------------"
	@#echo "$(YELLOW)Example argument to install Packer:$(RESET) make install packer=true"
	@echo "$(YELLOW)Example for auto-approve:$(RESET) make globalvars-apply auto=yes"
	@echo "$(YELLOW)Debug mode show command not run:$(RESET) make vpc-apply debug=yes"
	@echo "$(YELLOW)Output json format:$(RESET) make output c=vpc-aws f=json \n"
#help:
#	@grep -E '[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

info:   ## for developer information
	@echo "Show..."
	@echo "$(GREEN)Project Name:$(RESET) ${PROJECT_NAME}"
	@#echo "---------------------------------------------------------------------"
	@echo "$(YELLOW)os:$(RESET) ${os} $(YELLOW)arch:$(RESET) ${arch}"
	@echo "$(YELLOW)UNAME:$(RESET) ${UNAME} $(YELLOW)BUILD_DATE:$(RESET) ${BUILD_DATE}"
	@echo "$(YELLOW)CURRENT_DIR:$(RESET) ${CURRENT_DIR}"
	@test -z "${ADMIN_USERNAME}" || echo "$(GREEN)Admin Username:$(RESET) ${ADMIN_USERNAME}"
	@test -z "${ADMIN_PASSWORD}" || echo "$(GREEN)Admin Password:$(RESET) ${ADMIN_PASSWORD}"
	@#echo "$(GREEN)Run command open in browser:$(RESET) open https://${NGINX_HOST}/admin/"
	@#echo "$(GREEN)$(RESET)"
	@test -z "${AWS_ACCESS_KEY_ID}" || echo "export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}"
	@test -z "${AWS_PROFILE}" || echo "AWS CLI profile: ${AWS_PROFILE}"
	@test -z "${REGION}" || test -z "${AWS_ACCESS_KEY}" || echo "AWS region: ${REGION}"
	@echo "---------------------------------------------------------------------"

checker:
ifeq ($(shell test -e $(VERSION_FILE) && echo -n yes),yes)
	@$(eval VERSION=$(shell cat $(VERSION_FILE)))
else
	@echo File $(VERSION_FILE) does not exist
	@exit 0;
endif

init: ## File initialization and structure
	@echo "Установку зависемостей"
	@#curl -L https://raw.githubusercontent.com/flant/multiwerf/master/get.sh | bash
	@#./scripts/init-mysql.sh

major:  ## Set major version
	@#git tag -a v1.0.1 -m 'version v1.0.10' && git push --tags
	@git tag $$(svu major)
	@git push --tags
	@#goreleaser --rm-dist

minor:  ##  Set minor version
	@git tag $$(svu minor)
	@git push --tags
	@#goreleaser --rm-dist

patch:  ##  Set patch version
	@git tag $$(svu patch)
	@git push --tags
	@#goreleaser --rm-dist

# update:  ## Update project
# 	@echo "Start Update project..."
# 	@git pull
# 	@cp ./env ./.env
# 	@sudo ./pre.sh
# 	@sudo systemctl start docker-compose@elasticsearch.service

build:  ## Build Werf
	@echo "Start Build project..."
	@goreleaser build --rm-dist -p 1 --single-target --snapshot
	@cp $(CURRENT_DIR)/dist/${PROJECT_NAME}_${os}_${arch}/${APP_NAME} $(CURRENT_DIR)/${APP_NAME}
	@rm -Rf $(CURRENT_DIR)/dist
	@echo "$(GREEN)Example RUN:$(RESET)"
	@echo "./${APP_NAME} infra -c sample --stack my-dev --verbose"
# 	@werf build --stages-storage :local --introspect-before-error

# release:  ## Build and publish Release
# 	@echo "Start Build and publish Release ${WERF_IMAGES_REPO}:${TAG}"
# 	@werf build-and-publish --stages-storage :local --images-repo=${WERF_IMAGES_REPO} --tag-custom=${TAG}

# restart:   ## Restart Docker container
# 	@echo "Restart..."
# 	@#sudo systemctl restart docker-compose@gw.service
# 	@docker-compose down
# 	@docker-compose up --build -d

# stop:   ## останавливаем локального контейнера
# 	@echo "Stop docker-compose..."
# 	@docker-compose down

# start:   ## запуск локального контейнера в режиме daemon
# 	@echo "Start docker-compose..."
# 	@sudo chmod -R 0777 ./data/mongodb
# 	@docker-compose up --build -d

# pull:
# 	@echo "Download image: ${WERF_IMAGES_REPO}/${APP_NAME}:${TAG}"
# 	@mkdir -p /opt/data
# 	@docker-compose pull -q

# up: pull  ## запуск локального контейнера
# 	@docker-compose up --build

# down:  ## остановка контейнера с удалением volume
# 	@docker-compose down -v

# connect:  ## одключаемся к контейнеру
# 	@docker-compose exec app /bin/bash

# clean:  ## очищаем после сборки проекта
# 	@echo "Start clean..."
# 	@werf stages cleanup --stages-storage :local --images-repo=${WERF_IMAGES_REPO}

# ps:   ## List containers
# 	@docker-compose ps

# logs:   ## Show ALL logs
# 	@docker-compose logs -f

# logs-app:   ## Show App logs
# 	@docker-compose logs -f app

# logs-mysql:   ## Show mysql logs
# 	@docker-compose logs -f mysql

# shell-mysql:   ## Connection to MySQL container
# 	@docker-compose exec mysql bash

# mysql-cli:   ## Connection to MySQL cli
# 	@docker-compose exec -e MYSQL_PWD=${DATASOURCE_PASSWORD} mysql mysql -u${DATASOURCE_USERNAME} -D ${MYSQL_DATABASE}

# logs-mongodb:   ## Show MongoDB logs
# 	@docker-compose logs -f mongodb

# shell-mongodb:   ## Connection to MongoDB container
# 	@docker-compose exec mongodb bash

