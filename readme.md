# dbmux
_Distributed reverse proxy dbms load balancing multiplexer_

## A short explanation
DbMux is a database reverse proxy. (Sort of) DbMux is a tool to manage distributed clustered database environments, which could be well suited for large scale distributed systems. It offers you the ability to setup and manage your own database clusters and to replicate data across different data sources. It also acts as a way to offer a single source of truth in distributed application architectures and will do some work to speed up frequest database queries with the help of Redis.

## The idea behind dbmux

Instead of targeting your db of choice directly, what you effectively do is execute queries against the service instead as a database proxy. From there you can then manage what actually happens with the request to the db. The service also has built in support for Redis, so that could help speed up queries that are made very frequently.

## Supported databases
- MySql
- MsSql
- Postgres

## Why is this written in Go?
The short and sweet is that Go is great for handling concurrency. In large scale production environments that is very important, if we are trying to manage multiple concurrent database requests. It is also really quick and responding to TCP/IP & HTTP requests which also makes it a great option.

## Development plan

### Short term vision

- Build out TCP client to listen for incoming connections
  - Proxy that connection to the actual mysql client that is running locally.
  
- Setup the code that would work as the server, and code that would operate as the client. The server is responsible for accepting TCP connections from the host, and the clients listen for TCP connections from the server. The host is whatever application is making use of the service. A single instance of a ramjet service can act as a server and client, which would probably be the most common configuration.
- Ramjet services can ping each other to find out what would be the fastest route for the request to be handled. According to that ping chart the replication can take place.

### Long term vision

- Tunnel database requests through a redis cache instance and validate whether the request has been made before. If it is, then the request can be accepted.
- View the modes that are connected to each other on a browser interface to be able to view the health of the replication cluster.
- Look at ways of handling database sharding and partitioning for db replication.
