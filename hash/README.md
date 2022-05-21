example

consistent hash:
```
/*
learn from stathat.com/c/consistent
*/
package main

import (
    "context"
    "fmt"
    "github.com/wandore/mytool/redis"
    "time"
)

var addr = "your redis addr"
var pwd = "your redis password"
var db = 0

func grantLease(ctx context.Context, locker *redis.Locker) {
    t := time.NewTicker(time.Second * 3)
    defer t.Stop()

    for {
        select {
        case <-ctx.Done():
            fmt.Println("stop granting lease")
            return
        case <-t.C:
            locker.Keep(time.Second * 5)
            fmt.Println("grant lease")
        }
    }
}

func main() {
    myredis, err := redis.New(addr, pwd, db)
    if err != nil {
        panic(err)
    }

    key := "lock"
    value := "test"
    ttl := time.Second * 5
    var locker *redis.Locker
    for {
        locker, err = myredis.TryLock(key, value, ttl)
        if err != nil {
            fmt.Println(err, "sleep 3s")
            time.Sleep(time.Second * 3)
        } else {
            fmt.Println(locker)
            break
        }
    }

    ctx, cancelFunc := context.WithCancel(context.Background())

    go grantLease(ctx, locker)

    fmt.Println("get lock, dispatching job")

    time.Sleep(time.Second * 30)

    cancelFunc()
    locker.Unlock()
}
```