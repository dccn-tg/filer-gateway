#!/bin/bash

source common.sh

# curl command prefix
CURL="curl -k -#"

# filer gateway connection detail
[ -z $API_URL ] && API_URL="http://localhost:8080/v1"
[ -z $API_KEY ] && API_KEY="demo"
[ -z $API_USER ] && API_USER="demo"

function usage() {

    echo "Usage: $1 <new|get|set>" 1>&2
    cat << EOF >&2

This script is to demo managing user qtree using the filer-gateway.
It requires "curl" and "jq".

Operations:
    new: creates new space for a user.
    get: retrieves filer information about the user.
    set: configures quota of the user's home space.
    del: delete the qtree of the user's home space.
EOF

}

function input() {
    usr_main=$( input2json userID )

    usr_storage=$( input2json quotaGb system )

    data=$( jq ".storage |= ${usr_storage}"  <<< ${usr_main} )

    echo $data | jq '.storage.quotaGb=(.storage.quotaGb|tonumber)' && return 0
}

function newUser() {
    data=$(input) || return 1

    echo $data | jq

    echo -n "create user space [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "n" ] && return 0 ||
        echo -n "password for api user ($API_USER): " > /dev/tty &&
        read -s pass &&
        out=$( ${CURL} -X POST -u "${API_USER}:${pass}" \
            -H 'content-type: application/json' \
            -H "X-API-Key: ${API_KEY}" \
            -d $(echo ${data} | jq -c -M '.storage.quotaGb=(.storage.quotaGb|tonumber)' ) \
            "${API_URL}/users" )
    echo
    echo $out | jq
    echo

    # waiting for task to reach the end state
    id=$(echo $out | jq '.taskID' | sed 's/"//g')
    [ "$id" == "null" ] && echo "cannot find task id" >&2 && return 1

    waitTask $id user
}

function setUser() {
    data=$(input) || return 1

    usr=$(echo $data | jq '.userID' | sed 's/"//g')

    data=$(echo $data | jq 'with_entries(select(.key != "userID"))')

    echo userID: $usr
    echo $data | jq

    echo -n "update user space [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "n" ] && return 0 ||
        echo -n "password for api user ($API_USER): " > /dev/tty &&
        read -s pass &&
        out=$( ${CURL} -X PATCH -u "${API_USER}:${pass}" \
            -H 'content-type: application/json' \
            -H "X-API-Key: ${API_KEY}" \
            -d $(echo ${data} | jq -c -M '.storage.quotaGb=(.storage.quotaGb|tonumber)' ) \
            "${API_URL}/users/${usr}" )
    echo
    echo $out | jq
    echo

    # waiting for task to reach the end state
    id=$(echo $out | jq '.taskID' | sed 's/"//g')
    [ "$id" == "null" ] && echo "cannot find task id" >&2 && return 1

    waitTask $id user
}

function delUser() {
    echo -n "userID: " > /dev/tty
    read ans && [ "$ans" == "" ] && return 1
    usr=$ans

    echo -n "delete user space ($usr) [y/N]:"
    read ans
    [ "${ans}" == "" ] && ans="n"
    [ "${ans,,}" == "n" ] && return 0 ||
        echo -n "password for api user ($API_USER): " > /dev/tty &&
        read -s pass &&
        out=$( ${CURL} -X DELETE -u "${API_USER}:${pass}" \
            -H 'content-type: application/json' \
            -H "X-API-Key: ${API_KEY}" \
            "${API_URL}/users/${usr}" )
    echo
    echo $out | jq
    echo

    # waiting for task to reach the end state
    id=$(echo $out | jq '.taskID' | sed 's/"//g')
    [ "$id" == "null" ] && echo "cannot find task id" >&2 && return 1

    waitTask $id user
}

function getUser() {
    echo -n "userID: " > /dev/tty
    read ans && [ "$ans" == "" ] && return 1
    $CURL -X GET "${API_URL}/users/$ans" | jq 
}

## Main program
[ $# -lt 1 ] && usage && exit 1

ops=$1

case $ops in
get)
    getUser
    ;;
new)
    newUser
    ;;
set)
    setUser
    ;;
del)
    delUser
    ;;
esac
