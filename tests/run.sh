#!/bin/bash

PWD=`pwd`

run_command() {
	$@
	if [ $? != 0 ]
	then
		echo "run \"$@\" failed " >&2
		return 1
	fi
	
	return 0
}

run_test_stop_on_error() {
	run_command $@
	if [ $? != 0 ]
	then
		echo "stop this test. please fix the error." >&2
		exit 1
	fi
}

run_test_not_stop_on_error() {
	run_command $@
	if [ $? != 0 ]
	then
		echo "error happend but go on running this test." >&2
	fi
}

run_test_stop_on_error bash ./tests/sourcemanager.sh ${PWD}/../sourcemanager
run_test_stop_on_error bash ./tests/jobmanager.sh ${PWD}/../jobmanager
run_test_stop_on_error bash ./tests/apiserver_jobmanager.sh ${PWD}/../apiserver
