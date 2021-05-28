package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/repository"
)


type redisSubscriber struct {
	rdb *redis.Client
}

func NewRedisSubscriber(rdb *redis.Client) Subscriber{
	return &redisSubscriber{rdb: rdb}
}

func (r *redisSubscriber) Subscriber(ctx context.Context, ruleRepository repository.RuleRepository, ruleEngine RuleEngine) error {
	pubsub := r.rdb.Subscribe(ctx, "ruleChannel")
	ch := pubsub.Channel()
	go r.run(ctx, ch, ruleRepository, ruleEngine)
	return nil
}

func (r *redisSubscriber) run(ctx context.Context, ch <- chan *redis.Message, ruleRepository repository.RuleRepository, ruleEngine RuleEngine) {
	for {
		select {
		case msg := <-ch:
			rule := entity.Rule{}
			err := json.Unmarshal([]byte(msg.Payload), &rule)
			if err != nil {
				fmt.Println(err)
			}
			err = ruleEngine.add(&rule)
			if err != nil {
				break
			}
			ruleRepository.Save(&rule)
		case <- ctx.Done():
			return
		}
	}
}
