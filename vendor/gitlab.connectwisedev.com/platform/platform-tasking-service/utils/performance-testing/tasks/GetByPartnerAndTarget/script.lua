-- wrk settings
partner_id = "1d4400c0"
path =  "tasking/v1/partners/" .. partner_id .. "/tasks/targets/%s"

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

-- Selecting managed endpoints
local tasks = assert(peer:execute(string.format("SELECT target FROM tasks WHERE partner_id='%s'", partner_id)))

-- Performing request
counter = 1
request = function()
    if counter > #tasks then
        counter = 1
    end
    counter = counter + 1
    url = wrk.path .. string.format(path, tasks[counter].target)
    return wrk.format(wrk.method, url)
end
