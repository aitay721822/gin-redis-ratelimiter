package ginredisratelimiter

type TokenBucketLuaRequest struct {
	valueKey		string
	timestampKey	string
	limit			int64
	interval		int64
	batchSize		int64
}

const TokenBucketLuaScript = `
-- Request value
local valueKey     = KEYS[1]
local timestampKey = KEYS[2]
local limit        = tonumber(ARGV[1])
local interval     = tonumber(ARGV[2]) -- milliseconds
local batchSize    = math.max(tonumber(ARGV[3]), 0)

-- Response value
local rejected     = false
local remainToken  = 0

redis.replicate_commands()

local time = redis.call('TIME')
local currentTime = math.floor(time[1] * 1000 + time[2] / 1000)
local modified = false
local lastRemainToken = redis.call('GET', valueKey)
local lastUpdateTime = false

if lastRemainToken == false then
   lastRemainToken = 0
   lastUpdateTime = currentTime - interval
else 
   lastUpdateTime = redis.call('GET', timestampKey)
   if lastUpdateTime == false then
      modified = true
      lastUpdateTime = currentTime - ((lastRemainToken / limit) * interval)
   end
end

local feedbackToken = math.max((currentTime - lastUpdateTime) / interval * limit, 0)
local token = math.min(lastRemainToken + feedbackToken, limit)
remainToken = token - batchSize

if remainToken < 0 then
   rejected = true
   remainToken = token
end

if rejected == false then
   redis.call('PSETEX', valueKey, interval, remainToken)
   if feedbackToken > 0 or modified then
      redis.call('PSETEX', timestampKey, interval, currentTime)
   else 
      redis.call('PEXPIRE', timestampKey, interval)
   end
end

return { rejected, remainToken }
`