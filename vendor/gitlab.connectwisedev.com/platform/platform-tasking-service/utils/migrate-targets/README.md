Transforms column targets with type set<frozen<target>> to the column task_targets with type map<uuid, boolean> in the table platform_tasking_db.tasks.

### Prerequisites

1. Make all to get binary (platform-tasking-service-migrate-targets by default, change Makefile if needed)
2. Generate output file

### Usage

./platform-tasking-service-migrate-targets [options]
1) <b>config</b>:  Path to the the config file (default: "config.json") [$K_CONFIG]
2) <b>help</b>: Show help

Parameters in the config file:
* <b>CassandraURL</b>             Cassandra URL (default: "localhost:9042")
* <b>CassandraKeyspace</b>        Cassandra keyspace where migrations should be fulfilled (default: "platform_tasking_db")
