#!/bin/bash

app='mimic-wac-go'
cmd=$1
pid=`ps -ef|grep ./main | grep -v grep|awk '{print $2}'`

startup(){
    cd /home/mimic-wac/wac-exb-go/
    nohup ./main 2>&1 > /home/mimic-wac/wac-go.log&
}


if [ ! $cmd ]; then
    echo "please specify command parameters 'start|restart|stop'"
    exit
fi

if [ $cmd == 'start' ]; then
    if [ ! $pid ]; then
        startup
    else
	echo "$app is already running! pid=$pid"
    fi
fi

if [ $cmd == 'restart' ]; then
    if [ $pid ]; then
        echo "$pid will be killed after 3 seconds!"
	sleep 1
	kill -9 $pid
    fi
    startup
fi

if [ $cmd == 'stop' ]; then
    if [ $pid ]; then
        echo "$pid whill be killed after 3 seconds!"
	sleep 1
	kill -9 $pid
    fi
    echo "$app is stoped"
fi

