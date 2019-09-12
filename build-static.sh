#! /bin/bash

set -e

apt-get update &&  apt-get install \
    ca-certificates git g++ gcc curl \
    make cmake automake libtool libssl-dev

CURRENT_PATH=$(pwd)
OUTPUT_PATH=$CURRENT_PATH/bin
BUILD_PATH=`mktemp -d`
LIB_SOURCE_PATH=$BUILD_PATH/src
export LIBRARY_PATH=$BUILD_PATH/lib
export C_INCLUDE_PATH=$BUILD_PATH/include
export CGO_LDFLAGS="-L $LIBRARY_PATH"
export CGO_CPPFLAGS="-I $C_INCLUDE_PATH"


echo Building Milagro Crypt C library

git clone https://github.com/apache/incubator-milagro-crypto-c.git $LIB_SOURCE_PATH/milagro-crypto-c
cd $LIB_SOURCE_PATH/milagro-crypto-c
git checkout 1.0.0
mkdir build
cd build
cmake \
  -DCMAKE_BUILD_TYPE=Release \
  -DBUILD_SHARED_LIBS=OFF \
  -DAMCL_CHUNK=64 \
  -DAMCL_CURVE="BLS381,SECP256K1" \
  -DAMCL_RSA="" \
  -DBUILD_PYTHON=OFF \
  -DBUILD_BLS=ON \
  -DBUILD_WCC=OFF \
  -DBUILD_MPIN=OFF \
  -DBUILD_X509=OFF \
  -DCMAKE_C_FLAGS="-fPIC" \
  -DCMAKE_INSTALL_PREFIX=$BUILD_PATH ..
make && make install 


echo Building LibOQS

git clone https://github.com/open-quantum-safe/liboqs.git $LIB_SOURCE_PATH/liboqs
cd $LIB_SOURCE_PATH/liboqs
git checkout 7cb03c3ce9182790c77e69cd21a6901e270781d6
autoreconf -i
./configure \
  --prefix=$BUILD_PATH \
  --disable-shared \
  --disable-aes-ni \
  --disable-kem-bike \
  --disable-kem-frodokem \
  --disable-kem-newhope \
  --disable-kem-kyber \
  --disable-sig-qtesla \
  --disable-doxygen-doc
make -j && make install


echo Building pqnist

mkdir -p $LIB_SOURCE_PATH/pqnist
cd $LIB_SOURCE_PATH/pqnist
cmake \
  -DCMAKE_BUILD_TYPE=Release\
  -DBUILD_SHARED_LIBS=OFF \
  -DCMAKE_INSTALL_PREFIX=$BUILD_PATH \
  $CURRENT_PATH/libs/crypto/libpqnist
make && make install


echo Downloading Go

curl -o "$BUILD_PATH/go.tar.gz" "https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz"
cd $BUILD_PATH && tar xzvf go.tar.gz
export GOROOT=$BUILD_PATH/go

cd $CURRENT_PATH

GO111MODULES=on \
CGO_ENABLED=1 \
$GOROOT/bin/go build -o target/out \
  -ldflags '-w -linkmode external -extldflags "-static"' \
  -o $OUTPUT_PATH/milagro \
  github.com/apache/incubator-milagro-dta/cmd/service

