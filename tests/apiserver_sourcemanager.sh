#!/bin/bash

cd $1/tests
go test base_test.go sourcemanager_test.go jobmanager_test.go udfmanager_test.go --run "TestEngineMap|TestPingSource|TestCreateSource|TestUpdateSource|TestDescribeSource|TestListSource|TestCreateWorkTable|TestUpdateWorkTable|TestDescribeWorkTable|TestListWorkTable|TestDeleteWorkTable|TestDisableSource|TestEnableSource|TestDeleteSource"
exit $?
