package config

import (
	"path"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"

	"github.com/kentalee/errors"
	"github.com/kentalee/eventbus"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns/common/defines"
	"github.com/kentalee/ddns/common/events"
	"github.com/kentalee/ddns/common/typeUtils/cueAST"
	"github.com/kentalee/ddns/package/utils"
)

const (
	DefinesExportFlag = "cDef"
	DefinesCfgEntry   = "defines"
	DefinesFakePath   = "__internal__"
	EvtCollectDefines = "config.CollectConst"
)

var definesContent []byte

type DefinesCollector struct {
	runtime  cue.Runtime
	rwMutex  sync.RWMutex
	instance *cue.Instance
	buildDcl []ast.Decl
}

func (c *DefinesCollector) Runtime() *cue.Runtime {
	return &c.runtime
}

func (c *DefinesCollector) Init() (err error) {
	if c.instance, err = (&c.runtime).Compile(path.Join(DefinesFakePath, "fakePath.cue"), ""); err != nil {
		return err
	}
	c.buildDcl = []ast.Decl{
		&ast.Package{Name: ast.NewIdent(defines.ConfigCUEPackageName)},
	}
	return nil
}

func (c *DefinesCollector) Add(val interface{}, path ...string) (err error) {
	var newInstance *cue.Instance
	newInstance, err = c.instance.Fill(val, append([]string{DefinesCfgEntry}, path...)...)
	if err != nil {
		return errors.Note(err)
	}
	c.rwMutex.Lock()
	c.instance = newInstance
	c.rwMutex.Unlock()

	return nil
}

func (c *DefinesCollector) Expr(expr ast.Expr, path ...string) {
	var decl = &ast.Field{Label: ast.NewIdent(DefinesCfgEntry)}
	var _decl = decl
	for _, pathPart := range path {
		subField := &ast.Field{Label: ast.NewIdent(pathPart)}
		_decl.Value = ast.NewStruct(subField)
		_decl = subField
	}
	_decl.Value = expr
	c.buildDcl = append(c.buildDcl, decl)
}

func (c *DefinesCollector) Content() (content []byte, err error) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	var file *ast.File
	if file, err = cueAST.MergeRuntimeInstance(c.buildDcl, c.instance); err != nil {
		return nil, errors.Note(err)
	}
	if content, err = format.Node(file); err != nil {
		return nil, errors.Note(err)
	}
	return content, nil
}

func (c *DefinesCollector) Event() eventbus.Key {
	return EvtCollectDefines
}

func loadDefines(c *DefinesCollector) (err error) {
	if err = events.Global().Post(c); err != nil {
		return err
	}
	if definesContent, err = c.Content(); err != nil {
		return err
	}
	return nil
}

func exportDefines(path string) (err error) {
	err = utils.FileWrite(path, definesContent, false)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
