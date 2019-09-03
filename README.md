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

# Milagro Distributed Trust Authority
---
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)

Milagro D-TA is a collaborative key management server 

Milagro D-TA facilitates secure and auditable communication between people who to use key pairs (Principal) and service providers who can keep the secret keys safe (Master Fiduciary). It is written in Go and uses REST services based on the [GoKit microservices framework](https://gokit.io), it uses IPFS to create a shared immutable log of transactions and relies on Milagro-Crypto-C for it's crypto.

## Plugins
Milagro D-TA provides a basic set of services for creating identities for actors in the system, and passing encrypted communication between them but it assumes that different service providers will have their own "special sauce" for securely storing secret keys, so the vanilla services can be extended using a plugin framework. Two basic plugins are included in this release to give you an idea of how this can be done.
1. **BitcoinPlugin**  Generates a Bitcoin address and reveals the corresponding secret key
2. **SafeGuardSecret** Encrypts a string and decrypts it again

## Installation
Below are instructions on how to build and run the Milagro D-TA either using Docker, or static or dynamic builds for Linux.

These instructions will build the service with default settings including an embeded IPFS node connected to a Public IPFS network. This will get you up and running quickly but will turn your D-TA into a public IPFS relay. **Not recommended for production use!**

### Docker
To see Milagro D-TA in action you can run Milagro D-TA in a docker container.  This is currently the preferred method to build and run the D-TA.  Tested on Ubuntu 19.04 and MacOS but should run in any Docker environment.

#### Prerequisites
[Docker](https://www.docker.com/)

#### Build & Run Instructions
```
git clone https://github.com/apache/incubator-milagro-dta.git

cd incubator-milagro-dta

docker build -t mydta .

docker run -p5556:5556 mydta
```

### Linux Static Build
Note that this has only been tested on a clean build of Debian 10 and installs all required libraries statically

#### Prerequisites
[Golang](https://golang.org/doc/install)

#### Build & Run Instructions
```
sudo apt-get update
sudo apt-get upgrade
git clone https://github.com/apache/incubator-milagro-dta.git
cd incubator-milagro-dta
./build-static.sh
cd bin
./milagro init
./milagro daemon
```

### Linux Dynamic Build
Note that the static build above doesn't work on Ubuntu, so instead please use these instructions.  Note that these have only been tested on a clean build of Ubuntu 18.04.

#### Prerequisites
To correctly build the software on Ubuntu 18.04 you need to install the following packages:

#### Development Tools
```
sudo apt-get update
sudo apt-get install \
     ca-certificates \
     cmake \
     g++ \
     gcc \
     git \
     make \
     libtool \
     automake \
     libssl-dev \
     jq \
     curl
```

#### liboqs

[liboqs](https://github.com/open-quantum-safe/liboqs) is a C library for
quantum-resistant cryptographic algorithms. It is a API level on top of the
NIST round two submissions.

```
git clone https://github.com/open-quantum-safe/liboqs.git
cd liboqs
git checkout 7cb03c3ce9182790c77e69cd21a6901e270781d6 
autoreconf -i
./configure --disable-shared --disable-aes-ni --disable-kem-bike --disable-kem-frodokem --disable-kem-newhope --disable-kem-kyber --disable-sig-qtesla 
make clean
make -j
sudo make install
```

#### AMCL

[AMCL](https://github.com/apache/incubator-milagro-crypto-c) is required

Build and install the AMCL library

```
git clone https://github.com/apache/incubator-milagro-crypto-c.git
cd incubator-milagro-crypto-c
mkdir build
cd build
cmake -D CMAKE_BUILD_TYPE=Release -D BUILD_SHARED_LIBS=ON -D AMCL_CHUNK=64 -D AMCL_CURVE="BLS381,SECP256K1" -D AMCL_RSA="" -D BUILD_PYTHON=OFF -D BUILD_BLS=ON -D BUILD_WCC=OFF -D BUILD_MPIN=OFF -D BUILD_X509=OFF -D CMAKE_INSTALL_PREFIX=/usr/local -DCMAKE_C_FLAGS="-fPIC" ..
make
make test
sudo make install
```

#### pqnist

```
git clone https://github.com/apache/incubator-milagro-dta.git
cd incubator-milagro-dta/libs/crypto/libpqnist
mkdir build
cd build
cmake -D CMAKE_INSTALL_PREFIX=/usr/local -D BUILD_SHARED_LIBS=ON ..
make
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
make test
sudo make install
```

#### golang

Download and install [Golang](https://golang.org/dl/)


#### Set the library path
```
export C_INCLUDE_PATH=$C_INCLUDE_PATH:/usr/local/lib
```

#### Build & Run Instructions
```
./build.sh
```

To run the service with default settings:

```
./target/milagro init
./target/milagro daemon 
```

### MacOS 

These instructions have been tested on MaxOS Mojave 10.14.6

#### Prerequisites
To correctly build the software on MacOS you need to install the following packages:

#### Brew
To allow the additional packages below to be installed, [Brew](https://docs.brew.sh/Installation) is required.

#### Development Tools
```
brew install \
     cmake \
     autoconf \
     automake \
     libtool 
```

#### liboqs

[liboqs](https://github.com/open-quantum-safe/liboqs) is a C library for
quantum-resistant cryptographic algorithms. It is a API level on top of the
NIST round two submissions.

```
git clone https://github.com/open-quantum-safe/liboqs.git
cd liboqs
git checkout 7cb03c3ce9182790c77e69cd21a6901e270781d6 
autoreconf -i
./configure --disable-shared --disable-aes-ni --disable-kem-bike --disable-kem-frodokem --disable-kem-newhope --disable-kem-kyber --disable-sig-qtesla 
make clean
make -j
sudo make install
```

#### AMCL

[AMCL](https://github.com/apache/incubator-milagro-crypto-c) is required

Build and install the AMCL library

```
git clone https://github.com/apache/incubator-milagro-crypto-c.git
cd incubator-milagro-crypto-c
mkdir build
cd build
cmake -D CMAKE_BUILD_TYPE=Release -D BUILD_SHARED_LIBS=ON -D AMCL_CHUNK=64 -D AMCL_CURVE="BLS381,SECP256K1" -D AMCL_RSA="" -D BUILD_PYTHON=OFF -D BUILD_BLS=ON -D BUILD_WCC=OFF -D BUILD_MPIN=OFF -D BUILD_X509=OFF -D CMAKE_INSTALL_PREFIX=/usr/local -DCMAKE_C_FLAGS="-fPIC" ..
make
make test
sudo make install
```

#### pqnist

```
git clone https://github.com/apache/incubator-milagro-dta.git
cd incubator-milagro-dta/libs/crypto/libpqnist
mkdir build
cd build
cmake -D CMAKE_INSTALL_PREFIX=/usr/local -D BUILD_SHARED_LIBS=ON ..
make
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
make test
sudo make install
```

#### golang

Download and install [Golang](https://golang.org/dl/)


#### Set the library path
```
export C_INCLUDE_PATH=$C_INCLUDE_PATH:/usr/local/lib
```

#### Build & Run Instructions
```
./build.sh
```

To run the service with default settings:

```
./target/milagro init
./target/milagro daemon 
```

## Documentation

You can find documentation for Milagro D-TA in the main [Milagro docs site](https://milagro.apache.org/) 

Which includes a quick start guide that will show you how to get Milagro D-TA to "do stuff"


## Contributing

 Key pairs are becoming central to our online lives, and keeping secret keys safe is a growing industry, we hope to create an ecosystem of custodial service providers who collaborate to make the Internet a safer place for everyone. We are keen to get contributions and feedback from anyone in this space. This is a brand new project so our development processes are still being figured out, but if you have suggestions, questions or wish to make contributions please go ahead raise an issue and someone on the team will get right on it.


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

## Credits

* Design and Concept... [Brian Spector](https://github.com/spector-in-london)
* Core Algorithm and Services... [Chris Morris](https://github.com/fluidjax)
* Framework and Refactoring... [Stanislav Mihaylov](https://github.com/smihaylov)
* Crypto Genius... [Kealan McCusker](https://github.com/kealan)
* Keeper of "The Apache Way"... [John McCane-Whitney](https://github.com/johnmcw)
* Prototype and Cat Herding... [Howard Kitto](https://github.com/howardkitto)