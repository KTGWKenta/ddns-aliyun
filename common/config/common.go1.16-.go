// +build !go1.16

package config

import (
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/kentalee/errors"
	"github.com/kentalee/eventbus"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns-aliyun/common/events"
	"github.com/kentalee/ddns-aliyun/package/utils"
)

var innerTemplate string

const templatePathFlag = "cTpl"

type templateSubscriber struct {
	path string
}

func (f templateSubscriber) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyInitialize: events.PriorityINI_CheckFlag,
	}
}

func (f *templateSubscriber) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case events.EventKeyInitialize:
		return loadTemplate(f.path)
	}
	return nil
}

func init() {
	subscriber := templateSubscriber{}
	registerTemplateFlags(&subscriber.path)
	if err := events.Global().Register(&subscriber); err != nil {
		log.Fatal(err)
	}
}

func registerTemplateFlags(path *string) {
	flag.StringVar(path, templatePathFlag, "", "Path to config template file "+
		"(this feature will be deprecated in go1.16 (use embed template instead))")
}

func loadTemplate(path string) (err error) {
	if path == "" {
		log.Fatalf("config template not found. please specify template path with command-line arg `%s`", templatePathFlag)
		return nil
	}
	var pathStat os.FileInfo
	if path, pathStat, err = utils.FileStat(path); err != nil {
		if os.IsNotExist(err) {
			return errors.Note(ErrConfigNotFound)
		}
		return errors.Because(ErrInvalidConfigPath, err, zap.String("path", path))
	}
	if pathStat.IsDir() || !configPathPattern.MatchString(pathStat.Name()) {
		return errors.Note(ErrInvalidConfigPath, zap.String("expected", ".cue file"), zap.String("got", "directory"))
	}
	var byteContent []byte
	if byteContent, err = utils.FileContent(path); err != nil {
		return errors.Because(ErrInvalidConfigFile, err)
	}
	innerTemplate = string(byteContent)
	return nil
}
