# Milagro-DTA
---
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)



Milagro Distributed Trust Authority creates an ecosystem in which service providers can issue and protect secrets. When a node is connected to the network it will discover service providers who are able to offer secure longterm storage of highly sensitive digital assets. 

It is written in Go and uses Qredo's implmentation of the GoKit microservices framework: https://gokit.io/

## Installation

The Milagro DTA has two major dependedncies

### Redis

Instructions


### IPFS

Instructions

export IPFS_PATH=$HOME/.ipfs-custody


## Linux / MacOS

Click here to download the binary

To build the code

Clone it and run 

```
$ go get ./..
```

## Integration

Milagro-DTA is designed to be used by any organisation that might have digital assets such as keys and tokens that need to be secured, its is designed to be integrated in your existing back-office systems via a simple REST api.[The specification of which can be seen here](/swagger)

## Developer Notes

If you would like further information about how to integrate or extend Milagro-DTA the docs can be found here [https://milagro.apache.org/docs/d-ta-overview/](https://milagro.apache.org/docs/d-ta-overview/)

You may need to install protobufs

If you change the portobufs definition run 

```
$ protoc -I=. --go_out=. ./docs.proto
```




