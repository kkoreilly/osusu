package osusu

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// A Group is a group of users that can determine what to eat together
type Group struct {
	ID       int64
	Owner    int64
	Code     string
	Name     string
	Members  []int64
	Cuisines []string
}

// Groups is a slice of groups
type Groups []Group

// GroupJoin is the struct that contains data used to have a person join a group
type GroupJoin struct {
	GroupCode string
	UserID    int64
}

// CurrentGroup returns the current group state value
func CurrentGroup(ctx app.Context) Group {
	var group Group
	ctx.GetState("currentGroup", &group)
	return group
}

// SetCurrentGroup sets the current group state value to the given group
func SetCurrentGroup(group Group, ctx app.Context) {
	ctx.SetState("currentGroup", group, app.Persist)
}

// IsGroupNew gets whether the current group is a new group
func IsGroupNew(ctx app.Context) bool {
	var isGroupNew bool
	ctx.GetState("isGroupNew", &isGroupNew)
	return isGroupNew
}

// SetIsGroupNew sets whether the current group is a new group
func SetIsGroupNew(isGroupNew bool, ctx app.Context) {
	ctx.SetState("isGroupNew", isGroupNew, app.Persist)
}
