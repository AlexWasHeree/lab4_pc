#!/bin/bash

args=`find dataset -type f | xargs`

time bash go/concurrent/run.sh $args
time bash go/serial/run.sh $args
time bash go/serial-partial/run.sh $args