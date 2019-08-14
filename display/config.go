package display

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/linuxdeepin/go-x11-client/ext/randr"
	"pkg.deepin.io/lib/log"
)

type Config map[string]*ScreenConfig

type ScreenConfig struct {
	Custom  []*CustomModeConfig
	Mirror  *MirrorModeConfig
	Extend  *ExtendModeConfig
	OnlyOne *OnlyOneModeConfig
	Single  *MonitorConfig
}

type CustomModeConfig struct {
	Name     string
	Monitors []*MonitorConfig
}

type MirrorModeConfig struct {
	Monitors []*MonitorConfig
}

type ExtendModeConfig struct {
	Monitors []*MonitorConfig
}

type OnlyOneModeConfig struct {
	Monitors []*MonitorConfig
}

func (s *ScreenConfig) getMonitorConfigs(mode uint8, customName string) []*MonitorConfig {
	switch mode {
	case DisplayModeCustom:
		for _, custom := range s.Custom {
			if custom.Name == customName {
				return custom.Monitors
			}
		}
	case DisplayModeMirror:
		if s.Mirror == nil {
			return nil
		}
		return s.Mirror.Monitors

	case DisplayModeExtend:
		if s.Extend == nil {
			return nil
		}
		return s.Extend.Monitors

	case DisplayModeOnlyOne:
		if s.OnlyOne == nil {
			return nil
		}
		return s.OnlyOne.Monitors
	}

	return nil
}

func (s *ScreenConfig) getMonitorConfig(single bool, mode uint8, name, uuid string) *MonitorConfig {
	if single {
		return s.Single
	}

	switch mode {
	case DisplayModeCustom:
		for _, custom := range s.Custom {
			if custom.Name == name {
				return getMonitorConfigByUuid(custom.Monitors, uuid)
			}
		}
	case DisplayModeMirror:
		if s.Mirror == nil {
			return nil
		}
		return getMonitorConfigByUuid(s.Mirror.Monitors, uuid)

	case DisplayModeExtend:
		if s.Extend == nil {
			return nil
		}
		return getMonitorConfigByUuid(s.Extend.Monitors, uuid)

	case DisplayModeOnlyOne:
		if s.OnlyOne == nil {
			return nil
		}
		return getMonitorConfigByUuid(s.OnlyOne.Monitors, uuid)
	}

	return nil
}

func getMonitorConfigByUuid(configs []*MonitorConfig, uuid string) *MonitorConfig {
	for _, mc := range configs {
		if mc.UUID == uuid {
			return mc
		}
	}
	return nil
}

func setMonitorConfigsPrimary(configs []*MonitorConfig, uuid string) {
	for _, mc := range configs {
		if mc.UUID == uuid {
			mc.Primary = true
		} else {
			mc.Primary = false
		}
	}
}

func updateMonitorConfigsName(configs []*MonitorConfig, monitorMap map[randr.Output]*Monitor) {
	for _, mc := range configs {
		for _, m := range monitorMap {
			if mc.UUID == m.uuid {
				mc.Name = m.Name
				break
			}
		}
	}
}

func (s *ScreenConfig) setMonitorConfigs(mode uint8, customName string, configs []*MonitorConfig) {
	switch mode {
	case DisplayModeCustom:
		foundName := false
		for _, custom := range s.Custom {
			if custom.Name == customName {
				foundName = true
				custom.Monitors = configs
			}
		}

		// new custom
		if !foundName {
			s.Custom = []*CustomModeConfig{
				{
					Name:     customName,
					Monitors: configs,
				},
			}
		}

	case DisplayModeMirror:
		if s.Mirror == nil {
			s.Mirror = &MirrorModeConfig{}
		}
		s.Mirror.Monitors = configs

	case DisplayModeExtend:
		if s.Extend == nil {
			s.Extend = &ExtendModeConfig{}
		}
		s.Extend.Monitors = configs

	case DisplayModeOnlyOne:
		s.setMonitorConfigsOnlyOne(configs)
	}
}

func (s *ScreenConfig) setMonitorConfigsOnlyOne(configs []*MonitorConfig) {
	if s.OnlyOne == nil {
		s.OnlyOne = &OnlyOneModeConfig{}
	}
	oldConfigs := s.OnlyOne.Monitors
	var newConfigs []*MonitorConfig
	for _, cfg := range configs {
		if !cfg.Enabled {
			oldCfg := getMonitorConfigByUuid(oldConfigs, cfg.UUID)
			if oldCfg != nil {
				// 不设置 X,Y 是因为它们总是 0
				cfg.Width = oldCfg.Width
				cfg.Height = oldCfg.Height
				cfg.RefreshRate = oldCfg.RefreshRate
				cfg.Rotation = oldCfg.Rotation
				cfg.Reflect = oldCfg.Reflect
			} else {
				continue
			}
		}
		newConfigs = append(newConfigs, cfg)
	}
	s.OnlyOne.Monitors = newConfigs
}

type MonitorConfig struct {
	UUID        string
	Name        string
	Enabled     bool
	X           int16
	Y           int16
	Width       uint16
	Height      uint16
	Rotation    uint16
	Reflect     uint16
	RefreshRate float64
	Primary     bool
}

func loadConfig(filename string) (Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c Config) save(filename string) error {
	var data []byte
	var err error
	if logger.GetLogLevel() == log.LevelDebug {
		data, err = json.MarshalIndent(c, "", "    ")
		if err != nil {
			return err
		}
	} else {
		data, err = json.Marshal(c)
		if err != nil {
			return err
		}
	}

	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
