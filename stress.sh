#!/bin/bash
for i in {1..10000}; do 
    curl localhost:30443/categories/1/products?=
    echo "Request: {$1}"
    sleep $1
done