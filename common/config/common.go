package config

import (
	"context"
	"path"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"go.uber.org/zap"

	"github.com/kentalee/errors"
)

const (
	definesContentPath = "__definesContent__"
	innerTemplatePath  = "__innerTemplate__"
)

var configInstance *cue.Instance
var configRWLocker sync.RWMutex

func Lookup(fieldPath ...string) cue.Value {
	configRWLocker.RLock()
	defer configRWLocker.RUnlock()
	if configInstance == nil {
		return cue.Value{}
	}
	val := configInstance.Lookup(fieldPath...)
	return val
}

func Set(value interface{}, fieldPath ...string) error {
	configRWLocker.Lock()
	defer configRWLocker.Unlock()
	instance, err := configInstance.Fill(value, fieldPath...)
	if err != nil {
		return err
	}
	configInstance = instance
	return nil
}

type CUEDecoder interface {
	DecodeCUE(cfg cue.Value) error
}

type CUECtxDecoder interface {
	DecodeCUE(ctx context.Context, cfg cue.Value) error
}

func Decode(cfg cue.Value, target interface{}) (err error) {
	if decoder, ok := target.(CUEDecoder); ok {
		if err = decoder.DecodeCUE(cfg); err != nil {
			return errors.Note(err, zap.String("position", RefInfo(cfg)))
		}
	} else if decoder, ok := target.(CUECtxDecoder); ok {
		if err = decoder.DecodeCUE(context.Background(), cfg); err != nil {
			return errors.Note(err, zap.String("position", RefInfo(cfg)))
		}
	} else {
		if err = cfg.Decode(target); err != nil {
			return errors.Note(err, zap.String("position", RefInfo(cfg)))
		}
	}
	return nil
}

func DecodeWithCtx(ctx context.Context, cfg cue.Value, target interface{}) (err error) {
	if decoder, ok := target.(CUEDecoder); ok {
		if err = decoder.DecodeCUE(cfg); err != nil {
			return errors.Note(err, zap.String("position", RefInfo(cfg)))
		}
	} else if decoder, ok := target.(CUECtxDecoder); ok {
		if err = decoder.DecodeCUE(ctx, cfg); err != nil {
			return errors.Note(err, zap.String("position", RefInfo(cfg)))
		}
	} else {
		if err = cfg.Decode(target); err != nil {
			return errors.Note(err, zap.String("position", RefInfo(cfg)))
		}
	}
	return nil
}

func RefInfo(val cue.Value) string {
	ins, pathArr := val.Reference()
	if ins != nil {
		return path.Clean(path.Join(ins.Dir, ins.DisplayName)) + ":" + strings.Join(pathArr, ".")
	} else {
		return val.Pos().String()
	}
}

var FieldsLookupOption = []cue.Option{
	cue.Final(),
	cue.Concrete(true),
	cue.DisallowCycles(true),
}
