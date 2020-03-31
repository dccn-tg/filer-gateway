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
