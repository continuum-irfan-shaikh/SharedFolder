path =  "tasking/v1/partners/" .. partner_id .. "/tasks-definitions/%s"

host_cassandra_local = "10.128.233.149"
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

-- Getting IDs
local IDs = assert(peer:execute("SELECT id FROM task_definitions"))
peer:close()

-- Performing request
counter = 1
request = function()
	if counter > #IDs then
        counter = 1
    end
    id = IDs[counter]
    url = wrk.path .. string.format(path, id.id)
    counter = counter + 1
    return wrk.format(wrk.method, url)
end