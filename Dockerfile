# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

FROM ubuntu:latest as libs_builder

RUN apt-get update &&  apt-get install -y --no-install-recommends \
    ca-certificates \
    cmake \
    g++ \
    gcc \
    git \
    make \
    libtool \
    automake \
    libssl-dev

ENV BUILD_PATH=/tmp/milagro-dta-build
ENV LIBRARY_PATH=$BUILD_PATH/lib
ENV C_INCLUDE_PATH=$BUILD_PATH/include

WORKDIR /root

# Milagro Crypto C Library
RUN echo Building Milagro Crypt C library && \
	git clone https://github.com/apache/incubator-milagro-crypto-c.git && \
	cd incubator-milagro-crypto-c && \
	git checkout feecb27e7cc24ba5809405cc3efe74b0f7ac69ef && \
    mkdir build && \
    cd build && \
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
    -DCMAKE_INSTALL_PREFIX=$BUILD_PATH \
    .. && \
    make && make install 

# LibOQS
RUN echo Building LibOQS && \
	git clone https://github.com/open-quantum-safe/liboqs.git && \
	cd liboqs && \
    git checkout 7cb03c3ce9182790c77e69cd21a6901e270781d6 && \
    autoreconf -i && \
    ./configure \
    --prefix=$BUILD_PATH \
    --disable-shared \
    --disable-aes-ni \
    --disable-kem-bike \
    --disable-kem-frodokem \
    --disable-kem-newhope \
    --disable-kem-kyber \
    --disable-sig-qtesla \
    --disable-doxygen-doc && \
    make -j && make install


# Lib pqnist
ADD libs/crypto/libpqnist pqnist/
RUN mkdir -p pqnist/build && \
	cd pqnist/build && \
	cmake \
	-DCMAKE_BUILD_TYPE=Release\
	-DBUILD_SHARED_LIBS=OFF \
    -DCMAKE_INSTALL_PREFIX=$BUILD_PATH \
	.. && \
	make && make install


FROM golang:1.13 as go_builder

ENV LIBS_PATH=/tmp/milagro-dta-build
ENV LIBRARY_PATH=$LIBS_PATH/lib
ENV C_INCLUDE_PATH=$LIBS_PATH/include
ENV PROJECT_PATH=/src/github.com/apache/incubator-milagro-dta
ENV CGO_LDFLAGS="-L $LIBRARY_PATH"
ENV CGO_CPPFLAGS="-I $C_INCLUDE_PATH"

COPY --from=libs_builder $LIBS_PATH $LIBS_PATH
ADD . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN CGO_ENABLED=1 \
    GO111MODULES=on \
    go build \
      -ldflags '-w -linkmode external -extldflags "-static"' \
      -o $GOPATH/bin/milagro \
      github.com/apache/incubator-milagro-dta/cmd/service

RUN $GOPATH/bin/milagro init

FROM alpine as certs
RUN apk add --no-cache ca-certificates


FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go_builder /root/.milagro .milagro
COPY --from=go_builder /go/bin/milagro /

ENTRYPOINT ["/milagro", "daemon"]
