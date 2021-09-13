#! /bin/bash

FHASH=`shasum -a 256 $1` 
while true; do
    NHASH=`shasum -a 256 $1`
    if [ "$NHASH" != "$FHASH" ]; then
        ./mdp -file $1 FHASH=$NHASH
    fi
sleep 5

done