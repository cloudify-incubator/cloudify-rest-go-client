# cloudify-rest-go-client
[![Circle CI](https://circleci.com/gh/cloudify-incubator/cloudify-rest-go-client/tree/master.svg?style=shield)](https://circleci.com/gh/cloudify-incubator/cloudify-rest-go-client/tree/master)

cfy-go implements CLI for cloudify client.
If we compare to official cfy command cfy-go has implementation for only external commands.

# install

```shell
sudo apt-get install golang-go
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`
go get github.com/cloudify-incubator/cloudify-rest-go-client/cfy-go
rm bin/cfy-go
ln -s src/github.com/cloudify-incubator/cloudify-rest-go-client/Makefile Makefile
make all
```

# reformat code

```shell
make reformat
```

Additional information you can check on [godoc](https://godoc.org/github.com/cloudify-incubator/cloudify-rest-go-client/cfy-go).

# create container for run python version

For use 'sudo bin/cfy-go container run -base container-place/base/  -- /usr/bin/cfy profile use local'

```shell
export ALPINE_MIRROR=http://dl-cdn.alpinelinux.org/alpine/
export ALPINE_VERSION=2.9.1-r0
export ALPINE_ROOT=alpine-root
export ALPINE_BRANCH=v3.7
mkdir ${ALPINE_ROOT}
wget ${ALPINE_MIRROR}/latest-stable/main/x86_64/apk-tools-static-${ALPINE_VERSION}.apk
tar -xzf apk-tools-static-${ALPINE_VERSION}.apk
sudo ./sbin/apk.static -X ${ALPINE_MIRROR}/latest-stable/main -U --allow-untrusted --root ${ALPINE_ROOT} --initdb add alpine-base
echo -e 'nameserver 208.67.222.222\nnameserver 2620:0:ccc::2' | sudo tee ${ALPINE_ROOT}/etc/resolv.conf
mkdir -p ${ALPINE_ROOT}/etc/apk
echo "${ALPINE_MIRROR}${ALPINE_BRANCH}/main" | sudo tee ${ALPINE_ROOT}/etc/apk/repositories

# image with python
sudo chroot ${ALPINE_ROOT}/ /sbin/apk add python2 ca-certificates py-setuptools py2-pip
sudo du -sh ${ALPINE_ROOT}

# install prerequirements for cloudify
sudo chroot ${ALPINE_ROOT}/ /sbin/apk add build-base python2-dev
sudo du -sh ${ALPINE_ROOT}

# install cloudify
sudo chroot ${ALPINE_ROOT}/ /usr/bin/pip install cloudify==4.3
sudo du -sh ${ALPINE_ROOT}

# sh to image
sudo chroot ${ALPINE_ROOT}/ /bin/sh
```
