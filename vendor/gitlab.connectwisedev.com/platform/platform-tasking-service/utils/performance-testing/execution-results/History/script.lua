-- wrk settings
partner_id = "1d4400c0"
path =  "tasking/v1/partners/" .. partner_id .. "/task-execution-results/tasks/%s/managed-endpoints/%s/history"

host_cassandra_local = "127.0.0.1"
host_cassandra_remote = "172.28.48.6"

-- Setting cassandra connection
local cassandra = require "cassandra"
local peer = assert(cassandra.new {
    host = host_cassandra_local,
    port = 9042,
    keyspace = "platform_tasking_db"
})
peer:settimeout(1000)
assert(peer:connect())

-- Getting tasks
local IDs = assert(peer:execute(string.format("SELECT id FROM tasks WHERE partner_id='%s'", partner_id)))
local TARGETS = assert(peer:execute(string.format("SELECT target FROM tasks WHERE partner_id='%s'", partner_id)))

counter = 1

request = function()
	if counter > #IDs then
        counter = 1
    end
    counter = counter + 1
    url = wrk.path .. string.format(path, IDs[counter].id, TARGETS[counter].target)
    return wrk.format(wrk.method, url)
end
