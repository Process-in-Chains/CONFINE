package prosessMiningAlgorithms

import (
	"fmt"
	"sort"
	"strings"
)

type Set struct {
	data map[string]struct{}
}

type SetOfSets struct {
	Data map[string]*Set
}

func NewSet() *Set {
	return &Set{
		data: make(map[string]struct{}),
	}
}

func NewSetOfSets() *SetOfSets {
	return &SetOfSets{
		Data: make(map[string]*Set),
	}
}

func (s *Set) Add(element string) {
	s.data[element] = struct{}{}
}

func (s *Set) Contains(element string) bool {
	_, exists := s.data[element]
	return exists
}

func (s *SetOfSets) Add(set *Set) {
	key := set.String()
	s.Data[key] = set
}

func (s *SetOfSets) Contains(set *Set) bool {
	key := set.String()
	_, exists := s.Data[key]
	return exists
}

func (s *Set) String() string {
	elements := make([]string, 0, len(s.data))
	for element := range s.data {
		elements = append(elements, element)
	}
	sort.Strings(elements)
	return strings.Join(elements, ",")
}

func (s *SetOfSets) Print() {
	fmt.Printf("{\t")
	for k, _ := range s.Data {
		fmt.Printf("%v\t", k)
	}
	fmt.Printf("}")
}
