#!/bin/bash

functions=("fibonacci" "bbs")

for func in "${functions[@]}"; do
    mkdir -p "./output/$func"
done

for func in "${functions[@]}"; do
    for i in $(seq 1 10); do
        echo "Running $func tentative $i..."
        go run main.go "$func" > "./output/$func/exit_${i}.txt"
    done
done
