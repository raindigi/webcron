#!/bin/sh

case $1 in 
	start)
		nohup ./webcron 2>&1 >> info.log 2>&1 /dev/null &
		echo "Service started..."
		sleep 1
	;;
	stop)
		killall webcron
		echo "Service stopped..."
		sleep 1
	;;
	restart)
		killall webcron
		sleep 1
		nohup ./webcron 2>&1 >> info.log 2>&1 /dev/null &
		echo "Service restarted..."
		sleep 1
	;;
	*) 
		echo "$0 {start|stop|restart}"
		exit 4
	;;
esac

