package service

import (
	"context"

	"github.com/seperhakimi90/arvanChallengeClient/repository"
)

type Subscriber interface {
	Subscriber(context.Context ,repository.RuleRepository, RuleEngine) error
}
