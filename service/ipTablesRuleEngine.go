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

func NewIpTableRuleEngine() (RuleEngine, error) {
	ipt, err := iptables.New()
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "NewIpTableRuleEngine", err)
		return nil, err
	}
	iptRuleEngine := &ipTableRuleEngine{ipt: ipt}
	iptRuleEngine.Reset()
	return iptRuleEngine, nil
}

func (i *ipTableRuleEngine) Add(rule *entity.Rule) error {
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

	err = i.ipt.Append(tableName, chainName, "-s", rule.Suspect, "-d", rule.IP, "-j", dropTargetName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "add", err)
		return err
	}
	return nil
}

func (i *ipTableRuleEngine) Delete(rule *entity.Rule) error {
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

	exist, err = i.ipt.Exists(tableName, chainName, "-s", rule.Suspect,"-d", rule.IP ,"-j", dropTargetName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}
	if !exist {
		err = errors.New(fmt.Sprintf("rule for restrict  %s on %s(%s) does not exist\n", rule.Domain, rule.Suspect, rule.IP))
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}

	err = i.ipt.Delete(tableName, chainName, "-s", rule.Suspect, "-d", rule.IP ,"-j", dropTargetName)
	if err != nil {
		utils.LogError("ipTablesRuleEngine", "delete", err)
		return err
	}

	return nil
}

func (i *ipTableRuleEngine) Reset() error {
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
	i.ipt.Insert(tableName, "INPUT", 1,"-j", chainName)
	return nil
}
