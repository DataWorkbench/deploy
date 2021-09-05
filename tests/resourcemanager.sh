#!/bin/bash

cd $1/tests
go test
exit $?
