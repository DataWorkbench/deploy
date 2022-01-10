# Copyright 2020 The Dataomnis Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

pull:
	docker pull dockerhub.qingcloud.com/dataomnis/builder
	@echo "ok"

push:
	docker push dockerhub.qingcloud.com/dataomnis/builder
	docker push dockerhub.qingcloud.com/dataomnis/builder:zeppelin
	@echo "ok"

run:
	docker run --rm -it -v `pwd`:/root dockerhub.qingcloud.com/dataomnis/builder

clean:
	@echo "ok"

build:
	docker build -t dockerhub.qingcloud.com/dataomnis/builder .
	docker build -t dockerhub.qingcloud.com/dataomnis/builder:zeppelin -f Dockerfile.zeppelin .
	@echo "ok"