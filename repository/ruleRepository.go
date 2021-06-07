package repository

import (
	"github.com/seperhakimi90/arvanChallengeClient/entity"
)

type RuleRepository interface {
	Save(*entity.Rule) (*entity.Rule, error)
	GetExpiredRules() ([]entity.Rule, error)
	DeleteById(uint) error
	GetNextExpiringRuleTime() (*entity.Rule, error)
}
