package config

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	_ "cuelang.org/go/pkg"
	"go.uber.org/zap"

	"github.com/kentalee/errors"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns-aliyun/package/utils"
)

const (
	configPathFlag  = "c"
	configPathRegex = `.*\.cue$`
)

var (
	configPath               string
	configPathPattern        = regexp.MustCompile(configPathRegex)
	ErrUnspecifiedConfigPath = errors.New("9e54a67400030001", "unspecified config path")
	ErrConfigNotFound        = errors.New("9e54a67400030002", "config not found")
	ErrInvalidConfigPath     = errors.New("9e54a67400030003", "invalid config path")
	ErrInvalidConfigFile     = errors.New("9e54a67400030004", "invalid config file")
	ErrEmptyConfigDir        = errors.New("9e54a67400030005", "empty config dir")
)

func LoadFromPath(path string) (*cue.Instance, error) {
	var err error
	var runtime = &cue.Runtime{}
	// empty config path
	if path = filepath.Clean(strings.TrimSpace(path)); path == "" {
		return nil, errors.Note(ErrUnspecifiedConfigPath)
	}
	// check config status
	var pathStat os.FileInfo
	var pathInfo = zap.String("path", path)
	if path, pathStat, err = utils.FileStat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Note(ErrConfigNotFound)
		}
		return nil, errors.Because(ErrInvalidConfigPath, err, pathInfo)
	}
	var instanceBuilder = build.NewContext().NewInstance(path, nil)
	// load global defines
	if err = instanceBuilder.AddFile(definesContentPath, definesContent); err != nil {
		return nil, errors.Because(ErrInvalidConfigFile, err, zap.String("path", definesContentPath))
	}
	// load internal template
	if err = instanceBuilder.AddFile(innerTemplatePath, innerTemplate); err != nil {
		return nil, errors.Because(ErrInvalidConfigFile, err, zap.String("path", innerTemplatePath))
	}
	// load config files
	if pathStat.IsDir() {
		var pathList []string
		if pathList, err = utils.DirList(path, configPathPattern, true); err != nil {
			return nil, errors.Because(ErrInvalidConfigPath, err, pathInfo)
		}
		if len(pathList) == 0 {
			return nil, errors.Because(ErrInvalidConfigPath, ErrEmptyConfigDir, pathInfo)
		}
		for _, subPath := range pathList {
			if err = instanceBuilder.AddFile(subPath, nil); err != nil {
				return nil, errors.Because(ErrInvalidConfigFile, err, zap.String("path", subPath))
			}
		}
	} else if err = instanceBuilder.AddFile(path, nil); err != nil {
		return nil, errors.Because(ErrInvalidConfigFile, err, pathInfo)
	}
	if err = instanceBuilder.Complete(); err != nil {
		return nil, errors.Note(err, pathInfo)
	}
	var cfgInstance *cue.Instance
	if cfgInstance, err = runtime.Build(instanceBuilder); err != nil {
		return nil, errors.Because(ErrInvalidConfigPath, err, pathInfo)
	}
	return cfgInstance, nil
}

func registerConfigFlags() {
	flag.StringVar(&configPath, configPathFlag, "config", "Path to config file or dir")
}

func checkFlags() (err error) {
	configPath = filepath.Clean(strings.TrimSpace(configPath))
	if configPath == "" {
		log.Warn("Please specific config path with arg `-c`")
		return errors.Note(ErrUnspecifiedConfigPath)
	}
	return nil
}

func loadConfigs() error {
	instance, err := LoadFromPath(configPath)
	if err != nil {
		log.Error("Failed to load configs!")
		err = errors.Note(err)
		return err
	}
	configRWLocker.Lock()
	configInstance = instance
	configRWLocker.Unlock()
	return nil
}

func checkConfig() error {
	// 预留，用作额外的配置检查
	return nil
}
