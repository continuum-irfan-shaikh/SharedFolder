-- wrk settings
partner_id = "1d4400c0"
general_partner = "00000000-0000-0000-0000-000000000000"
path =  "tasking/v1/partners/" .. partner_id .. "/tasks-definition-templates/script/%s"

host_cassandra_local = "127.0.0.1"
host_cassandra_remote = "172.28.48.6"

-- Setting cassandra connection
local cassandra = require "cassandra"
local peer = assert(cassandra.new {
    host = host_cassandra_local,
    port = 9042,
    keyspace = "platform_scripting_db"
})
peer:settimeout(1000)
assert(peer:connect())

-- Getting scripts
local scripts = assert(peer:execute(string.format("SELECT * FROM scripts WHERE partner_id in ('%s', '%s')", general_partner, partner_id)))
peer:close()

-- Performing request
counter = 1
request = function()
    if counter > #scripts then
        counter = 1
    end
    script = scripts[counter]
    counter = counter + 1
    url = wrk.path .. string.format(path, script.id)
    return wrk.format(wrk.method, url)
end
