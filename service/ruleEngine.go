package service

import (
	"github.com/seperhakimi90/arvanChallengeClient/entity"
)

type RuleEngine interface {
	add(*entity.Rule) error
	delete(*entity.Rule) error
	reset() error
}
