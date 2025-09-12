#!/bin/bash

source common.sh

#trap 'echo "interrupted"' SIGINT

# curl command prefix
CURL="curl -k -#"

# filer gateway connection detail
[ -z $API_URL ] && API_URL="http://localhost:8080/v1"
[ -z $API_KEY ] && API_KEY="demo"
[ -z $API_USER ] && API_USER="demo"

function usage() {

    echo "Usage: $1 <new|get|set>" 1>&2
    cat << EOF >&2

This script is to demo managing research-related data (RRD) flexvol/qtree using the filer-gateway.
It requires "curl" and "jq".

Operations:
    new: creates new space for a project rrd storage.
    get: retrieves information of a project rrd storage.
    set: configures quota of the project rrd storage.
EOF

}

function input() {
    prj_main=$( input2json projectID )

    data=$( jq ".members |= []"  <<< ${prj_main} )

    echo -n "add member [Y/n]:" > /dev/tty
    read ans && [ "$ans" == "" ] && ans="y"

    cmembers=0
    while [ "${ans,,}" == "y" ]; do
        prj_member=$( input2json userID role )
        data=$( jq ".members += [${prj_member}]" <<< ${data} )
	cmembers=$(($cmembers + 1))

        echo -n "add another member [y/N]:" > /dev/tty
        read ans && [ "$ans" == "" ] && ans="n"
    done

    if [ $cmembers -gt 0 ]; then
        echo -n "add members recursively [y/N]:" > /dev/tty
        read ans && [ "$ans" == "" ] && ans="n"
        [ "${ans,,}" == "y" ] && data=$( jq ".recursion = true" <<< ${data} )
    fi

    echo $data | jq && return 0
}

function newRrd() {
    data=$(input) || return 1

    echo $data | jq

    echo -n "create project RRD [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "n" ] && return 0 ||
        echo -n "password for api user ($API_USER): " > /dev/tty &&
        read -s pass &&
        out=$( ${CURL} -X POST -u "${API_USER}:${pass}" \
            -H 'content-type: application/json' \
            -H "X-API-Key: ${API_KEY}" \
            -d $(echo ${data} | jq -c -M ) \
            "${API_URL}/rrds" )
    echo
    echo $out | jq
    echo

    # waiting for task to reach the end state
    id=$(echo $out | jq '.taskID' | sed 's/"//g')
    [ "$id" == "null" ] && echo "cannot find task id" >&2 && return 1

    waitTask $id project
}

function setRrd() {
    data=$(input) || return 1

    prj=$(echo $data | jq '.projectID' | sed 's/"//g')

    data=$(echo $data | jq 'with_entries(select(.key != "projectID"))')

    echo project: $prj
    echo $data | jq

    echo -n "update project [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "n" ] && return 0 ||
        echo -n "password for api user ($API_USER): " > /dev/tty &&
        read -s pass &&
        out=$( ${CURL} -X PATCH -u "${API_USER}:${pass}" \
            -H 'content-type: application/json' \
            -H "X-API-Key: ${API_KEY}" \
            -d $(echo ${data} | jq -c -M ) \
            "${API_URL}/rrds/${prj}" )
    echo
    echo $out | jq
    echo

    # waiting for task to reach the end state
    id=$(echo $out | jq '.taskID' | sed 's/"//g')
    [ "$id" == "null" ] && echo "cannot find task id" >&2 && return 1

    waitTask $id project
}

function getRrd() {
    echo -n "projectID: " > /dev/tty
    read ans && [ "$ans" == "" ] && return 1
    $CURL -X GET "${API_URL}/rrds/$ans" | jq 
}

function waitTask() {
    id=$1
    ns=$2
    while [ 1 -eq 1 ]; do
        s=$(taskPoll $id $ns)
        [ $? -ne 0 ] && break
        echo "task $s" > /dev/tty
        if [[ "$s" =~ ^(failed|succeeded)$ ]]; then
            break
        else
            sleep 2
        fi
    done
}

function taskPoll() {

    id=$1
    ns=$2

    out=$( ${CURL} -X GET "${API_URL}/tasks/${ns}/${id}" )
    [ $? -ne 0 ] && echo "fail to poll task $id" >&2 && return 1

    status=$( echo $out | jq '.taskStatus.status' | sed 's/"//g' )

    if [[ "$status" =~ ^(waiting|processing|failed|succeeded|canceled)$ ]]; then
        echo $status
        return 0
    else
        echo $out >&2
        return 1
    fi
}

## Main program
[ $# -lt 1 ] && usage && exit 1

ops=$1

case $ops in
get)
    getRrd
    ;;
new)
    newRrd
    ;;
set)
    setRrd
    ;;
esac
