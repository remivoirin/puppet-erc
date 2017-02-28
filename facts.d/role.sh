#!/bin/bash

trap 'exit 1' ERR

SERVER="localhost"

myhostname=$(hostname -f)
myrole=$(curl -Ss http://${SERVER}:14002/role/fulltext/${myhostname})

[[ -z "${myrole}" ]] && exit 1

echo "role=${myrole}" && exit 0
