package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/repository"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

func StartUp(ruleRepository repository.RuleRepository, ruleEngine RuleEngine, rulesUrl string) error{
	log.Println("Retrieving active rules")
	rules := make([]entity.Rule, 0, 1)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(rulesUrl)
	if err != nil {
		utils.LogError("startUp", "StartUp", err)
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&rules)
	if err != nil {
		utils.LogError("startUp", "StartUp", err)
		return err
	}
	log.Printf("loading %d active rules", len(rules))
	for _, rule := range rules {
		IPs, err := utils.GetDomainIPv4(rule.Domain)
		if err != nil {
			log.Println("Error in address Lookup==>", err)
			return err
		} else {
			for _, ip := range IPs {
				//rule.ID = 0
				rule.IP = ip.String()
				err = ruleEngine.Add(&rule)
				if err != nil {
					utils.LogError("redisSubscribe", "run", err)
					continue
				}
				ruleRepository.Save(&rule)
			}
		}
	}
	return nil
}
