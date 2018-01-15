package util

import "sort"

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func Reverse(data sort.Interface) sort.Interface {
	return &reverse{data}
}
