/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"fmt"
)

var (
	errPropInvalid      = fmt.Errorf("Invalid xsettings property")
	errPropNotFound     = fmt.Errorf("This property not found")
	errPropTypeNotMatch = fmt.Errorf("This property's type not match")
)

func (m *XSManager) SetInteger(prop string, v int32) error {
	var setting = xsSetting{
		sType: settingTypeInteger,
		prop:  prop,
		value: v,
	}

	err := m.setSettings([]xsSetting{setting})
	if err != nil {
		logger.Debugf("Set '%s' to '%v' failed: %v", prop, v, err)
		return err
	}
	m.setGSettingsByXProp(prop, v)

	return nil
}

func (m *XSManager) GetInteger(prop string) (int32, error) {
	v, sType, err := m.getSettingValue(prop)
	if err != nil {
		logger.Debug("Get '%s' value failed: %v", prop, err)
		return -1, err
	}

	if sType != settingTypeInteger {
		return -1, errPropTypeNotMatch
	}

	return v.(*integerValueInfo).value, nil
}

func (m *XSManager) SetString(prop, v string) error {
	var setting = xsSetting{
		sType: settingTypeString,
		prop:  prop,
		value: v,
	}

	err := m.setSettings([]xsSetting{setting})
	if err != nil {
		logger.Debugf("Set '%s' to '%v' failed: %v", prop, v, err)
		return err
	}
	m.setGSettingsByXProp(prop, v)

	return nil
}

func (m *XSManager) GetString(prop string) (string, error) {
	v, sType, err := m.getSettingValue(prop)
	if err != nil {
		logger.Debug("Get '%s' value failed: %v", prop, err)
		return "", err
	}

	if sType != settingTypeString {
		return "", errPropTypeNotMatch
	}

	return v.(*stringValueInfo).value, nil
}

func (m *XSManager) SetColor(prop string, v [4]int16) error {
	var setting = xsSetting{
		sType: settingTypeColor,
		prop:  prop,
		value: v,
	}

	err := m.setSettings([]xsSetting{setting})
	if err != nil {
		logger.Debugf("Set '%s' to '%v' failed: %v", prop, v, err)
		return err
	}
	m.setGSettingsByXProp(prop, v)

	return nil
}

func (m *XSManager) GetColor(prop string) ([4]int16, error) {
	v, sType, err := m.getSettingValue(prop)
	if err != nil {
		logger.Debug("Get '%s' value failed: %v", prop, err)
		return [4]int16{}, err
	}

	if sType != settingTypeColor {
		return [4]int16{}, errPropTypeNotMatch
	}

	tmp := v.(*colorValueInfo)

	return [4]int16{tmp.red, tmp.blue, tmp.green, tmp.alpha}, nil
}

func (m *XSManager) getSettingValue(prop string) (interface{}, int8, error) {
	datas, err := getSettingPropValue(m.owner, m.conn)
	if err != nil {
		return nil, -1, err
	}

	xsInfo := marshalSettingData(datas)
	item := xsInfo.getPropItem(prop)
	if item == nil {
		return nil, -1, errPropNotFound
	}

	return item.value, item.header.sType, nil
}

func (m *XSManager) setGSettingsByXProp(prop string, v interface{}) {
	info := gsInfos.getInfoByXSKey(prop)
	if info == nil {
		return
	}

	info.setKeyValue(m.gs, v)
}
