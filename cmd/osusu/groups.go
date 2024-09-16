package main

import (
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/views"
	"github.com/kkoreilly/osusu/osusu"
)

func groups(b *core.Body) {
	d := core.NewBody("groups")
	core.NewText(d).SetType(core.TextHeadlineMedium).SetText("Join group")
	groupCode := ""
	views.NewValue(d, &groupCode)
	core.NewButton(d).SetText("Join group")
	core.NewText(d).SetType(core.TextHeadlineMedium).SetText("Create group")
	newGroup := &osusu.Group{OwnerID: curUser.ID, Owner: *curUser, Members: []osusu.User{*curUser}}
	views.NewStructView(d).SetStruct(newGroup)
	core.NewButton(d).SetText("Create group").OnClick(func(e events.Event) {
		err := osusu.DB.Create(newGroup).Error
		if err != nil {
			core.ErrorDialog(d, err)
			return
		}
		curGroup = newGroup
		curUser.GroupID = newGroup.ID
		err = osusu.DB.Save(curUser).Error
		if err != nil {
			core.ErrorDialog(d, err)
		}
		d.Close()
	})
	d.NewFullDialog(b).Run()
}
