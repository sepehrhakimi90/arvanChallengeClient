package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/repository"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

type redisSubscriber struct {
	rdb       *redis.Client
	collector ExpiredCollector
}

func NewRedisSubscriber(rdb *redis.Client,expiredCollector ExpiredCollector) Subscriber {
	return &redisSubscriber{
		rdb: rdb,
		collector: expiredCollector,
	}
}

func (r *redisSubscriber) Subscriber(ctx context.Context, ruleRepository repository.RuleRepository, ruleEngine RuleEngine) error {
	pubsub := r.rdb.Subscribe(ctx, "ruleChannel")
	ch := pubsub.Channel()
	go r.run(ctx, ch, ruleRepository, ruleEngine)
	return nil
}

// Todo exponential backoff for client reconnection to redis
func (r *redisSubscriber) run(ctx context.Context, ch <-chan *redis.Message, ruleRepository repository.RuleRepository, ruleEngine RuleEngine) {
	for {
		select {
		case msg := <-ch:
			rule := entity.Rule{}
			err := json.Unmarshal([]byte(msg.Payload), &rule)
			if err != nil {
				utils.LogError("redisSubscribe", "run", err)
				continue
			}

			IPs, err := utils.GetDomainIPv4(rule.Domain)
			if err != nil {
				log.Println("Error in address Lookup==>", err)
			} else {
				for _, ip := range IPs {
					rule.ID = 0
					rule.IP = ip.String()
					err = ruleEngine.Add(&rule)
					if err != nil {
						utils.LogError("redisSubscribe", "run", err)
						continue
					}
					ruleRepository.Save(&rule)
				}
				r.collector.Reload()
			}
		case <-ctx.Done():
			return
		}
	}
}
