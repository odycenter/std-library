package redis_test

import (
	"context"
	"github.com/odycenter/std-library/redis"
	"github.com/pkg/errors"
	"log"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redis.Init(&redis.Opt{
		IsCluster:    false,
		Addrs:        []string{"127.0.0.1:6379"},
		Password:     "",
		PoolSize:     40,
		MinIdleConns: 20,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	s, err := redis.RDB().WithCtx(ctx).Get("A")
	if err != nil {
		if err != redis.Nil {
			log.Fatal(errors.Wrap(err, "Redis get failed"))
		}
		err = nil
	}
	time.Sleep(time.Second * 2)
	log.Println("-1-:", s)
	err = redis.RDB().WithCtx(ctx).Set("A", time.Now().Format(time.RFC850))
	if err != nil {
		if err != redis.Nil {
			log.Fatal(errors.Wrap(err, "Redis set failed"))
		}
		err = nil
	}
	s, err = redis.RDB().WithCtx(ctx).Get("A")
	if err != nil {
		if err != redis.Nil {
			log.Fatal(errors.Wrap(err, "Redis get failed"))
		}
		err = nil
	}
	log.Println("-2-:", s)
	//err = redis.RDB().WithCtx(ctx).Set("A", time.Now().Format(time.RFC850), 10*time.Second)
	err = redis.RDB().WithCtx(ctx).Set("A", time.Now().Format(time.RFC850), 2)
	if err != nil {
		if err != redis.Nil {
			log.Fatal(errors.Wrap(err, "Redis set failed"))
		}
		err = nil
	}
	time.Sleep(time.Second * 2)
	s, err = redis.RDB().WithCtx(ctx).Get("A")
	if err != nil {
		if err != redis.Nil {
			log.Fatal(errors.Wrap(err, "Redis get failed"))
		}
		err = nil
	}
	log.Println("-3-:", s)
}
