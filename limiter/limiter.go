package limiter

import (
	"github.com/toegg/ecurrent_limiter/utils"
	"github.com/go-redis/redis"
	"log"
	"time"
)

const(
	//本地连接
	redisHost = "127.0.0.1"
	redisPort = "6379"
	redisPassWord = ""
)

type Limiter struct {
	conn *redis.Client
}

//创建Redis连接
func NewCache() (*Limiter, error){
	addr := redisHost + ":" + redisPort
	conn := redis.NewClient(&redis.Options{
		Addr:addr,
		Password: redisPassWord,
	})
	_, err := conn.Ping().Result()
	if err != nil{
		return nil, err
	}

	limiter := new(Limiter)
	limiter.conn = conn
	return limiter, nil
}

// ----------------------------------------计数器限流
//核心：通过incr+设置过期时间
//1. 先Get key，判断有没有超过上限count
//2. 没超过上限，可以直接放行，Incr为1的话则说明是时间区间内第一个请求，需要设置ttl过期时间
//3. 超过上限，需要判断ttl是否没设置(因为存在第2步的Incr成功了，但是Expire失败了)
//4. 设置了ttl的，说明在限定时间内超过上限，限流不放行
//5. 未设置ttl的，用Set+px参数原子性操作设置为1，成功则放行，失败则限流

//不保证流程原子性，存在并发竞争问题
func (l *Limiter) CountLimit(key string, count uint, ttl int64) bool{

	reqCounts, _ := l.conn.Get(key).Int64()
	if uint(reqCounts) < count {
		reqCounts, _ =l.conn.Incr(key).Result()
		if reqCounts == 1 {
			l.conn.Expire(key, time.Duration(ttl) * time.Second)
		}
		return true
	}

	if l.conn.TTL(key).Val() <= 0 {
		err := l.conn.Set(key, 1, time.Duration(ttl) * time.Second).Err()
		if err!=nil{
			log.Println("CountLimit Set Expire Err:", err)
			return false
		}

		return true
	}
	return false
}

//Lua脚本保证流程原子性，并发安全
func (l *Limiter) SyncCountLimit(key string, count uint, ttl int64) bool{

	var luaScript string
	luaScript =
		" local key = KEYS[1] " +
			" local ttl = ARGV[2] " +
			" local count = ARGV[1] " +
			" local reqCounts = redis.call('get', key) " +
			" if (not reqCounts or tonumber(reqCounts) < tonumber(count)) then " +
			"	 reqCounts = redis.call('incr', key) " +
			"	 if tonumber(reqCounts) == 1 then " +
			"		 redis.call('expire', key, tonumber(ttl)) " +
			"	 end " +
			"	 return 1 " +
			" end " +
			" if tonumber(redis.call('ttl', key)) <= 0 then " +
			"	 local res = redis.call('set', key, 1, 'ex', tonumber(ttl)) " +
			"	 redis.log(redis.LOG_NOTICE, key..\" not set expire\")	" +
			"	 if res.ok ~= \"OK\" then " +
			"	 	 redis.log(redis.LOG_NOTICE, key..\" set again err\") 	" +
			"		 return 2 " +
			"	 end " +
			"	 return 1 " +
			" end " +
			" return 2 "

	result, err := l.conn.Eval(luaScript, []string{key}, count, ttl).Result()

	if err !=nil{
		log.Println("SyncCountLimit error:", err)
		return false
	}

	if utils.GetInt(result) != 1 {
		return false
	}

	return true
}


// ----------------------------------------滑动窗口限流
//核心：利用list队列左进右出，个数占位推进代替时间推进
//1. 判断list队列长度是否超过上限count
//2. 没超过上限，直接放行，把当前时间放进去队列
//3. 超过上限，判断队列最右边占位的时间和当前时间差是否大于窗口时间
//4. 小于窗口时间，说明在窗口时间内达到上限，限流不放行
//5. 大于窗口时间，说明已推进到新窗口，移除最右边的，放入当前时间，放行

//不保证流程原子性，存在并发竞争问题
func (l *Limiter) WindowLimit(key string, count uint, windowTime int64) bool{

	time := time.Now().Unix()
	len := l.conn.LLen(key).Val()
	if uint(len) < count{
		l.conn.LPush(key, time)
		return true
	}

	earlyTime,_ := l.conn.LIndex(key, len - 1).Int64()
	if time - earlyTime < windowTime{
		return false
	}

	l.conn.RPop(key)
	l.conn.LPush(key, time)

	return true
}

//Lua脚本保证流程原子性，并发安全
func (l *Limiter) SyncWindowLimit(key string, count uint, windowTime int64) bool{

	time := time.Now().Unix()

	var luaScript string
	luaScript =
		"local key = KEYS[1] " +
			"local time = ARGV[3] " +
			"local windowTime = ARGV[2] " +
			"local count = ARGV[1] " +
			"local len = redis.call('llen', key) " +
			"if tonumber(len) < tonumber(count) then " +
			"   redis.call('lpush', key, time) " +
			"	return 1 " +
			"end " +
			"local earlyTime = redis.call('lindex', key, tonumber(len) - 1) " +
			"if tonumber(time) - tonumber(earlyTime) < tonumber(windowTime) then " +
			"	return 2 " +
			"end " +
			"redis.call('rpop', key) " +
			"redis.call('lpush', key, time) " +
			"return 1 "

	result, err := l.conn.Eval(luaScript, []string{key}, count, windowTime, time).Result()
	if err !=nil{
		log.Println("SyncWindowLimit error:", err)
		return false
	}

	if utils.GetInt(result) != 1 {
		return false
	}

	return true
}
