package openapi

import (
	"fmt"
	"strconv"
	"strings"
)

type APIVersion struct {
	Major    int
	Level    string
	Revision int
}

var levelOrders = []string{"alpha", "beta", "stable"}

func orderIndexOfLevel(level string) int {
	for i, l := range levelOrders {
		if l == level {
			return i
		}
	}
	return -1
}

func (v APIVersion) LessThan(another APIVersion) bool {
	if v.Major < another.Major {
		return true
	}
	if v.Major > another.Major {
		return false
	}
	if v.Level == another.Level {
		return v.Revision < another.Revision
	}
	return orderIndexOfLevel(v.Level) < orderIndexOfLevel(another.Level)
}

func ParseAPIVersion(version string) (*APIVersion, error) {
	result := &APIVersion{
		Major:    0,
		Level:    "stable",
		Revision: 0,
	}
	if !strings.HasPrefix(version, "v") {
		return nil, fmt.Errorf("version %s does not start with v", version)
	}
	if strings.Contains(version, "alpha") {
		result.Level = "alpha"
	} else if strings.Contains(version, "beta") {
		result.Level = "beta"
	} else {
		result.Level = "stable"
	}
	if result.Level == "stable" {
		major, err := strconv.Atoi(strings.Trim(version, "v"))
		if err != nil {
			return nil, err
		}
		result.Major = major
	} else {
		major, err := strconv.Atoi(strings.Trim(strings.Split(version, result.Level)[0], "v"))
		if err != nil {
			return nil, err
		}
		result.Major = major
		revision, err := strconv.Atoi(strings.Trim(strings.Split(version, result.Level)[1], "v"))
		if err != nil {
			return nil, err
		}
		result.Revision = revision
	}
	return result, nil
}
