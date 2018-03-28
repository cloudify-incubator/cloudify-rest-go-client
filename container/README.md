Solution for containerized python:

micropython-bin - shared linked, 3.5 but is not fully compatible with CPython, 450kb
Source: https://github.com/micropython/micropython
Documentation: http://docs.micropython.org/en/latest/pyboard/reference/index.html
Create script:
```shell
git clone https://github.com/micropython/micropython.git
cd micropython/
git submodule update --init
cd ports/unix
make axtls
make
```

Debian/Debootstrap - full operation system:
Based on https://wiki.debian.org/Debootstrap
```shell
sudo debootstrap --variant=minbase stable stable-chroot http://deb.debian.org/debian/
```

Size without python: 172M
Size with after "apt update": 248M
Size with python: 284M
Size with python without update: 254M
Size bz2: 215M

Alpine - full operation system:
Based on https://wiki.alpinelinux.org/wiki/Installing_Alpine_Linux_in_a_chroot
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
Size without python: 7,6M
Size with python: 68M
Size with cloudify: 266M

Ok, lets stop on alpine.

# For run cloudify we need some devices

Like if no urandom: NotImplementedError: /dev/urandom (or equivalent) not found
```shell
export ALPINE_ROOT=alpine-root
sudo mknod -m 666 ${ALPINE_ROOT}/dev/full c 1 7
sudo mknod -m 666 ${ALPINE_ROOT}/dev/ptmx c 5 2
sudo mknod -m 644 ${ALPINE_ROOT}/dev/random c 1 8
sudo mknod -m 644 ${ALPINE_ROOT}/dev/urandom c 1 9
sudo mknod -m 666 ${ALPINE_ROOT}/dev/zero c 1 5
sudo mknod -m 666 ${ALPINE_ROOT}/dev/tty c 5 0
```
