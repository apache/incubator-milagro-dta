# Milagro-Custody-DTA
---
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)

Milagro Custody DTA creates an ecosystem in which service providers can issue and protect secrets. When a node is connected to the network it is able to discover service providers who are able to offer secure long term storage of highly sensitive digital assets. It is written in Go and uses REST services based on the GoKit microservices framework: https://gokit.io/



## Linux / MacOS

Click here to down load the binary

To build the code

Clone it and run 

```
$ ./build.sh

```

## Developer Notes

You need to install protobufs

If you change the portobufs definition run 

$ protoc -I=. --go_out=. ./docs.proto

To add a new endpoint to the goKit Microservices framework

1. First define the contract in milagro/pkg/milagroservice/proto.go

Add structs for http transport
Add responses to milagro/swagger/swagger.config,yaml

2. Add and endpoint definition

milagro/pkg/milagroendpoints/endpoints.go

3. Create a handler factory





