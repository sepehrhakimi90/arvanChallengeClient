package service

import (
	"errors"
	"fmt"

	"github.com/coreos/go-iptables/iptables"

	"github.com/seperhakimi90/arvanChallengeClient/entity"
	"github.com/seperhakimi90/arvanChallengeClient/utils"
)

var (
	tableName      = "filter"
	chainName      = "arvanChain"
	dropTargetName = "DROP"
)

type ipTableRuleEngine struct {
	ipt *iptables.IPTables
}

func (i *ipTableRuleEngine) add(rule *entity.Rule) error {
	exist, err := i.ipt.ChainExists(tableName, chainName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "add", err)
		return err
	}
	if !exist {
		err = i.createChain(chainName)
		if err != nil {
			utils.LogError("ipTablesRuleEngine", "add", err)
			return err
		}
	}

	err = i.ipt.Append(tableName, chainName, "-s", rule.Suspect, "-j", dropTargetName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "add", err)
		return err
	}
	return nil
}

func (i *ipTableRuleEngine) delete(rule *entity.Rule) error {
	exist, err := i.ipt.ChainExists(tableName, chainName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}
	if !exist {
		err = errors.New(fmt.Sprintf("chain %s does not exist\n"))
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}

	exist, err = i.ipt.Exists(tableName, chainName, "-s", rule.Suspect, "-j", dropTargetName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}
	if !exist {
		err = errors.New(fmt.Sprintf("rule for restrict  %s on %s does not exist\n", rule.Suspect, rule.Domain))
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}

	err = i.ipt.Delete(tableName, chainName, "-s", rule.Suspect, )
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}

	return nil
}

func (i *ipTableRuleEngine) reset() error {
	exist, err := i.ipt.ChainExists(tableName, chainName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "reset", err)
		return err
	}
	if !exist {
		return nil
	}
	err = i.ipt.ClearChain(tableName, chainName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "reset", err)
		return err
	}
	return nil
}

func (i *ipTableRuleEngine) createChain(chainName string) error {
	err := i.ipt.NewChain(tableName, chainName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "createChain", err)
		return err
	}
	return nil
}
