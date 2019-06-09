#!/bin/bash

[ -e ./terraform.d/plugins/linux_amd64 ] || mkdir -p ./terraform.d/plugins/linux_amd64
make build
