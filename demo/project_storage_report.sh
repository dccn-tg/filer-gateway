#!/bin/bash

###########################################################################
#
# This script receives a list of project and prints out (for each project)
# the following storage information in CSV format
#
#   - storage system
#   - storage quota in GiB
#   - storage usage in GiB
#
# Usage:
#
#   $ ./project_storage_report.sh <project1> [<project2>] [...]
#
###########################################################################
for p in $@; do
    d=$(echo $p | ./demo.sh project get | jq '.storage')
    q=$(echo $d | jq '.quotaGb')
    u=$(echo $d | jq '.usageMb')
    s=$(echo $d | jq '.system')
    printf "%s,%s,%d,%d\n" $p $s $q $(( $u / 1024 ))
done
