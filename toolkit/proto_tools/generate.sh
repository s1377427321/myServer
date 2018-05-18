#!/bin/bash
rm -rf ../../protocol/*
rm -rf ../../src/protocol/*
python ./merge.py

cd ../../protocol
protoc --go_out=../src/protocol/  ./*.proto
#./protoc --go_out=../../protocol/  --proto_path=../../../protocol/*.proto
#read -p "Press any key to continue." var
