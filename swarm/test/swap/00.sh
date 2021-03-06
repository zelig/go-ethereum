#!/bin/bash
echo "TEST swap/00:"
echo " two nodes that do not sync and do not have any funds"
echo " cannot retrieve content from each other"

dir=`dirname $0`
source $dir/../../cmd/swarm/test.sh

FILE_00=/tmp/1K.0
randomfile 1 > $FILE_00
ls -l $FILE_00
mininginterval=50
key=/tmp/key

swarm init 2
sleep $wait
swarm attach 00 -exec "'bzz.noSync(true)'"
swarm attach 01 -exec "'bzz.noSync(true)'"
swarm up 00 $FILE_00|tail -n1 > $key
swarm needs 00 $key $FILE_00
swarm needs 01 $key $FILE_00 | tail -1| grep -ql "PASS" && echo "FAIL" || echo "PASS <3"

FILE_01=/tmp/1K.1
randomfile 1 > $FILE_01
swarm up 01 $FILE_01|tail -1 > $key
swarm needs 01 $key $FILE_01
swarm needs 00 $key $FILE_01 | tail -1| grep -ql "PASS" && echo "FAIL" || echo "PASS <3"

swarm stop all