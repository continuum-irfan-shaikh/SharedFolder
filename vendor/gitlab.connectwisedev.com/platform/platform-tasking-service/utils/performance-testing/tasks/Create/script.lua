-- don't change this value, change the variable `regularity` in wrk.sh instead
local do_not_change_me = "regularity"
--
local path = "tasking/v1/partners/1d4400c0/tasks/managed-endpoints"
--
local open = io.open
local socket = require'socket'
local function read_file(path)
    local file = open(path, "rb")
    if not file then return nil end
    local content = file:read "*a"
    file:close()
    return content
end

local file = read_file("./templates/task_" .. do_not_change_me .. "_template.json")
math.randomseed(socket.gettime()*1000)
local random = math.random

local function uuid()
    local template ='xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[x]', function (c)
        local v = (c == 'x') and random(0, 0xf) or random(8, 0xb)
        return string.format('%x', v)
    end)
end

request = function ()
    wrk.body = string.format(file, uuid())
    return wrk.format(wrk.method, wrk.path .. path)
end
wrk.method = "POST"
wrk.headers["uid"] = "admin"
