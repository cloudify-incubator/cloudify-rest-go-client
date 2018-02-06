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
