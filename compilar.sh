#!/bin/bash

gzip -d ./main.c.gz
gcc ./main.c -o ./portalrobot
gzip ./main.c

go build

