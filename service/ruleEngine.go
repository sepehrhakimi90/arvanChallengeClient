package service

import (
	"github.com/seperhakimi90/arvanChallengeClient/entity"
)

type RuleEngine interface {
	Add(*entity.Rule) error
	Delete(*entity.Rule) error
	Reset() error
}
