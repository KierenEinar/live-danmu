#!/bin/bash

if [ -z $GOPATH ]; then
        export GOPATH=/www/go
        echo "system not set gopath, set gopath to $GOPATH"
fi;

cd $GOPATH/src/live-danmu

git reset --hard && git pull

echo "=======================================start build============================================"
go clean && glide install &&  go build
echo "=======================================finished build========================================="

WORK_DIR=/www/workspace/live-danmu

if [ ! -d $WORK_DIR ]; then
   mkdir -p $WORK_DIR
fi;

cp $GOPATH/src/live-danmu/live-danmu $WORK_DIR/app
cp -r $GOPATH/src/live-danmu/conf $WORK_DIR/conf

SSH_USER="danmu"
SSH_PORT="22"

SSH_IP=("danmu-2" "danmu-3")


echo "======================================upload to cluster======================================="
for IP in ${SSH_IP[@]};
do
        ssh ${SSH_USER}@${IP} "if [ ! -d $WORK_DIR ]; then mkdir -p $WORK_DIR; fi;"
        scp -r ${WORK_DIR} ${SSH_USER}@${IP}:${WORK_DIR} && scp -r /www/bin/ ${SSH_USER}@${IP}:/www/
done