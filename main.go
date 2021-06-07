package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/repository"
	"github.com/seperhakimi90/arvanChallengeClient/service"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

var (
	serverHost = os.Getenv("SERVER_HOST")
	serverPort = os.Getenv("SERVER_PORT")
)

func main() {
	db, err := gorm.Open(sqlite.Open(".client.db"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database\n", err)
	}
	db.AutoMigrate(&entity.Rule{})

	ruleRepo := repository.NewSqliteRuleRepository(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "192.168.122.1:6379",
		Password: "",
		DB:       0,
	})

	ruleEngine, err := service.NewIpTableRuleEngine()
	if err != nil {
		log.Fatalln(err)
	}

	err = service.StartUp(ruleRepo, ruleEngine, fmt.Sprintf("http://%s:%s/rules", serverHost, serverPort))
	if err != nil {
		utils.LogError("main", "main", err)
	}

	expireCollector := service.NewSimpleExpiredCollector(ruleEngine, ruleRepo)
	go expireCollector.Run()

	redisService := service.NewRedisSubscriber(redisClient, expireCollector)
	redisService.Subscriber(context.Background(), ruleRepo, ruleEngine)

	/*go func() {
		time.Sleep(15 * time.Second)
		for {
			fmt.Println("Checker running")
			rules, err := ruleRepo.GetExpiredRules()
			if err != nil {
				log.Println(err)
			}
			//fmt.Printf("%#v\n", rules)
			for _, r := range rules {
				err = ruleEngine.Delete(&r)
				if err != nil {
					log.Printf("error in deleting rule: %#v\n", r)
				}
				err = ruleRepo.DeleteById(r.ID)
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("%#v has been deleted\n",r)
			}
			rule, err := ruleRepo.GetNextExpiringRuleTime()
			if err != nil {
				log.Println("Error in reading rule repo\n", err)
			}
			dur := time.Unix(rule.EndTime, 0).Sub(time.Now())
			if rule.EndTime == 0 {
				fmt.Println("do not have active role")
				time.Sleep(1 * time.Second)
				continue
			}
			for dur < 0 {
				fmt.Println(dur)
				err := ruleRepo.DeleteById(rule.ID)
				if err != nil {
					log.Println(err)
				}
				rule, err = ruleRepo.GetNextExpiringRuleTime()
				if err != nil {
					log.Println("Error in reading rule repo\n", err)
				}
			}
			fmt.Printf("going to sleep for %f\n", dur.Seconds())
			time.Sleep(dur)
		}

	}()*/


	var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	wg.Add(1)
	go func(){
		defer wg.Done()
		for sig := range c {
			if sig == os.Kill || sig == os.Interrupt {
				return
			}
		}
	}()

	wg.Wait()
}
