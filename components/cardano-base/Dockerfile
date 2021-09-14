# Download latest base image from ubuntu
FROM debian:11-slim

# LABEL about the custom image
LABEL version="0.1"
LABEL description="This is custom Docker Image for a Cardano sandbox."

# Starting work directory
WORKDIR /root

# Disable Prompt During Packages Installation
ARG DEBIAN_FRONTEND=noninteractive

# Update and upgrade Ubuntu software repository and install various packages
# from ubuntu repository
RUN apt-get update \
    && apt-get install -y \
        automake \
        build-essential \
        pkg-config \
        make \
        g++ \
        tmux \
        git \
        jq \
        wget \
        libtool \
        autoconf \
        curl \
        libffi-dev \
        libgmp-dev \
        libssl-dev \
        libtinfo-dev \
        libsystemd-dev \
        zlib1g-dev \
        libncursesw5

# Install GHC and Cabal
RUN curl --proto '=https' --tlsv1.2 -sSf https://get-ghcup.haskell.org -sSf | sh -s -- -y

# Add ghcup bin to PATH
ENV PATH="/root/.ghcup/bin:$PATH"

# Install specific ghc version
RUN ghcup install ghc 8.10.4
RUN ghcup set ghc 8.10.4

# Make a Cardano source code directory
RUN mkdir -p ~/cardano-src

# Download and install libsodium
RUN cd ~/cardano-src && git clone https://github.com/input-output-hk/libsodium \
    && cd libsodium && git checkout 66f017f1 && ./autogen.sh && ./configure \
    && make && make install

# Add the following environment variables so the compiler can be aware 
# libsodium is installed
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"
ENV PKG_CONFIG_PATH="/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH"

# Download cardano-node
RUN cd ~/cardano-src && git clone https://github.com/input-output-hk/cardano-node.git \
    && cd cardano-node && git fetch --all --recurse-submodules --tags \
    && git checkout tags/1.29.0

# Configure cardano-node
RUN cd ~/cardano-src/cardano-node && cabal configure --with-compiler=ghc-8.10.4 \
    && echo "package cardano-crypto-praos" >>  cabal.project.local \
    && echo "  flags: -external-libsodium-vrf" >>  cabal.project.local

# Build  cardano-node
RUN cd ~/cardano-src/cardano-node && cabal build all

# Install cardano-node
RUN mkdir -p ~/.local/bin && cd ~/cardano-src/cardano-node \
    && cp -p "$(./scripts/bin-path.sh cardano-node)" ~/.local/bin/ \
    && cp -p "$(./scripts/bin-path.sh cardano-cli)" ~/.local/bin/

# Add ~/.local/bin to PATH so system can find cardano-node and cardano-cli
ENV PATH="~/.local/bin:$PATH"

SHELL ["/bin/bash", "-c"]


