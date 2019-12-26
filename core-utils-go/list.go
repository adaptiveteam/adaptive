package core_utils_go

import (
	"math/rand"
	"time"
)

func InAAndB(a, b []string) []string {
	m := make(map[string]uint8)
	for _, k := range a {
		m[k] |= 1 << 0
	}
	for _, k := range b {
		m[k] |= 1 << 1
	}
	var inAAndB []string
	for k, v := range m {
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		switch {
		case a && b:
			inAAndB = append(inAAndB, k)
		}
	}
	return inAAndB
}

func InAButNotB(a, b []string) []string {
	m := make(map[string]uint8)
	for _, k := range a {
		m[k] |= 1 << 0
	}
	for _, k := range b {
		m[k] |= 1 << 1
	}
	var inAButNotB []string
	for k, v := range m {
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		switch {
		case a && !b:
			inAButNotB = append(inAButNotB, k)
		}
	}
	return inAButNotB
}

func InBButNotA(a, b []string) []string {
	m := make(map[string]uint8)
	for _, k := range a {
		m[k] |= 1 << 0
	}
	for _, k := range b {
		m[k] |= 1 << 1
	}
	var inBButNotA []string
	for k, v := range m {
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		switch {
		case !a && b:
			inBButNotA = append(inBButNotA, k)
		}
	}
	return inBButNotA
}

func RandomString(options []string) (rv string) {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	rv = options[rand.Intn(len(options))]
	return rv
}

// Distinct returns unique elements of the provided string slice
func Distinct(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
