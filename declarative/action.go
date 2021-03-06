// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"errors"
)

import (
	"github.com/lxn/walk"
)

type Action struct {
	AssignTo    **walk.Action
	Text        string
	Image       interface{}
	OnTriggered walk.EventHandler
}

func (a Action) createAction(menu *walk.Menu) (*walk.Action, error) {
	action := walk.NewAction()

	if err := action.SetText(a.Text); err != nil {
		return nil, err
	}
	if err := setActionImage(action, a.Image); err != nil {
		return nil, err
	}

	if a.OnTriggered != nil {
		action.Triggered().Attach(a.OnTriggered)
	}

	if menu != nil {
		if err := menu.Actions().Add(action); err != nil {
			return nil, err
		}
	}

	if a.AssignTo != nil {
		*a.AssignTo = action
	}

	return action, nil
}

type ActionRef struct {
	Action *walk.Action
}

func (ar ActionRef) createAction(menu *walk.Menu) (*walk.Action, error) {
	if menu != nil {
		if err := menu.Actions().Add(ar.Action); err != nil {
			return nil, err
		}
	}

	return ar.Action, nil
}

type Menu struct {
	AssignTo       **walk.Menu
	AssignActionTo **walk.Action
	Text           string
	Image          interface{}
	Items          []MenuItem
}

func (m Menu) createAction(menu *walk.Menu) (*walk.Action, error) {
	if menu == nil {
		var err error
		if menu, err = walk.NewMenu(); err != nil {
			return nil, err
		}
	}

	subMenu, err := walk.NewMenu()
	if err != nil {
		return nil, err
	}

	action, err := menu.Actions().AddMenu(subMenu)
	if err != nil {
		return nil, err
	}

	if err := action.SetText(m.Text); err != nil {
		return nil, err
	}
	if err := setActionImage(action, m.Image); err != nil {
		return nil, err
	}

	for _, item := range m.Items {
		if _, err := item.createAction(subMenu); err != nil {
			return nil, err
		}
	}

	if m.AssignActionTo != nil {
		*m.AssignActionTo = action
	}
	if m.AssignTo != nil {
		*m.AssignTo = subMenu
	}

	return action, nil
}

type Separator struct {
}

func (s Separator) createAction(menu *walk.Menu) (*walk.Action, error) {
	action := walk.NewAction()

	if err := action.SetText("-"); err != nil {
		return nil, err
	}

	if menu != nil {
		if err := menu.Actions().Add(action); err != nil {
			return nil, err
		}
	}

	return action, nil
}

func addToActionList(list *walk.ActionList, actions []*walk.Action) error {
	for _, a := range actions {
		if err := list.Add(a); err != nil {
			return err
		}
	}

	return nil
}

func setActionImage(action *walk.Action, image interface{}) (err error) {
	var img *walk.Bitmap

	switch image := image.(type) {
	case nil:
		return nil

	case *walk.Bitmap:
		img = image

	case string:
		if img, err = walk.NewBitmapFromFile(image); err != nil {
			return
		}

	default:
		return errors.New("invalid type for Image")
	}

	return action.SetImage(img)
}

func CreateActions(items ...MenuItem) ([]*walk.Action, error) {
	var actions []*walk.Action

	for _, item := range items {
		action, err := item.createAction(nil)
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}
