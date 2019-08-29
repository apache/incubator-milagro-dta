<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# Milagro-Custody-DTA
---
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)

Milagro Custody DTA creates an ecosystem in which service providers can issue and protect secrets. When a node is connected to the network it is able to discover service providers who are able to offer secure long term storage of highly sensitive digital assets. It is written in Go and uses REST services based on the GoKit microservices framework: https://gokit.io/


## Dependencies

To correctly build the software on Ubuntu 18.04 you need to install the following packages;

```
sudo apt-get update
sudo apt-get install -y --no-install-recommends \
     ca-certificates \
     cmake \
     g++ \
     gcc \
     git \
     make \
     libtool \
     automake \
     libssl-dev
sudo apt-get clean
```

### liboqs

[liboqs](https://github.com/open-quantum-safe/liboqs) is a C library for
quantum-resistant cryptographic algorithms. It is a API level on top of the
NIST round two submissions.

```
git clone https://github.com/open-quantum-safe/liboqs.git
cd liboq
git checkout 7cb03c3ce9182790c77e69cd21a6901e270781d6 
autoreconf -i
./configure --disable-shared --disable-aes-ni --disable-kem-bike --disable-kem-frodokem --disable-kem-newhope --disable-kem-kyber --disable-sig-qtesla 
make clean
make -j
sudo make install
```

### AMCL

[AMCL](https://github.com/apache/incubator-milagro-crypto-c) is required

Build and install the AMCL library

```
git clone https://github.com/apache/incubator-milagro-crypto-c.git
cd incubator-milagro-crypto-c
mkdir build
cd build
cmake -D CMAKE_BUILD_TYPE=Release -D BUILD_SHARED_LIBS=OFF -D AMCL_CHUNK=64 -D AMCL_CURVE="BLS381,SECP256K1" -D AMCL_RSA="" -D BUILD_PYTHON=OFF -D BUILD_BLS=ON  -D BUILD_WCC=OFF -D BUILD_MPIN=OFF -D BUILD_X509=OFF -D CMAKE_INSTALL_PREFIX=/usr/local ..
make
make test
sudo make install
```

### Install pqnist

cd incubator-milagro-dta/libs/crypto/libpqnist
mkdir build
cd build
cmake -D CMAKE_INSTALL_PREFIX=/usr/local ..
make
make test
sudo make install

### golang

The code is written in golang primarily with a wrapper around some C code.

```
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
tar -xzf go1.12.linux-amd64.tar.gz
sudo cp -r go /usr/local
echo 'export PATH=$PATH:/usr/local/go/bin' >> ${HOME}/.bashrc
```

#### configure GO

```
mkdir -p ${HOME}/go/bin 
mkdir -p ${HGME}/go/pkg 
mkdir -p ${HOME}/go/src 
echo 'export GOPATH=${HOME}/go' >> ${HOME}/.bashrc 
echo 'export PATH=$GOPATH/bin:$PATH' >> ${HOME}/.bashrc
```

This package is needed for testing.

```
go get github.com/stretchr/testify/assert
```

## Run Service

This script will build the service 

```
./build.sh
```

To run the service with default settings

```
./target/service
```


## Crypto Notice

This distribution includes cryptographic software. The country in which you
currently reside may have restrictions on the import, possession, use, and/or
re-export to another country, of encryption software. BEFORE using any
encryption software, please check your country's laws, regulations and
policies concerning the import, possession, or use, and re-export of encryption
software, to see if this is permitted. See <http://www.wassenaar.org/> for
more information.

The Apache Software Foundation has classified this software as Export Commodity
Control Number (ECCN) 5D002, which includes information security software using
or performing cryptographic functions with asymmetric algorithms. The form and
manner of this Apache Software Foundation distribution makes it eligible for
export under the "publicly available" Section 742.15(b) exemption (see the BIS
Export Administration Regulations, Section 742.15(b)) for both object code and
source code.


## Disclaimer

Apache Milagro is an effort undergoing incubation at The Apache Software Foundation (ASF), sponsored by the Apache Incubator. Incubation is required of all newly accepted projects until a further review indicates that the infrastructure, communications, and decision making process have stabilized in a manner consistent with other successful ASF projects. While incubation status is not necessarily a reflection of the completeness or stability of the code, it does indicate that the project has yet to be fully endorsed by the ASF.




