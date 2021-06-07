package service

import (
	"fmt"
	"log"
	"time"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/repository"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

type simpleExpiredCollector struct {
	ruleEngine RuleEngine
	ruleRepository repository.RuleRepository
	internalChannel chan interface{}
}

func NewSimpleExpiredCollector(ruleEngine RuleEngine, ruleRepository repository.RuleRepository) ExpiredCollector {
	return &simpleExpiredCollector{
		ruleEngine:      ruleEngine,
		ruleRepository:  ruleRepository,
		internalChannel: make(chan interface{}, 15),
	}
}

func (s *simpleExpiredCollector) Run() {
	sleepTime := time.Duration(1) * time.Second
	var err error
	for {
		select {
		case <-s.internalChannel:
			log.Println("scan demanded")
			sleepTime, err = s.scanCleanUp()
			fmt.Println(sleepTime.Seconds())
			if err != nil {
				utils.LogError("simpleExpiredCollector", "Run", err)
				sleepTime = time.Duration(1 * time.Second)
			}
			log.Printf("next clean up planned for %fs\n", sleepTime.Seconds())
		case <-time.After(sleepTime):

			sleepTime, err = s.scanCleanUp()
			if err != nil {
				utils.LogError("simpleExpiredCollector", "Run", err)
				sleepTime = 1 * time.Second
			}
			if sleepTime == 0 {
				sleepTime = 1 * time.Second
			}
			log.Printf("next clean up planned for %fs\n", sleepTime.Seconds())

		}
	}
}

func (s *simpleExpiredCollector) Reload() {
	s.internalChannel <- true
}

func (s *simpleExpiredCollector) scanCleanUp() (time.Duration, error){
	// Delete expired rules from DB and iptables
	rules, err := s.ruleRepository.GetExpiredRules()
	if err != nil {
		utils.LogError("simpleExpiredCollector", "scanCleanUp", err)
		return 0, err
	}
	for _, r := range rules {
		err := s.deleteRule(&r)
		if err != nil {
			utils.LogError("simpleExpiredCollector", "scanCleanUp", err)
		}
	}

	// Scan for next expired rule
	rule, err := s.ruleRepository.GetNextExpiringRuleTime()
	if err != nil {
		utils.LogError("simpleExpiredCollector", "scanCleanUp", err)
	}
	if rule.EndTime == 0 {
		log.Println("do not have active role")
		time.Sleep(1 * time.Second)
		return 0, nil
	}

	dur := time.Unix(rule.EndTime, 0).Sub(time.Now())
	for dur < 0 {
		fmt.Println(dur)
		err := s.deleteRule(rule)
		if err != nil {
			utils.LogError("simpleExpiredCollector", "scanCleanUp", err)
		}
		rule, err = s.ruleRepository.GetNextExpiringRuleTime()
		if err != nil {
			log.Println("Error in reading rule repo\n", err)
		}
	}

	return dur, nil
}

func (s *simpleExpiredCollector) deleteRule(rule *entity.Rule) error {
	err := s.ruleEngine.Delete(rule)
	if err != nil {
		utils.LogError("simpleExpiredCollector", "deleteRule", err)
		return err
	} else {
		log.Printf("rule %s->%s(%s) has been disabled\n", rule.Suspect, rule.Domain, rule.IP)
	}
	err = s.ruleRepository.DeleteById(rule.ID)
	if err != nil {
		utils.LogError("simpleExpiredCollector", "deleteRule", err)
		return err
	}
	log.Printf("rule %s->%s(%s) has been deleted\n", rule.Suspect, rule.Domain, rule.IP)
	return nil
}