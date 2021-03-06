#!/usr/bin/env bash

# run this scripts in docker to compile service

set -e

SERVICES=""
OUTPUT_DIR="${GOPATH}/bin"

usage(){
  echo "compile.sh -s services [-o output_dir] [-p program]"
  echo "OPTIONS:"
  echo "    -s required, service wanted to compile, split by comma, default all services"
  echo "    -o output dir of compiled service, default ${GOPATH}/bin"
  echo "    -g the program name of compiled service, default same as services"
}

while getopts ":hs:o:p:" opt;
do
  case $opt in
    s)
      SERVICES=(`echo $OPTARG | tr ',' ' '`)
      ;;
    o)
      OUTPUT_DIR=$OPTARG
      ;;
    p)
      PROGRAM=$OPTARG
      ;;
    h) #help
      usage
      exit 0
      ;;
    ?)
      usage
      exit 1
  esac
done

if [ ! -n "${SERVICES}" ]; then
  echo "SERVICES is required."
  usage
  exit 1
fi

_compile(){
  _service=$1
  cd "${current_path}"/../../../${_service} || exit 1
  if [ ! -n "${PROGRAM}" ];then
    OUTPUT="${OUTPUT_DIR}/${_service}"
  else
    OUTPUT="${OUTPUT_DIR}/${PROGRAM}"
  fi

  go build --tags "${TAGS}" -ldflags " \
  -X ${COMMON_MODULE}.GoVersion=$(go version|awk '{print $3}') \
  -X ${COMMON_MODULE}.CompileBy=$(git config user.email) \
  -X ${COMMON_MODULE}.CompileTime=$(date '+%Y-%m-%d:%H:%M:%S') \
  -X ${COMMON_MODULE}.GitBranch=$(git rev-parse --abbrev-ref HEAD) \
  -X ${COMMON_MODULE}.GitCommit=$(git rev-parse --short HEAD) \
  -X ${COMMON_MODULE}.OsArch=$(uname)/$(uname -m)" \
  -v -o ${OUTPUT} .
}

mkdir -p ${OUTPUT_DIR}
current_path=$(cd "$(dirname "${0}")" || exit 1; pwd)
COMMON_MODULE="github.com/DataWorkbench/common/utils/buildinfo"

if [[ "${BUILD_MODE}" == "release" ]]; then
    TAGS="netgo jsoniter ${BUILD_MODE}"
else
    TAGS="netgo jsoniter"
fi

for service in ${SERVICES[@]};
do
  echo "compile service: ${service}"
  _compile $service
done
