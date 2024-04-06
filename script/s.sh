#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <proto_path> <output_file> <filename> <destination_folder>"
    exit 1
fi

proto_path="$1"
output_file="$2"
filename="$3"
destination_folder="$4"

protoc --proto_path="$proto_path" \
    --plugin=protoc-gen-ts_proto=../user-svc/node_modules/.bin/protoc-gen-ts_proto \
    --ts_proto_out="$output_file" \
    --ts_proto_opt=nestJs=true "$filename.proto"

cp "$proto_path/$filename.ts" "$destination_folder"
cp "$proto_path/$filename.proto" "$destination_folder"