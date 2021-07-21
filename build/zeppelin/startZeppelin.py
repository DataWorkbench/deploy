#!/opt/conda/bin/python3

import json
import sys
import os
import requests
import time

fileName = "/zeppelin/conf/interpreter.json"
TiniCmd = "/usr/bin/tini"
zeppelinCmd = "/zeppelin/bin/zeppelin.sh"
#libUrl = os.getenv("LIB_URL")
#libName = os.getenv("LIB_NAMES")
libDir = "/opt/zeppelin/lib/"

def log(msg):
    print(msg, file=sys.stderr)

def rewriteZeppelinConf():
    f = open(fileName,'r')
    load_dict = json.load(f)
    load_dict["interpreterSettings"]["flink"]["option"]["perNote"] = "isolated"
    load_dict["interpreterSettings"]["flink"]["option"]["perUser"] = ""
    f.close()

    f = open(fileName,'w')
    json.dump(load_dict, f, indent=2)
    f.close()
    log("set flink perNote to 'isolated', perUser to ''")

#def downloadLib():
#    for lib in libName.split(','):
#        log("download " + lib + " to " + libDir + " is running")
#        r = requests.get(libUrl + lib)
#        with open(libDir + lib, "wb") as f:
#            f.write(r.content)
#        log("download " + lib + " to " + libDir + " done")


def startRealZeppelin(arg):
    cmd = TiniCmd + " " + arg
    log("start zeppelin with command: " + cmd)
    if os.system(cmd) != 0:
        raise RuntimeError("start zeppelin failed")

def main(arg):
    try:
        rewriteZeppelinConf()
#        downloadLib()
        startRealZeppelin(arg)
    except Exception as err:
        log(err)
        sys.exit(1)

if __name__ == '__main__':
    if len(sys.argv) == 1:
        arg = zeppelinCmd
    else:
        arg = " ".join(sys.argv[1:])

    main(arg)
