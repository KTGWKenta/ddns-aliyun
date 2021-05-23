package config

import (
	"flag"
	"os"

	"github.com/kentalee/eventbus"

	"github.com/kentalee/ddns-aliyun/common/events"
)

// 命令行参数注册、解析

type flagSubscriberEvents struct{}

func (f *flagSubscriberEvents) Init() error {
	registerConfigFlags()
	return nil
}

func (f *flagSubscriberEvents) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyInitialize: events.PriorityINI_CheckFlag,
	}
}

func (f *flagSubscriberEvents) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case events.EventKeyInitialize:
		return checkFlags()
	}
	return nil
}

// 配置文件加载、检查

type configSubscriberEvents struct{}

func (e *configSubscriberEvents) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyInitialize: events.PriorityINI_ParseConfig,
		events.EventKeyAfterCheck: events.PriorityAFT_CheckConfig,
	}
}

func (e *configSubscriberEvents) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case events.EventKeyInitialize:
		return loadConfigs()
	case events.EventKeyAfterCheck:
		return checkConfig()
	}
	return nil
}

type definesCollectorEvents struct {
	definesExportPath string
	definesCollector  *DefinesCollector
}

func (d *definesCollectorEvents) Init() error {
	flag.StringVar(&d.definesExportPath, DefinesExportFlag, "", "Path to export internal config defines")
	return nil
}

func (d *definesCollectorEvents) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyInitialize: events.PriorityINI_CollectDefines,
	}
}

func (d *definesCollectorEvents) OnEvent(key eventbus.Key, _ eventbus.Publisher) (err error) {
	switch key {
	case events.EventKeyInitialize:
		d.definesCollector = &DefinesCollector{}
		if err = loadDefines(d.definesCollector); err != nil {
			return err
		}
		if d.definesExportPath != "" {
			if err = exportDefines(d.definesExportPath); err != nil {
				return err
			}
			os.Exit(0)
		}
	}
	return nil
}

func init() {
	var err error
	err = events.Global().Register(&flagSubscriberEvents{})
	if err != nil {
		panic(err)
	}
	err = events.Global().Register(&configSubscriberEvents{})
	if err != nil {
		panic(err)
	}
	err = events.Global().Register(&definesCollectorEvents{})
	if err != nil {
		panic(err)
	}
}
