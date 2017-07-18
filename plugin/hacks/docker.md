# Short, Unsupported guide to do this using docker

* *DISCLAIMER*: Seriously, this is a unsupported. Don't email me.

You don't have linux?

* Install `docker` in your system (provided your system supports `docker`).

* You need an image, As of 2017.07.17, `golang:1.8.3` should suffice. I did mine based on ubuntu, here is the `Dockerfile` #YMMV.

```
FROM ubuntu:16.04

MAINTAINER Herman Junge "chpdg42@gmail.com"

RUN apt-get update && \
    apt-get install -y \
        curl \
        git \
        make \
        vim-gnome \
        build-essential && \
    apt-get upgrade -y

WORKDIR /tmp

RUN curl -O https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && \
    tar -xvf go1.8.3.linux-amd64.tar.gz && \
    GOROOT=/tmp/go GOPATH=/tmp/gopath /tmp/go/bin/go get golang.org/x/tools/cmd/goimports && \
    mv /tmp/go /usr/local && \
    mv /tmp/gopath/bin/goimports /usr/local/bin/goimports

RUN echo 'export GOROOT=/usr/local/go' >> /root/.bashrc \
    && echo 'export GOPATH=/go' >> /root/.bashrc \
    && echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /root/.bashrc
```

* Create the image (will take a while) with

```
docker build -t golang-ubuntu:1.8.3 -f Dockerfile .
```

* Now, to do the less damage possible, we will run this image, only mounting `github.com/ipfs/go-ipfs` and `github.com/ipfs/go-ipld/eth`.

```
docker run \
	-ti --rm \
	--name golang-linux-ipfs-compiler \
	-v $HOME/.golang-linux-ipfs-compiler:/go \
	-v $GOPATH/src/github.com/ipfs/go-ipfs:/go/src/github.com/ipfs/go-ipfs \
	-v $GOPATH/src/github.com/ipfs/go-ipld-eth:/go/src/github.com/ipfs/go-ipld-eth \
	-w /go/src/github.com/ipfs \
	ubuntu-golang:1.8.3 /bin/bash
```

This will mount your directories and put _everything_ you add into the directory `~/.golang-linux-ipfs-compiler`

* You should be inside the container. Now comes the _damaging part_. Please go into the `go-ipfs` directory and delete `gx` and `gx-go`.

```
# We start at /go/src/github.com/ipfs

cd go-ipfs

rm -rf bin/gx*
```

Why? Because they may have been compiled for your host system.

* Install `IPFS`. (This really can take a while. Undust Skyrim, grind your character to Destruction 100, it rules).

```
make install
```

* Finally. Build your plugin

```
cd /go/src/github.com/ipfs/go-ipld-eth/plugin

make
```

* And move the `ethereum.so` file into the IPFS' `plugins` directory

```
# Make it an executable
chmod +x ethereum.so

# Did you ipfs init before?
ipfs init

# You may need to mkdir this one
mkdir /root/.ipfs/plugins

# Do it
mv /go/src/github.com/ipfs/go-ipld-etc/plugin/ethereum.so /root/.ipfs/plugins/.
```

* Profit!

```
# In one console
ipfs daemon

# In a different one
cd go/src/github.com/ipfs/go-ipld-etc
cat test_data/block-with-txs.bin | ipfs dag put --input-enc=raw --format=eth
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw

# You should get

{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","difficulty":12555463106190,"extra":"14MBAwOER2V0aIdnbzEuNC4yhWxpbnV4","gaslimit":3141592,"gasused":231000,"mixdigest":"0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0","nonce":"0xf491f46b60fe04b3","number":999999,"parent":{"/":"z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3"},"receipts":{"/":"z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhGQtNDNDk9m9N2BZA"},"root":{"/":"z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMpoh"},"time":1455404037,"tx":{"/":"z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq3Dc51"},"uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```
