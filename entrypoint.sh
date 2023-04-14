#!/bin/bash

while true; do
	echo "attempting to connect to vmi $1 at namespace $2 streaming console log to stdout"
	/usr/bin/virtctl console $1 --namespace $2
	sleep 2
done

