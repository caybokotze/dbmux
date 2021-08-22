# DbMux - *database multiplexer*
_Distributed reverse proxy dbms load balancing multiplexer_

## What does DbMux currently do?
DbMux currently can act as a proxy for your mysql connection. So if you want to tunnel your mysql connection through from 3306 to 3308 on the server you are running mysql, you can use dbmux to do that.

## What do I want the project to become?
1.) A tool for database replication. I have looked into this quite a bit and I think this service could be of use for workloads that require an element of database replication be that for redundancy, security or anything else. For me the advantage of doing this over TCP means you can achive this very quickly. You can replicate a query/command to 2 or multiple databases, even across multiple server environments.

2.) Database sharding. If the db proxy acts as the ow-level middleman (over TCP) between the client and the database there is quite a lot of granular contorl this would allow you to have. You can hash common queries and target several mysql connections on any other environment.

3.) Managing distributed workloads. More than just sharding and replication, if a single service is reposible for the connection between multiple databases you could gain significant performance improvements across multi-server environments.

4.) Caching. The db-mux service can manage your database caching on the TCP level and do so very quickly for frequently used queries.

## My vision for dbmux.
I beleive DbMux could be a useful tool for companies where managing diverse and complex workloads can become difficult to manage. It could be a great option for comapnies that want to make use of large multi-cluster environments and need something that can manage that database workload automatically and very effectively. I think dbmux can become that tool with enough time and effort spent. It could solve some of those big complicated problems without thinking about it too much. Caching, Sharding, Replciation and shared state between multiple services.

## Why is this written in Go?
My short explination is becuase Go handles concurrency quite well and this will be handy when managing multiple SQL Clients at the same time and managing that workload effectively. It also makes it very easy to manage TCP clients / connections. It is relatively low level and a quick programming langauge and serves this type of use case quite well. 

## Development plan
- [X] TCP Proxy (Currently works without any major issues but tests still need to be written)
- [ ] Reading / Accessing Commands and queries from the SQL Client.
- [ ] Caching db queries / commands
- [ ] Managing database sharding and partitioned environments effectively.
- [ ] Managing multi cluster environements with many SQL clients and replicating the tool accross multiple server environments.
- [ ] Supporting more than just MySql, but perhaps also Oracle, MSSql and Postgres.
