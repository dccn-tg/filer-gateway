#!/bin/bash

function usage() {

cat <<EOF

Usage: $0 <user|project> <new|get|set>

It requires the three env. variables to be set:

- API_USER
- API_URL
- API_KEY

EOF

}

[ $# -lt 2 ] && usage && exit 1

script="./${1}.sh"

# run demo against filer-gateway, it requires the following env. variables
#  - API_USER
#  - API_URL
#  - API_KEY
#
# you could either set those env. variables before running this script, or modifiy the following
# line with values fixed in this script.  E.g.
#
#  API_USER=user API_URL=https://filer-gateway/v1 API_KEY=key ${script} ${@:2}
#
${script} ${@:2}
