package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

type sqliteRuleRepo struct {
	db *gorm.DB
}

func NewSqliteRuleRepository(db *gorm.DB) RuleRepository{
	repo := sqliteRuleRepo{db: db}
	return &repo
}

func (s *sqliteRuleRepo) Save(rule *entity.Rule) (*entity.Rule, error) {
	rule.EndTime = getEndTime(rule.StartTime, rule.TTL)
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

func (s *sqliteRuleRepo) DeleteById(id int) error {
	result := s.db.Delete(&entity.Rule{}, id)
	if result.Error != nil {
		utils.LogError("sqlLiteRuleRepository", "DeleteById", result.Error)
		return result.Error
	}
	return nil
}

func getEndTime(startTime time.Time, ttl int) int64{
	return startTime.Add(time.Duration(ttl) * time.Second).Unix()
}
