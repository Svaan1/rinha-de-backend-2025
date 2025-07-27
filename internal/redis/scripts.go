package redis

import "github.com/redis/go-redis/v9"

var createPaymentScript = redis.NewScript(`
local paymentKey = KEYS[1]
local sortedSetKey = KEYS[2]
local timestamp = ARGV[1]
local amountCents = ARGV[2]

redis.call('ZADD', sortedSetKey, timestamp, paymentKey .. ':' .. amountCents)
return 'OK'
`)

var getPaymentSummaryScript = redis.NewScript(`
local sortedSetKey = KEYS[1]
local fromTimestamp = ARGV[1]
local toTimestamp = ARGV[2]

local paymentEntries = redis.call('ZRANGEBYSCORE', sortedSetKey, fromTimestamp, toTimestamp)

local defaultRequests = 0
local defaultAmountCents = 0
local fallbackRequests = 0
local fallbackAmountCents = 0

for i = 1, #paymentEntries do
    local entry = paymentEntries[i]

	local colonPos = string.find(entry, ':')
    local secondColonPos = string.find(entry, ':', colonPos + 1)
    
    if secondColonPos then
        local processorType = string.sub(entry, 1, colonPos - 1)
        local amountStr = string.sub(entry, secondColonPos + 1)
        local amount = tonumber(amountStr)
        
        if processorType == "default" then
            defaultRequests = defaultRequests + 1
            defaultAmountCents = defaultAmountCents + amount
        elseif processorType == "fallback" then
            fallbackRequests = fallbackRequests + 1
            fallbackAmountCents = fallbackAmountCents + amount
        end
    end
end

return {defaultRequests, defaultAmountCents, fallbackRequests, fallbackAmountCents}
`)
