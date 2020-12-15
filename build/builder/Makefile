# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

default:
	docker build -t datawh/builder .
	@echo "ok"

pull:
	docker pull datawh/builder
	@echo "ok"

run:
	docker run --rm -it -v `pwd`:/root datawh/builder

clean:
	@echo "ok"
