# DbMux - *database multiplexer*
_Distributed reverse proxy dbms load balancing multiplexer_

## A short explanation
DbMux is a database reverse proxy multiplexer which allows you to proxy your mysql db connections to multiple mysql instances. The final goal for this project is to create a tool which can assist with TCP level query caching, database sharding, database syncing and connection multiplexing. 

### What is TCP caching?

The basic idea behind TCP level caching is to 'query' data on the TCP level for a specified amount of time. The tables which you would like to have cached can be set in the `appconfig.json` file. The advantages of this strategy is to allow caching to multiple client connections without the need for configuring third-party tools (like redis) on every endpoint where caching is required.
### What is multiplexing?
The connection multiplexer is responsible for proxying / tunneling database tcp connections to other instances for real-time data replication and potentially database sharding.
### What is Sharding?
Database sharding is when you distribute a single data set across multiple databases, which could be running on different machines. DbMux would be responsible for managing those connections and load-balancing the requests made to each database instance.
### What is database syncing?
Database syncing can be configured to sync the datasets of multiple database instances. This feature can be used with or without database sharding activated.

## The core idea behind db-mux

Provide large and small users alike with a means to easily be able to solve some of the biggest pain-points in large-scale distributed systems development.

## Supported databases
- MySql

(Hopefully others will be added soonish)

## Why GoLang?
Go is a relatively low-level language and handles concurrency quite well, as well as being easy to read and understand. Go is also quite popular for these types of applications because of those qualities. Examples include 
The short and sweet is that Go is great for handling concurrency. In large scale production environments that is very important, if we are trying to manage multiple concurrent database requests. It is also really quick and responding to TCP/IP & HTTP requests which also makes it a great option.

## Development plan

### Short term goals
  -[X] Proxy TCP connections for mysql.
  -[ ] Allow for the caching of queries by hashing and caching the query response in memory.
  -[ ] Replicate inserts onto two different versions of mysql.

### Long term goals

  -[ ] Implement sharding and partitioning strategies.
  -[ ] Load-balance database requests between multiple instances.
