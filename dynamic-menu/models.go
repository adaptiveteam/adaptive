package dynamic_menu

import (
	"github.com/adaptiveteam/adaptive/checks"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

type VisibilityCheck struct {
	Name     string
	Check    checks.CheckFunction
	ShouldBe bool
}

type OptionSpecification struct {
	Option ebm.MenuOption
	Checks []VisibilityCheck
}

type GroupSpecification struct {
	Group   ebm.MenuOptionGroup
	Options []OptionSpecification
}

type DynamicMenuSpecification []GroupSpecification

type FunctionBindings map[string]string
