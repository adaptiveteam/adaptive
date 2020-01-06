package EngagementBuilder

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
)

// MenuOptionBuilder pattern code
// deprecated. Use MenuOption constructor function
type MenuOptionBuilder struct {
	menuOption *model.MenuOption
}

func NewMenuOptionBuilder() *MenuOptionBuilder {
	menuOption := &model.MenuOption{}
	b := &MenuOptionBuilder{menuOption: menuOption}
	return b
}

func (b *MenuOptionBuilder) Text(text string) *MenuOptionBuilder {
	b.menuOption.Text = text
	return b
}

func (b *MenuOptionBuilder) Value(value string) *MenuOptionBuilder {
	b.menuOption.Value = value
	return b
}

func (b *MenuOptionBuilder) Description(text string) *MenuOptionBuilder {
	b.menuOption.Description = text
	return b
}

func (b *MenuOptionBuilder) Build() (*model.MenuOption, error) {
	return b.menuOption, nil
}

// MenuOptionGroup builder pattern code
type MenuOptionGroupBuilder struct {
	menuOptionGroup *model.MenuOptionGroup
}

func NewMenuOptionGroupBuilder() *MenuOptionGroupBuilder {
	menuOptionGroup := &model.MenuOptionGroup{}
	b := &MenuOptionGroupBuilder{menuOptionGroup: menuOptionGroup}
	return b
}

func (b *MenuOptionGroupBuilder) Text(text string) *MenuOptionGroupBuilder {
	b.menuOptionGroup.Text = text
	return b
}

func (b *MenuOptionGroupBuilder) Options(options []model.MenuOption) *MenuOptionGroupBuilder {
	b.menuOptionGroup.Options = options
	return b
}

func (b *MenuOptionGroupBuilder) Build() (*model.MenuOptionGroup, error) {
	return b.menuOptionGroup, nil
}
