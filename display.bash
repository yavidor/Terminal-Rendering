#!/bin/bash
transmit_png() {
    data=$(base64 "$1")
    data="${data//[[:space:]]}"
    local pos=0
    local chunk_size=4096
    while [ $pos -lt ${#data} ]; do
        printf "\e_G"
        [ $pos = "0" ] && printf "a=T,f=100,"
        local chunk="${data:$pos:$chunk_size}"
        pos=$(($pos+$chunk_size))
        [ $pos -lt ${#data} ] && printf "m=1"
        [ ${#chunk} -gt 0 ] && printf ";%s" "${chunk}"
        printf "\e\\"
    done
    printf "aaa"
}

transmit_png "$1"
