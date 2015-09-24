#!/bin/bash
for i in `seq 1 $1`;
do
	 ./map-coloring -file input/usa.in -heuristic=$2 > /dev/null
done 
