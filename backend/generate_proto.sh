#!/bin/bash
CURRENT_DIR=${1:-$(pwd)}

# Generate protobuf files from todo_grpc_proto directory
for x in $(find ${CURRENT_DIR}/todo_protos/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/todo_protos -I /usr/local/include --go_out=${CURRENT_DIR} \
   --go-grpc_out=require_unimplemented_servers=false:${CURRENT_DIR} ${x}/*.proto
done

# If no subdirectories, generate from proto files directly
if [ ! -d "${CURRENT_DIR}/todo_protos" ]; then
  echo "Error: todo_grpc_proto directory not found"
  exit 1
fi

# Generate from proto files in the main directory
protoc -I=${CURRENT_DIR}/todo_protos -I /usr/local/include --go_out=${CURRENT_DIR} \
 --go-grpc_out=require_unimplemented_servers=false:${CURRENT_DIR} ${CURRENT_DIR}/todo_protos/*.proto

echo "Proto generation completed!"
