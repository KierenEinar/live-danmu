#!/bin/bash

SSH_USER=`whoami`
SSH_IP=("danmu-1" "danmu-2" "danmu-3")
BIN_DIR=/www/bin
function stopAll(){

        for IP in ${SSH_IP[@]}
        do
                ssh -f -n -P22  ${SSH_USER}@${IP} "bash -c ${BIN_DIR}/stop.sh"
        done
        echo "stop cluster all process success"

}

function startAll(){
        for IP in ${SSH_IP[@]}
        do
                echo "server ${IP} start process ..."
                echo "ssh ${SSH_USER}@${IP} bash -c ${BIN_DIR}/start.sh"
                ssh -f -n -P22 ${SSH_USER}@${IP} "${BIN_DIR}/start.sh "
        done
        echo "start cluster all process success"
}

function restartAll(){
        stopAll && startAll
}

while getopts "s:" opt; do
        case $opt in
                s)
                        [ "start" == $OPTARG ] && startAll
                        [ "stop" == $OPTARG ] && stopAll
                        [ "restart" == $OPTARG ] && restartAll
                ;;

                ?)
                        echo "invalid options" && exit 0
                ;;
        esac
done

