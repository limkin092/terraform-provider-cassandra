# Basic terraform provider for Apache Cassandra

It will provide the following features with respect to CQL 3.0.0 spec
- Manage Keyspace(s)
- Manage Role(s)
- Managing Grants

## Initialising the provider

```java
provider "cassandra" {
  username = "cluster_username"
  password = "cluster_password"
  port     = "9042"
  hosts    = [ "localhost" ]
}
```

### Configuration

#### username

Cassandra client username. Required variable

#### password

Cassandra client password. Required variable

#### port

Cassandra client port. Default value is __9042__

#### hosts

Array of hosts pointing to nodes in the cassandra cluster

#### connection_timeout

Connection timeout to the cluster in milliseconds. Default value is __1000__

#### root_ca

Optional value, only used if you are connecting to cluster using certificates.

#### use_ssl

Optional value, it is __false__ by default. Only turned on when connecting to cluster with ssl

#### min_tls_version

Default value is __TLS1.2__. It is only applicable when use_ssl is __true__

## Resources

### Creating a Keyspace

```java
locals {
  stategy_options = {
    replication_factor = 1
  }
}

resource "cassandra_keyspace" "keyspace" {
  name                 = "some_keyspace_name"
  replication_strategy = "SimpleStrategy"
  strategy_options     = "${local.strategy_options}"
}

```

Parameters

#### name

name of the keyspace, must be between 1 and 48 characters.

#### replication_strategy

name of the replication strategy, only the built in replication strategies are supported. That is either __SimpleStrategy__ or __NetworkTopologyStrategy__

#### strategy_options
A map containing any extra options that are required by the selected replication strategy.

For simple strategy, **replication_factor** must be passed. While for network topology strategy must contain keys which corresspond to the data center names and values which match their desired replication factor

#### durable_writes

Enables or disables durable writes. The default value is __true__. It is not reccomend to turn this off.


### Creating a role

```java

resource "cassandra_role" "role" {
  name = "app_user"
  password = "sup3rS3cr3tPa$$w0rd"
}

```

Parameters

#### name

Name of the role. Must contain between 1 and 256 characters.

#### super_user

Allow role to create and manage other roles. It is __false__ by default

#### login

Enables role to be able to login. It defaults to __true__

#### password

Password for user when using cassandra internal authentication.
It is required. It has the restriction of being between 40 and 512 characters.

### Creating a Grant

```java
resource "cassandra_grant" "all_access_to_keyspace" {
  privilege = "all"
  resource_type = "keyspace"
  keyspace_name = "test"
  grantee = "migration"
}
```

Parameters

#### privilege

Type of access we are granting against a resource

One of either all, create, alter, drop, select, modify, authorize, describe and execute

See official cassandra docs for more [information](https://docs.datastax.com/en/cql/3.3/cql/cql_reference/cqlGrant.html)


#### grantee

The name of the cassandra role which we are granting privileges to

#### resource_type

Enables one to qualify/restrict the grant to a particular resource(s)

This can take any of the following values

- all functions
- all functions in keyspace
- function
- all keyspaces
- keyspace
- table
- all roles
- role
- roles
- mbean
- mbeans
- all mbeans

For more info please see official [docs](https://docs.datastax.com/en/cql/3.3/cql/cql_reference/cqlGrant.html)

#### keyspace_name

keyspace qualifier to the resource, only applicable when resource_type takes the following values

- all functions in keyspace
- function
- keyspace
- table

#### function_name

Represents name of the function we are granting access to. Its only applicable when resource_type is function


#### table_name

Represents name of the table we are granting access to. Its only applicable when resource_type is table

#### role_name

represents name of the role we are granting access to. Only applicable for resource_type is role

#### mbean_name

Represents name of the mbean we are granting access to. Only applicable for resource_type is mbean

#### mbean_pattern

Represents a pattern, which will grant access to all mbeans which satisfy this pattern. Only works when resource_type is mbeans
