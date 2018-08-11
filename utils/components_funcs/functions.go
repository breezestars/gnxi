package components_funcs

import (
	"github.com/derekparker/trie"
	"strings"
	log "github.com/golang/glog"
	"fmt"
)

type dataSetFunc func([]string, []string, string, bool) (error)
type dataDelFunc func([]string, []string, string, bool) (error)

type path2DataFunc struct {
	path    []string
	setFunc dataSetFunc
	delFunc dataDelFunc
}

var (
	setTrie *trie.Trie
	delTrie *trie.Trie
	setTbl  = []path2DataFunc{
		{
			path:    []string{"network-instances", "network-instance", "vlans", "vlan", "config", "vlan-id"},
			setFunc: dataSetFunc(SetVlan),
		},
		{
			path:    []string{"network-instances", "network-instance", "vlans", "vlan", "members", "member", "state", "interface", "config", "name"},
			setFunc: dataSetFunc(SetVlanMember),
		},
		{
			path:    []string{"interfaces", "interface", "config", "enabled"},
			setFunc: dataSetFunc(SetInterfaceConfigEnabled),
		},
		{
			path: []string{"interfaces","interface","aggregation", "config", "lag-type"},
			setFunc:dataSetFunc(SetInterfacePortchannel),
		},
	}

	delTbl = []path2DataFunc{
		{
			path:    []string{"network-instances", "network-instance", "vlans", "vlan", "config", "vlan-id"},
			delFunc: dataDelFunc(DelVlan),
		},
		{
			path:    []string{"network-instances", "network-instance", "vlans", "vlan", "members", "member", "state", "interface", "config", "name"},
			delFunc: dataDelFunc(DelVlanMember),
		},
		{
			path: []string{"interfaces","interface","aggregation", "config", "lag-type"},
			setFunc:dataSetFunc(DelInterfacePortchannel),
		},
	}
)

func delTriePopulate(t *trie.Trie) {
	for _, p2DF := range delTbl {
		n := t.Add(strings.Join(p2DF.path, ""), p2DF.delFunc)
		if n.Meta().(dataDelFunc) == nil {
			log.V(1).Infof("Failed to add trie node for %v with %v", p2DF.path, p2DF.delFunc)
		} else {
			log.V(2).Infof("Add trie node for %v with %v", p2DF.path, p2DF.delFunc)
		}
	}
}

func setTriePopulate(t *trie.Trie) {
	for _, p2DF := range setTbl {
		n := t.Add(strings.Join(p2DF.path, ""), p2DF.setFunc)
		if n.Meta().(dataSetFunc) == nil {
			log.V(1).Infof("Failed to add trie node for %v with %v", p2DF.path, p2DF.setFunc)
		} else {
			log.V(2).Infof("Add trie node for %v with %v", p2DF.path, p2DF.setFunc)
		}
	}
}

func init() {
	delTrie = trie.New()
	setTrie = trie.New()
	delTriePopulate(delTrie)
	setTriePopulate(setTrie)
}

func LookupSetFunc(path [] string) (dataSetFunc, error) {
	node, ok := setTrie.Find(strings.Join(path, ""))
	if ok {
		setter := node.Meta().(dataSetFunc)
		return setter, nil
	}
	return nil, fmt.Errorf("%v not found in setTrie tree", path)
}

func LookupDelFunc(path [] string) (dataDelFunc, error) {
	node, ok := delTrie.Find(strings.Join(path, ""))
	if ok {
		deleter := node.Meta().(dataDelFunc)
		return deleter, nil
	}
	return nil, fmt.Errorf("%v not found in delTrie tree", path)
}
