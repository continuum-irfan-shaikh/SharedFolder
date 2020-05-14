local path = "tasking/v1/partners/1d4400c0/task-definitions"

local open = io.open

local function read_file(path)
    local file = open(path, "rb")
    if not file then return nil end
    local content = file:read "*a"
    file:close()
    return content
end

local file = read_file("./templates/template.json")
request = function ()
    wrk.body = file
    return wrk.format(wrk.method, wrk.path .. path)
end
wrk.method = "POST"