// Copyright 2018 cg33.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package menu

import (
	"github.com/compty001/go-admin/modules/auth"
	"github.com/compty001/go-admin/modules/db"
	"github.com/compty001/go-admin/modules/language"
	"strconv"
	"sync"
	"sync/atomic"
)

type Menu struct {
	GlobalMenuList   []MenuItem
	GlobalMenuOption []map[string]string
	MaxOrder         int64
}

var GlobalMenu = &Menu{
	GlobalMenuList:   []MenuItem{},
	GlobalMenuOption: []map[string]string{},
	MaxOrder:         0,
}

var InitMenuOnce = &Once{
	done: 0,
}

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *Once) unlock() {
	atomic.StoreUint32(&o.done, 0)
}

func InitMenu(user auth.User) {
	SetGlobalMenu(user)
}

func GetGlobalMenu(user auth.User) *Menu {
	InitMenu(user)
	return GlobalMenu
}

func Unlock() {
	InitMenuOnce.unlock()
}

func SetGlobalMenu(user auth.User) {

	var (
		menus      []map[string]interface{}
		menuOption = make([]map[string]string, 0)
	)

	if user.IsSuperAdmin() {
		menus, _ = db.Table("goadmin_menu").
			Where("id", ">", 0).
			OrderBy("order", "asc").
			All()
	} else {

		var ids []interface{}
		for i := 0; i < len(user.Menus); i++ {
			ids = append(ids, user.Menus[i])
		}

		menus, _ = db.Table("goadmin_menu").
			WhereIn("id", ids).
			OrderBy("order", "asc").
			All()

	}

	var title string
	for i := 0; i < len(menus); i++ {

		if menus[i]["type"].(int64) == 1 {
			title = language.Get(menus[i]["title"].(string))
		} else {
			title = menus[i]["title"].(string)
		}

		menuOption = append(menuOption, map[string]string{
			"id":    strconv.FormatInt(menus[i]["id"].(int64), 10),
			"title": title,
		})
	}

	menulist := ConstructMenuTree(menus, 0)

	GlobalMenu.GlobalMenuOption = menuOption
	GlobalMenu.GlobalMenuList = menulist
	GlobalMenu.MaxOrder = menus[len(menus)-1]["parent_id"].(int64)
}

func (menu *Menu) SexMaxOrder(order int64) {
	menu.MaxOrder = order
}

func ConstructMenuTree(menus []map[string]interface{}, parentId int64) []MenuItem {

	branch := make([]MenuItem, 0)

	var title string
	for j := 0; j < len(menus); j++ {
		if parentId == menus[j]["parent_id"].(int64) {

			childList := ConstructMenuTree(menus, menus[j]["id"].(int64))

			if menus[j]["type"].(int64) == 1 {
				title = language.Get(menus[j]["title"].(string))
			} else {
				title = menus[j]["title"].(string)
			}

			child := MenuItem{
				Name:         title,
				ID:           strconv.FormatInt(menus[j]["id"].(int64), 10),
				Url:          menus[j]["uri"].(string),
				Icon:         menus[j]["icon"].(string),
				Active:       "",
				ChildrenList: childList,
			}

			branch = append(branch, child)
		}
	}

	return branch
}

func GetMenuItemById(id string) MenuItem {
	menu, _ := db.Table("goadmin_menu").Find(id)

	return MenuItem{
		Name:         menu["title"].(string),
		ID:           strconv.FormatInt(menu["id"].(int64), 10),
		Url:          menu["uri"].(string),
		Icon:         menu["icon"].(string),
		Active:       "",
		ChildrenList: []MenuItem{},
	}
}

func (menu *Menu) SetActiveClass(path string) *Menu {

	for i := 0; i < len((*menu).GlobalMenuList); i++ {
		(*menu).GlobalMenuList[i].Active = ""
	}

	for i := 0; i < len((*menu).GlobalMenuList); i++ {
		if (*menu).GlobalMenuList[i].Url == path {
			(*menu).GlobalMenuList[i].Active = "active"
			return menu
		} else {
			for j := 0; j < len((*menu).GlobalMenuList[i].ChildrenList); j++ {
				if (*menu).GlobalMenuList[i].ChildrenList[j].Url == path {
					(*menu).GlobalMenuList[i].Active = "active"
					return menu
				} else {
					(*menu).GlobalMenuList[i].Active = ""
				}
			}
		}
	}

	return menu
}

type MenuItem struct {
	Name         string
	ID           string
	Url          string
	Icon         string
	Active       string
	ChildrenList []MenuItem
}

func (menu *Menu) GetEditMenuList() []MenuItem {
	return (*menu).GlobalMenuList
}
