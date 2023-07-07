#/bin/bash

SERVICE_NAME="lesscode-usercenter"

BASE_PATH=`cd "$(dirname "$0")"; pwd`
cd $BASE_PATH/


./usercenter &
