#!/bin/bash

function input2json() {

    keys=($@)

    jq_args=()
    jq_query='.'

    for idx in ${!keys[@]}; do
        echo -n "${keys[$idx]}: " > /dev/tty
        read val

        jq_args+=( --arg "key$idx" "${keys[$idx]}" )
        jq_args+=( --arg "val$idx" "$val" )
        jq_query+=" | .[\$key${idx}]=\$val${idx}"
    done

    jq "${jq_args[@]}" "$jq_query" <<<'{}'
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
