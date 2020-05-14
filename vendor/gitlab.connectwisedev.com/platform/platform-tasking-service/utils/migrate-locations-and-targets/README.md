Transforms tasks and task_instances tables to add timezone info for tasks and task targets.

### Prerequisites

1. Make all to get binary (platform-tasking-service-migrate-locations-and-targets by default, change Makefile if needed)
2. Generate output file

### Usage

./platform-tasking-service-migrate-locations-and-targets [options]
1) <b>config</b>:  Path to the the config file (default: "config.json") [$K_CONFIG]
2) <b>help</b>: Show help

Parameters in the config file:
* <b>CassandraURL</b>             Cassandra URL (default: "localhost:9042")
* <b>CassandraKeyspace</b>        Cassandra keyspace where migrations should be fulfilled (default: "platform_tasking_db")
