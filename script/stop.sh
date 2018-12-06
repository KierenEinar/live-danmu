#!/bin/bash
puser=`whoami`
pname="app"
pids=`ps -fu ${puser}| grep app | grep -v grep | awk '{print $2}'`
for pid in $pids
do
        kill -9 $pid
done
echo "stop danmu app server..............."

