# Copyright 2020 The Databench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Repo:=dockerhub.databench.io/databench
TARG.Name:=databench
TRAG.Gopkg:=DataWorkbench
#TRAG.Version:=$(TRAG.Gopkg)/pkg/version

TAG:=dev
FLYWAY_IMAGE:=$(TARG.Repo)/flyway:$(TAG)
ZEPPELIN_IMAGE:=$(TARG.Repo)/zeppelin:0.9.0
FLINK_IMAGE:=$(TARG.Repo)/flinkutile:1.12.3-scala_2.11
BUILDER_IMAGE:=$(TARG.Repo)/builder:latest
BUILDER_IMAGE_ZEPPELIN:=$(TARG.Repo)/builder:zeppelin

LOCAL_CACHE:=`go env GOCACHE`
LOCAL_MODCACHE:=`go env GOPATH`/pkg
WORKDIR_IN_DOCKER=/$(TRAG.Gopkg)
RUN_IN_DOCKER:=docker run -it --rm -v `pwd`/..:$(WORKDIR_IN_DOCKER) -v ${LOCAL_CACHE}:/go/cache -v $(LOCAL_MODCACHE):/go/pkg -w $(WORKDIR_IN_DOCKER) $(BUILDER_IMAGE)

comma:= ,
empty:=
space:= $(empty) $(empty)
# the service that need to format/compile/build.., default all.
service=apiglobal,apiserver,spacemanager,flowmanager,jobmanager,jobdeveloper,jobwatcher,scheduler,sourcemanager,udfmanager,resourcemanager,notifier,account,enginemanager
SERVICE_ARRAY=$(subst ${comma},${space},$(service))
COMPOSE_DB_CTRL=databench-db-ctrl

.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: update-builder
update-builder: ## Pull databench-builder image
	docker pull $(BUILDER_IMAGE)
	docker pull $(BUILDER_IMAGE_ZEPPELIN)
	@echo "update-builder done"

.PHONY: compile
compile:
	$(RUN_IN_DOCKER) bash -c "time ./deploy/build/scripts/compile.sh -s $(service) -o $(WORKDIR_IN_DOCKER)/deploy/tmp/bin -c $(WORKDIR_IN_DOCKER)/deploy/tmp/conf";

.PHONY: build-flyway
build-flyway: ## Build flyway image for database migration
	cd .. && docker build -t $(FLYWAY_IMAGE) -f ./deploy/build/db/Dockerfile .

.PHONY: build-zeppelin
build-zeppelin: ## zeppelin, set perNote to isolate perUser to '', download lib from QingStor
	cd ./build/zeppelin && docker build -t $(ZEPPELIN_IMAGE) . && cd ../..

.PHONY: build-flink-utile
build-flink-utile:
	cd ./build/flink_utile/ && docker build -t $(FLINK_IMAGE) . && cd ../..

.PHONY: build-dev
build-dev: compile  ## Build databench image
	@cd .. && docker build -t $(TARG.Repo)/$(TARG.Name):$(TAG) -f ./deploy/Dockerfile.dev .
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: build-all  ## Build all images
build-all: build-flyway build-dev build-zeppelin

.PHONY: push-images  ## push all images
push-images:
	docker push $(TARG.Repo)/$(TARG.Name):$(TAG)
	docker push $(ZEPPELIN_IMAGE)
	docker push $(FLYWAY_IMAGE)

.PHONY: pull-images
pull-images: ## Pull images for docker-compose
	docker-compose pull --ignore-pull-failures
	@echo "pull-images done"

.PHONY: compose-migrate-db
compose-migrate-db: ## Migrate db in docker compose
	until docker-compose exec databench-db bash -c "echo 'SELECT VERSION();' | mysql -uroot -ppassword"; do echo "waiting for mysql"; sleep 2; done;
	docker-compose up $(COMPOSE_DB_CTRL)

.PHONY: compose-up
compose-up: ## Launch databench in docker compose
	docker-compose up -d databench-db
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

update: build-dev ## Update some service
	docker-compose stop $(SERVICE_ARRAY)
	docker-compose up --no-deps -d $(SERVICE_ARRAY)
	@echo "update service $(service) done"

update-service: build-dev ## Update databench service in k8s, eg: make update-k8s service=s1,s2
	docker push $(TARG.Repo)/$(TARG.Name):$(TAG)
	@echo "push image done"
	@$(foreach s,$(SERVICE_ARRAY),kubectl -n databench rollout restart deployment databench-$(s);)
	@echo "update service $(service) done"

.PHONY: clean
clean:
	rm -r ./tmp

.PHONY: test
test: ## Launch databench in docker compose
	@bash ./tests/run.

.DEFAULT_GOAL = help
# Target name % means that it is a rule that matches anything, @: is a recipe;
# the : means do nothing
%:
	@:
