#!/usr/bin/env bash

# run this scripts in docker to compile service
set -e

# setup variable env
BIN_DIR="${GOPATH}/bin"
CONF_DIR=""
current_path=$(cd "$(dirname "${0}")" || exit 1; pwd)
COMMON_MODULE="github.com/DataWorkbench/common/utils/buildinfo"
if [[ "${BUILD_MODE}" == "release" ]]; then
    TAGS="netgo jsoniter ${BUILD_MODE}"
else
    TAGS="netgo jsoniter"
fi


usage(){
  echo "compile.sh -s services [-o output_dir] [-p program_name] [-c config_dir]"
  echo "OPTIONS:"
  echo "    -s required, service wanted to compile, split by comma, default all services"
  echo "    -o output dir of compiled service, default ${GOPATH}/bin"
  echo "    -p the program name of compiled service, default same as services"
  echo "    -c the dir of the service config.yaml, if specified copy config.yaml to DIR/SERVICE.yaml , default ''"
}

# handle params
while getopts "hs:o:p:c:" opt;
do
  case $opt in
    s)
      echo "SERVICE_PARAM: ${OPTARG}"
      SERVICE_PARAM="${OPTARG}"
      ;;
    o)
      echo "BIN_DIR: ${OPTARG}"
      BIN_DIR="${OPTARG}"
      ;;
    c)
      echo "CONF_DIR: ${OPTARG}"
      CONF_DIR="${OPTARG}"
      ;;
    p)
      echo "p: ${OPTARG}"
      PROGRAM_PARAM="${OPTARG}"
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


_compile(){
  _service=$1
  cd "${current_path}/../../../${_service}" || exit 1
  if [ -n "${PROGRAM_PARAM}" ];then
    PROGRAM="${BIN_DIR}/${PROGRAM_PARAM}"
  else
    PROGRAM="${BIN_DIR}/${_service}"
  fi

  go build --tags "${TAGS}" -ldflags " \
  -X ${COMMON_MODULE}.GoVersion=$(go version|awk '{print $3}') \
  -X ${COMMON_MODULE}.CompileBy=$(git config user.email) \
  -X ${COMMON_MODULE}.CompileTime=$(date '+%Y-%m-%d:%H:%M:%S') \
  -X ${COMMON_MODULE}.GitBranch=$(git rev-parse --abbrev-ref HEAD) \
  -X ${COMMON_MODULE}.GitCommit=$(git rev-parse --short HEAD) \
  -X ${COMMON_MODULE}.OsArch=$(uname)/$(uname -m)" \
  -v -o "${PROGRAM}" .

  if [ -n "${CONF_DIR}" ]; then
    cp ./config/config.yaml "${CONF_DIR}/${_service}.yaml"
  fi

}


# handle param
mkdir -p "${BIN_DIR}"
if [ -n "${CONF_DIR}" ]; then
  mkdir -p "${CONF_DIR}"
fi

if [ -z "${SERVICE_PARAM}" ]; then
  echo "param service is required!"
  exit 1
else
  SERVICES=(${SERVICE_PARAM//,/ })
fi

# compile
# shellcheck disable=SC2068
for service in ${SERVICES[@]};
do
  echo "compile service: ${service}"
  _compile "${service}"
done
