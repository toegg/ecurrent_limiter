package main

import (
	"ecurrent_limiter/limiter"
	"log"
	"sync"
)

func main(){
  	l, err := limiter.NewCache()
  	if err != nil{
  		log.Fatalf("RedisInit Err:%v", err)
	}

	//计数器限流普通
  	i := 10
  	var wait sync.WaitGroup
  	for i > 0{
  		wait.Add(1)
  		go func(index int){
  			res :=l.CountLimit("CountLimit", 2, 5)
  			log.Printf("CountLimit-Index:%v,RequestAllow:%v", index, res)
  			wait.Done()
		}(i)
  		i--
	}
	wait.Wait()

  	//计数器限流并发安全
  	i = 10
	for i > 0{
		wait.Add(1)
		go func(index int){
			res :=l.SyncCountLimit("SyncCountLimit", 2, 5)
			log.Printf("SyncCountLimit-Index:%v,RequestAllow:%v", index, res)
			wait.Done()
		}(i)
		i--
	}
	wait.Wait()

	//滑动窗口限流普通
	i = 10
	for i > 0{
		wait.Add(1)
		go func(index int){
			res :=l.WindowLimit("WindowLimit", 2, 5)
			log.Printf("WindowLimit-Index:%v,RequestAllow:%v", index, res)
			wait.Done()
		}(i)
		i--
	}
	wait.Wait()

	//滑动窗口限流并发安全
	i = 10
	for i > 0{
		wait.Add(1)
		go func(index int){
			res :=l.SyncWindowLimit("SyncWindowLimit", 2, 5)
			log.Printf("SyncWindowLimit-Index:%v,RequestAllow:%v", index, res)
			wait.Done()
		}(i)
		i--
	}
	wait.Wait()

}