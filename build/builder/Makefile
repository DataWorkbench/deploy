# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

pull:
	docker pull dataworkbench/builder
	@echo "ok"

run:
	docker run --rm -it -v `pwd`:/root dataworkbench/builder

clean:
	@echo "ok"

build:
	docker build -t dataworkbench/builder .
	@echo "ok"