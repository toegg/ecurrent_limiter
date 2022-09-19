### Go+Redis实现的并发安全限流器

  
本限流器实现了**计数器限流**和**滑动窗口限流**，提供了非并发安全和并发安全的实现，方便两者的对比。  


#### 计数器限流:

设计思路：核心是通过Redis的Incr和Expire设置过期时间
```
三个参数, key:限流的Api存储key值，count：限流上限个数，ttl:限流单位时间
1. 先Get key，判断有没有超过限流上限count
2. 没超过上限，可以直接放行，执行Incr。Incr为1的话则说明是限流单位时间区间内第一个请求，需要设置ttl过期时间
3. 超过上限，需要判断ttl是否没设置(因为存在第2步的Incr成功了，但是Expire失败了)
4. 设置了ttl的，说明在限定时间内超过上限，限流不放行
5. 未设置ttl的，用Set+px参数原子性操作设置为1，成功则放行，失败则限流
```  

#### 滑动窗口限流:  

设计思路：核心是利用list队列左进右出，个数占位推进代替时间推进（空间代替时间推进的转换）
```
三个参数, key:限流的Api存储key值，count：限流上限个数，windowTime:滑动窗口时间
1. 判断list队列长度是否超过上限count
2. 没超过上限，直接放行，把当前时间戳(秒)放进去队列
3. 超过上限，判断队列最右边占位的时间戳和当前时间戳的差值是否大于windowTime
4. 小于窗口时间，说明在窗口时间内达到上限，限流不放行
5. 大于窗口时间，说明已推进到新窗口，移除最右边的并且放入当前时间戳到最左边，放行
```  

#### 使用：
```go

import "github.com/toegg/ecurrent_limiter/limiter"

//创建限流器
l, err := limiter.NewCache()
//计数器并发安全限流
res := l.SyncCountLimit(Api存储Redis的Key值, 限流个数, 限流单位时间)
if res {
    //逻辑处理
}
//计数器限流
res = l.CountLimit(Api存储Redis的Key值, 限流个数, 限流单位时间)

//滑动窗口并发安全限流
res = l.SyncWindowLimit(Api存储Redis的Key值, 限流个数, 滑动窗口时间)

//滑动窗口限流
res = l.WindowLimit(Api存储Redis的Key值, 限流个数, 滑动窗口时间)
```

**测试结果**：每个类型并发10个请求，设定5秒内限流放行2个，跑出来结果，只有SyncCountLimit和SyncWindowLimit是符合预期。
具体测试代码看**ecurrent_limiter.go**


