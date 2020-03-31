#!/bin/bash

source input2json.sh

#trap 'echo "interrupted"' SIGINT

# curl command prefix
CURL="curl -k -#"

# filer gateway endpoint
API_URL="http://localhost:8080/v1"

function usage() {

    echo "Usage: $1 <new|get|set>" 1>&2
    cat << EOF >&2

This script is to demo managing project flexvol/qtree using the filer-gateway.
It requires "curl" and "jq".

Operations:
    new: creates new space for a project or a user home.
    get: retrieves information of the space.
    set: configures quota of the project/user home space.
EOF

}

function input() {
    prj_main=$( input2json projectID )

    prj_storage=$( input2json quotaGb system )

    data=$( jq ".storage |= ${prj_storage} | .members |= []"  <<< ${prj_main} )

    echo -n "add member [Y/n]:" > /dev/tty
    read ans && [ "$ans" == "" ] && ans="y"

    while [ "${ans,,}" == "y" ]; do
        prj_member=$( input2json userID role )
        data=$( jq ".members += [${prj_member}]" <<< ${data} )

        echo -n "add another member [y/N]:" > /dev/tty
        read ans && [ "$ans" == "" ] && ans="n"
    done

    echo $data | jq '.storage.quotaGb=(.storage.quotaGb|tonumber)' && return 0
}

function newProject() {
    data=$(input) || return 1

    echo $data | jq

    echo -n "create project [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "y" ] && 
        out=$( ${CURL} -X POST \
            -H 'content-type: application/json' \
            -d $(echo ${data} | jq -c -M '.storage.quotaGb=(.storage.quotaGb|tonumber)' ) \
            "${API_URL}/projects" )
    echo
    echo $out | jq
    echo
}

function setProject() {
    data=$(input) || return 1

    prj=$(echo $data | jq '.projectID' | sed 's/"//g')

    data=$(echo $data | jq 'with_entries(select(.key != "projectID"))')

    echo project: $prj
    echo $data | jq

    echo -n "update project [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "y" ] && 
        out=$( ${CURL} -X PATCH \
            -H 'content-type: application/json' \
            -d $(echo ${data} | jq -c -M '.storage.quotaGb=(.storage.quotaGb|tonumber)' ) \
            "${API_URL}/projects/${prj}" )
    echo
    echo $out | jq
    echo
}

## Main program
[ $# -lt 1 ] && usage && exit 1

ops=$1

case $ops in
get)
    ;;
new)
    newProject
    ;;
set)
    setProject
    ;;
esac
