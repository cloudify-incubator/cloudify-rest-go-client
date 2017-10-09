# cloudify-rest-go-client

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
# Functionlity from original cfy client

* Common parameters:
    * `-host`: manager host
    * `-user`: manager user
    * `-password`: manager password
    * `-tenant`: manager tenant
* Example:

```shell
cfy-go status version -host <your manager host> -user admin -password secret -tenant default_tenant
```

## agents
Handle a deployment's agents
* Not Implemented

------

## kubernetes related commands

### init
Return json in kubernetes format for use as init script responce

```shell
cfy-go kubernetes init
```

### mount
Return json in kubernetes format for use as mount script responce

```shell
cfy-go kubernetes mount /tmp/someunxists '{"kubernetes.io/fsType":"ext4","kubernetes.io/pod.name":"nginx","kubernetes.io/pod.namespace":"default","kubernetes.io/pod.uid":"ecd89d9d-a44a-11e7-b34f-00505685ddd0","kubernetes.io/pvOrVolumeName":"someunxists","kubernetes.io/readwrite":"rw","kubernetes.io/serviceAccount.name":"default","size":"1000m","volumeID":"vol1","volumegroup":"kube_vg"}' -deployment slave -instance kubenetes_slave_*
```

### unmount
Return json in kubernetes format for use as unmount script responce

```shell
cfy-go kubernetes unmount /tmp/someunxists -deployment slave -instance kubenetes_slave_*
```

------

## blueprints
Handle blueprints on the manager

### create-requirements
Create pip-requirements
* Not Implemented

### delete
Delete a blueprint [manager only]

```shell
cfy-go blueprints delete blueprint
```

### download
Download a blueprint [manager only]

```shell
cfy-go blueprints download blueprint
```

### get
Retrieve blueprint information [manager only]

```shell
cfy-go blueprints list -blueprint blueprint
```

### inputs
Retrieve blueprint inputs [manager only]
* Not Implemented

### install-plugins
Install plugins [locally]
* Not Implemented

### list
List blueprints [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

```shell
cfy-go blueprints list
```

### package
Create a blueprint archive
* Not Implemented

### upload
Upload a blueprint [manager only]

```shell
cfy-go blueprints upload new-blueprint -path src/github.com/cloudify-incubator/cloudify-rest-go-client/examples/blueprint/Minimal.yaml
```

### validate
Validate a blueprint
* Not Implemented

------

## bootstrap
Bootstrap a manager
* Not Implemented

------

## cluster
Handle the Cloudify Manager cluster
* Not Implemented

------

## deployments
Handle deployments on the Manager

### create
Create a deployment [manager only]
* Partially implemented, you can set inputs only as json string.

```shell
cfy-go deployments create deployment  -blueprint blueprint --inputs '{"ip": "b"}'
```

### delete
Delete a deployment [manager only]

```shell
cfy-go deployments delete  deployment
```

### inputs
Show deployment inputs [manager only]
* Not Implemented

### list
List deployments [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

```shell
cfy-go deployments list
```

### outputs
Show deployment outputs [manager only]

```shell
cfy-go deployments inputs -deployment deployment
```

### update
Update a deployment [manager only]
* Not Implemented

------

## dev
Run fabric tasks [manager only]
* Not Implemented

------

## events
Show events from workflow executions

### delete
Delete deployment events [manager only]
* Not Implemented

### list
List deployments events [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

Supported filters:
* `blueprint`: The unique identifier for the blueprint
* `deployment`: The unique identifier for the deployment
* `execution`: The unique identifier for the execution

```shell
cfy-go events list
```

------

## executions
Handle workflow executions

### cancel
Cancel a workflow execution [manager only]
* Not Implemented

### get
Retrieve execution information [manager only]
* Not Implemented

### list
List deployment executions [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

```shell
cfy-go executions list
cfy-go executions list -deployment deployment

```

### start
Execute a workflow [manager only]
* Partially implemented, you can set params only as json string.

```shell
cfy-go executions start uninstall -deployment deployment
```

------

## groups
Handle deployment groups
* Not Implemented

------

## init
Initialize a working env
* Not Implemented

------

## install
Install an application blueprint [manager only]
* Not Implemented

------

## ldap
Set LDAP authenticator.
* Not Implemented

------

## logs
Handle manager service logs
* Not Implemented

------

## maintenance-mode
Handle the manager's maintenance-mode
* Not Implemented

------

## node-instances
Handle a deployment's node-instances

### get
Retrieve node-instance information [manager only]

```shell
cfy-go node-instances list -deployment deployment
```

### list
List node-instances for a deployment [manager only]

```shell
cfy-go node-instances list -deployment deployment
```

------

## nodes
Handle a deployment's nodes

### get
Retrieve node information [manager only]

```shell
cfy-go nodes list -node server -deployment deployment
```

### list
List nodes for a deployment [manager only]

```shell
cfy-go nodes list
```

------

## plugins
Handle plugins on the manager

### delete
Delete a plugin [manager only]

* Not Implemented

### download
Download a plugin [manager only]

* Not Implemented

### get
Retrieve plugin information [manager only]
* Not Implemented

### list
List plugins [manager only]
```shell
cfy-go plugins list
```

### upload
Upload a plugin [manager only]

* Not Implemented

### validate
Validate a plugin

* Not Implemented (requered )

------

## profiles
Handle Cloudify CLI profiles Each profile can...
* Not Implemented

------

## rollback
Rollback a manager to a previous version
* Not Implemented

------

## secrets
Handle Cloudify secrets (key-value pairs)
* Not Implemented

------

## snapshots
Handle manager snapshots
* Not Implemented

------

## ssh
Connect using SSH [manager only]
* Not Implemented

------

## status
Show manager status [manager only]

### Manager state
Show service list on manager

```shell
cfy-go status state
```

### Manager version
Show manager version

```shell
cfy-go status version
```

------

## teardown
Teardown a manager [manager only]
* Not Implemented

------

## tenants
Handle Cloudify tenants (Premium feature)
* Not Implemented

------

## uninstall
Uninstall an application blueprint [manager only]
* Not Implemented

------

## user-groups
Handle Cloudify user groups (Premium feature)
* Not Implemented

------

## users
Handle Cloudify users
* Not Implemented

------

## workflows
Handle deployment workflows
* Not Implemented
