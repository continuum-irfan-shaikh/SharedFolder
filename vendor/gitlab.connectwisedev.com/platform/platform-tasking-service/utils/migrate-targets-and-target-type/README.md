This util is develop to set target_type and targets to the proper values for all existing tasks

### Prerequisites

1. Make all to get binary (platform-tasking-service-migrate-targets-and-target-type by default, change Makefile if needed)
2. Generate output file

### Usage

./platform-tasking-service-migrate-targets-and-target-type [options]
1) <b>config</b>:  Path to the config file (default: "config.json") [$K_CONFIG]
2) <b>help</b>: Show help

Parameters in the config file:
* <b>CassandraURL</b>             Cassandra URL (default: "localhost:9042")
* <b>CassandraKeyspace</b>        Cassandra keyspace where migrations should be fulfilled (default: "platform_tasking_db")
* <b>CassandraTimeoutSec</b>      Cassandra Timeout before retrying to connect (default: 5)
* <b>CassandraConnNumber</b>      Cassandra Number Of Connections (default: 20)
