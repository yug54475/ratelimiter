-- KEYS[1] = rate limit key (e.g. "ratelimit:clientID")
-- ARGV[1] = limit
-- ARGV[2] = window in seconds

local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
end

if current > tonumber(ARGV[1]) then
    return 0
else
    return 1
end