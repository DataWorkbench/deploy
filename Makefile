# Copyright 2020 The Dataomnis Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

repo:=dockerhub.dataomnis.io
TRAG.Gopkg:=DataWorkbench

DEFAULT_TAG=dev
tag:=$(DEFAULT_TAG)
FLYWAY_IMAGE:=$(repo)/dataomnis/flyway:$(tag)
ZEPPELIN_IMAGE:=$(repo)/dataomnis/zeppelin:0.9.0
FLINK_IMAGE:=$(repo)/dataomnis/flinkutile:1.12.3-scala_2.11
BUILDER_IMAGE:=dockerhub.qingcloud.com/dataomnis/builder:latest
BUILDER_IMAGE_ZEPPELIN:=dockerhub.qingcloud.com/dataomnis/builder:zeppelin

LOCAL_CACHE:=`go env GOCACHE`
LOCAL_MODCACHE:=`go env GOPATH`/pkg
WORKDIR_IN_DOCKER=/$(TRAG.Gopkg)
RUN_IN_DOCKER:=docker run -it --rm -v `pwd`/..:$(WORKDIR_IN_DOCKER) -v ${LOCAL_CACHE}:/go/cache -v $(LOCAL_MODCACHE):/go/pkg -w $(WORKDIR_IN_DOCKER) $(BUILDER_IMAGE)

PWD_DIR:=$(shell pwd)
DOCKERFILE:=./deploy/Dockerfile
ifeq ($(tag),$(DEFAULT_TAG))
	DOCKERFILE=./deploy/Dockerfile.dev
endif

comma:= ,
empty:=
space:= $(empty) $(empty)
# the service that need to format/compile/build.., default all.
service=apiglobal,apiserver,spacemanager,scheduler,resourcemanager,account,enginemanager,notifier
SERVICE_ARRAY=$(subst ${comma},${space},$(service))
COMPOSE_DB_CTRL=dataomnis-db-ctrl

.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: update-builder
update-builder: ## Pull dataomnis-builder image
	docker pull $(BUILDER_IMAGE)
	docker pull $(BUILDER_IMAGE_ZEPPELIN)
	@echo "update-builder done"

.PHONY: compile
compile:
	$(RUN_IN_DOCKER) bash -c "time ./deploy/build/scripts/compile.sh -s $(service) -o $(WORKDIR_IN_DOCKER)/deploy/tmp/bin -c $(WORKDIR_IN_DOCKER)/deploy/tmp/conf";

.PHONY: build-flyway
build-flyway: ## Build flyway image for database migration
	cd ./build/db && docker build -t $(FLYWAY_IMAGE) .

.PHONY: build-zeppelin
build-zeppelin: ## zeppelin, set perNote to isolate perUser to '', download lib from QingStor
	cd ./build/zeppelin && docker build -t $(ZEPPELIN_IMAGE) . && cd ../..

.PHONY: build-flink-utile
build-flink-utile:
	cd ./build/flink_utile/ && docker build -t $(FLINK_IMAGE) . && cd ../..

.PHONY: build-image
build-image: compile  ## Build dataomnis image
	$(foreach S,$(SERVICE_ARRAY),cd $(PWD_DIR)/.. && docker build --build-arg SERVICE=$(S) -t $(repo)/dataomnis/$(S):$(tag) -f $(DOCKERFILE) .;)
	docker image prune -f 1>/dev/null 2>&1
	@echo "build $(service) done"

.PHONY: build-all  ## Build all images
build-all: build-flyway build-zeppelin build

.PHONY: push-flyway
push-flyway:
	docker push $(FLYWAY_IMAGE)
	@echo "push flyway image done"

.PHONY: push-zeppelin
push-zeppelin:
	docker push $(ZEPPELIN_IMAGE)
	@echo "push zeppelin image done"

.PHONY: push-image  ## push dataomnis service image
push-image:
	@$(foreach S,$(SERVICE_ARRAY),docker push $(repo)/dataomnis/$(S):$(tag);)
	@echo "push $(service) done"

.PHONY: pull-images
pull-images: ## Pull images for docker-compose
	docker-compose pull --ignore-pull-failures
	@echo "pull-images done"

.PHONY: compose-migrate-db
compose-migrate-db: ## Migrate db in docker compose
	until docker-compose exec dataomnis-db bash -c "echo 'SELECT VERSION();' | mysql -uroot -ppassword"; do echo "waiting for mysql"; sleep 2; done;
	docker-compose up $(COMPOSE_DB_CTRL)

.PHONY: compose-up
compose-up: ## Launch dataomnis in docker compose
	docker-compose up -d dataomnis-db
	make compose-migrate-db
	docker-compose up --remove-orphans -d
	@echo "compose-up done"

.PHONY: compose-down
compose-down: ## Shutdown service in docker compose
	docker-compose down
	@echo "compose-down done"

.PHONY: compose-logs-f
compose-logs-f: ## Follow service log in docker compose
	docker-compose logs --tail 10 -f $(SERVICE_ARRAY)

update: build-image ## Update some service
	docker-compose stop $(SERVICE_ARRAY)
	docker-compose up --no-deps -d $(SERVICE_ARRAY)
	@echo "update service $(service) done"

update-service: build-image push-image ## Update dataomnis service in k8s, eg: make update-k8s service=s1,s2
	@$(foreach s,$(SERVICE_ARRAY),kubectl -n dataomnis rollout restart deployment dataomnis-$(s);)
	@echo "update service $(service) done"

.PHONY: clean
clean:
	rm -r ./tmp

.PHONY: test
test: ## Launch dataomnis in docker compose
	@bash ./tests/run.

.DEFAULT_GOAL = help
# Target name % means that it is a rule that matches anything, @: is a recipe;
# the : means do nothing
%:
	@:
