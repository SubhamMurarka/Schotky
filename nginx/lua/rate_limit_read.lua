-- Import the Redis and MurmurHash3 libraries
local redis = require "resty.redis"
local murmur = require("murmurhash3")
if not murmur then
    ngx.log(ngx.ERR, "Failed to load murmurhash3 module")
else
    ngx.log(ngx.ERR, "Murmurhash3 module loaded successfully")
end

-- Attempt to log the hash32 function
-- ngx.log(ngx.ERR, "hash32 function: ", tostring(murmur.hash32))

-- Configuration
local redis_servers = {
    { host = "redis1", port = 6379 },
    { host = "redis2", port = 6379 },
    { host = "redis3", port = 6379 },
}

local rate_limit = 5  -- Maximum number of requests allowed per IP

-- Get client IP and use it as the key
local client_ip = ngx.var.http_x_forwarded_for or ngx.var.remote_addr
local key = client_ip  -- Define key as client IP

ngx.log(ngx.ERR, "Key: ", key)

-- Compute hash and determine Redis instance
local hash = murmur.hash32(key)  -- Using a fixed seed of 42 for hash computation
local redis_index = (hash % #redis_servers) + 1
local selected_redis = redis_servers[redis_index]

-- Connect to Redis
local red = redis:new()
red:set_timeout(1000)  -- 1 second

local ok, err = red:connect(selected_redis.host, selected_redis.port)
if not ok then
    ngx.log(ngx.ERR, "Failed to connect to Redis: ", err)
    return ngx.exit(500)
end

-- Get current count from Redis
local res, err = red:get(key)
if not res then
    ngx.log(ngx.ERR, "Failed to get key from Redis: ", err)
    return ngx.exit(500)
end

if res == ngx.null then
    ngx.status = 200
    return
else
    local count = tonumber(res)
    if count and count < rate_limit then
        ngx.status = 200
        return
    else
        -- Rate limit exceeded, block the request
        ngx.status = 429  -- 429 Too Many Requests
        ngx.say("Rate limit exceeded. Try again later.")
        return ngx.exit(429)
    end
end

-- Close the Redis connection
red:close()