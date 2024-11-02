package model

import (
	"fmt"
	"sort"
	"strings"
)

// Cloud represents an Openstack environment with env vars
type Cloud struct {
	Name   string
	Source string
	Env    map[string]string
}

type Clouds []Cloud

// Select a cloud by name
func (c Clouds) Select(name string) Cloud {
	for _, v := range c {
		if v.Name == name {
			return v
		}
	}
	return Cloud{}
}

func (c Cloud) String() string {
	var out strings.Builder

	keys := make([]string, 0, len(c.Env))
	var longest int
	for k := range c.Env {
		if len(k) > longest {
			longest = len(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out.WriteString("--- " + c.Name + " ---\n")

	var val string
	for _, k := range keys {
		if k == "OS_PASSWORD" {
			val = "****"
		} else {
			val = c.Env[k]
		}
		out.WriteString(fmt.Sprintf("%-*s %s\n", longest+1, k, val))
	}
	out.WriteString("(source: " + c.Source + ")")
	return out.String()
}
