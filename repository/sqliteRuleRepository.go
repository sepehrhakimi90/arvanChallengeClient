package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

const (
	nextExpiringRawQuery = "SELECT * FROM rules where end_time = (SELECT min(end_time) from rules)"
)

type sqliteRuleRepo struct {
	db *gorm.DB
}

func NewSqliteRuleRepository(db *gorm.DB) RuleRepository{
	repo := sqliteRuleRepo{db: db}
	db.Where("1=1").Delete(&entity.Rule{})
	return &repo
}

func (s *sqliteRuleRepo) Save(rule *entity.Rule) (*entity.Rule, error) {
	result := s.db.Create(&rule)
	if result.Error != nil {
		utils.LogError("sqlLiteRuleRepository", "Save", result.Error)
		return nil, result.Error
	}
	return rule, nil
}

func (s *sqliteRuleRepo) GetExpiredRules() ([]entity.Rule, error) {
	rules := make([]entity.Rule, 0)
	result := s.db.Where("end_time <= ?", time.Now().Unix()).Find(&rules)
	if result.Error != nil {
		utils.LogError("sqlLiteRuleRepository", "GetExpiredRules", result.Error)
		return nil, result.Error
	}
	return rules, nil
}

func (s *sqliteRuleRepo) DeleteById(id uint) error {
	result := s.db.Delete(&entity.Rule{}, id)
	if result.Error != nil {
		utils.LogError("sqlLiteRuleRepository", "DeleteById", result.Error)
		return result.Error
	}
	return nil
}

func (s *sqliteRuleRepo) GetNextExpiringRuleTime() (*entity.Rule, error) {
	rule := &entity.Rule{}
	result := s.db.Raw(nextExpiringRawQuery, time.Now().Unix()).Scan(rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return rule, nil
}
