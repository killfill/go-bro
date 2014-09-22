
# Go Bro'

![Logo](https://i.chzbgr.com/maxW500/6668368128/hB50D0768/)

GoBro is a [Cloud Foundry](http://cloudfoundry.org/) [Service Broker](http://docs.cloudfoundry.org/services/overview.html) written in [Go](http://golang.org/), that talks *Broker API v2*.

## Goals

### Be as simple as possible

It doesnt save any state, so does not depends on any database or store.
Should be easy to extend to add other services. Can be run as you like. Running it as a common CF app does work too.

It doesnt matter if the services are managed manually or by [Bosh](http://bosh.cloudfoundry.org/), clustered or not.
Just need to be accesible by the CF Controller and DEA...

### Support multiple services

The idea is to support multiple services, currently works for PostgreSQL and MySQL

It does support multiple instances of services too, so for example you can have one shared database cluster, and another dedicated one possible faster.


## Configurations

The file **config.json** contains GoBro's configuration settings, like credentials, plan limits and services data.

The file **catalog.json** contains the [catalog offerings](http://docs.cloudfoundry.org/services/catalog-metadata.html). Its a free file, and will not be parsed.

Just make sure the uuid's matches on both files.


## TODO

* Implement other services

## Usage example

```BASH
$ cf create-service-broker gobro us3r passw0rd http://10.0.0.227:3000
Creating service broker gobro as admin...
OK
```

```BASH
$ cf service-access
getting service access as admin...
broker: gobro
   service              plan            access   orgs
   postgres             free-5          none
   postgres             free-10         none
   postgres             free-20         none
   postgres             free-50         none
   postgres             free-100        none
   dedicated-postgres   dedicated-100   none
   dedicated-postgres   dedicated-200   none
   mysql                free-5          none
   mysql                free-10         none
```

```BASH
$ cf enable-service-access postgres
Enabling access to all plans of service postgres for all orgs as admin...
OK
```

```BASH
$ cf m
postgres          free-5, free-10, free-20, free-50, free-100   Shared PostgreSQL database service
```

```BASH
$ cf create-service postgres free-10 sql-common
Creating service sql-common in org virtu / space dev as admin...
OK
```

```BASH
$ cf bs myapp sql-common
Binding service sql-common to app myapp in org virtu / space dev as admin...
OK
TIP: Use 'cf restage' to ensure your env variable changes take effect
```

```BASH
$ cf us myapp sql-common
Unbinding app myapp from service sql-common in org virtu / space dev as admin...
OK
```

```BASH
$ cf ds sql-common

Really delete the service sql-common?> yes
Deleting service sql-common in org virtu / space dev as admin...
OK
```
