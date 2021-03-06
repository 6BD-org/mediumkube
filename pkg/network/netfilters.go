package network

import (
	"github.com/coreos/go-iptables/iptables"
	"k8s.io/klog/v2"
)

type IPMode int

const (
	// IPModAPP appends rule to filter chain
	IPModAPP IPMode = 1
	// IPModPREP prepend
	IPModPREP IPMode = 2

	table = "filter"
)

// InsertRuleIfNotExists Insert a rule into iptable if not exists
func InsertRuleIfNotExists(chain string, mode IPMode, rules ...string) {

	iptable, err := iptables.New()
	if err != nil {
		klog.Error(err)
		return
	}

	exists, err := iptable.Exists(
		table,
		chain,
		rules...,
	)
	if err != nil {
		klog.Error(err)
		return
	}

	if !exists {
		klog.Info("Appending: ", rules)
		if mode == IPModAPP {
			iptable.Append(table, chain, rules...)

		} else if mode == IPModPREP {
			iptable.Insert(table, chain, 1, rules...)

		}
	}
}

// DeleteIfExists deletes a rule if it exists
func DeleteIfExists(chain string, rules []string) {
	iptable, err := iptables.New()
	if err != nil {
		klog.Error(err)
		return
	}
	exists, err := iptable.Exists(
		table,
		chain,
		rules...,
	)
	if err != nil {
		klog.Error(err)
		return
	}
	if exists {
		klog.Info("Deleting: ", rules)
		iptable.Delete(table, chain, rules...)
	}
}
