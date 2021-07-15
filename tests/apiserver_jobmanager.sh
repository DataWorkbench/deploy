#!/bin/bash

cd $1/tests
go test base_test.go jobmanager_test.go sourcemanager_test.go --run "TestExplain|TestSyntaxCheck|TestPreview"
exit $?
