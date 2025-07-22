package redis

import "github.com/redis/go-redis/v9"

var fetchPaymentProcessorsScript = redis.NewScript(`
local defaultData = redis.call('HMGET', 'payment-processor', 'failing', 'minResponseTime')
local fallbackData = redis.call('HMGET', 'payment-processor-fallback', 'failing', 'minResponseTime')

return {defaultData[1], defaultData[2], fallbackData[1], fallbackData[2]}
`)

var createPaymentScript = redis.NewScript(`
local paymentKey = KEYS[1]
local sortedSetKey = KEYS[2]
local timestamp = ARGV[1]
local amountCents = ARGV[2]

redis.call('ZADD', sortedSetKey, timestamp, paymentKey)
redis.call('HSET', paymentKey, 'amount', amountCents)
return 'OK'
`)

var getPaymentSummaryScript = redis.NewScript(`
local sortedSetKey = KEYS[1]
local fromTimestamp = ARGV[1]
local toTimestamp = ARGV[2]

local paymentKeys = redis.call('ZRANGEBYSCORE', sortedSetKey, fromTimestamp, toTimestamp)
local defaultRequests = 0
local defaultAmountCents = 0
local fallbackRequests = 0
local fallbackAmountCents = 0

for i = 1, #paymentKeys do
    local paymentKey = paymentKeys[i]
    local amountCents = redis.call('HGET', paymentKey, 'amount')
    
	local amount = tonumber(amountCents)
	if string.match(paymentKey, "default:") then
		defaultRequests = defaultRequests + 1
		defaultAmountCents = defaultAmountCents + amount
	elseif string.match(paymentKey, "fallback:") then
		fallbackRequests = fallbackRequests + 1
		fallbackAmountCents = fallbackAmountCents + amount
	end
end

return {defaultRequests, defaultAmountCents, fallbackRequests, fallbackAmountCents}
`)

var purgeScript = redis.NewScript(`
local defaultData = redis.call('HMGET', 'payment-processor', 'failing', 'minResponseTime')
local fallbackData = redis.call('HMGET', 'payment-processor-fallback', 'failing', 'minResponseTime')
`)
