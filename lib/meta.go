package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	KFS = ".KFS_META"
)

var user [2]int
var userEnabled bool

func SetUser(u, g int) {
	user = [2]int{u, g}
	userEnabled = true
}

func NewMetaInfo() MetaInfo {
	return MetaInfo{
		Context: make(map[string]interface{}),
	}
}

type Meta struct {
	Root     string
	MetaInfo map[string]MetaInfo
	mu       *sync.Mutex
}

func NewMeta(path string) *Meta {
	m := Meta{
		MetaInfo: make(map[string]MetaInfo),
		mu:       &sync.Mutex{},
	}
	err := m.Load(path)
	if err != nil {
		m.init(path)
	}
	return &m
}

func (m *Meta) init(path string) {
	m.Root = path
	b, _ := json.MarshalIndent(m, "", "  ")
	p := filepath.Join(m.Root, KFS)
	os.WriteFile(p, b, 0644)
	if userEnabled {
		os.Chown(p, user[0], user[1])
	}
}

func (m *Meta) Load(path string) error {
	if m.MetaInfo == nil {
		m.MetaInfo = make(map[string]MetaInfo)
	}
	if m.mu == nil {
		m.mu = &sync.Mutex{}
	}
	metaFile := filepath.Join(path, KFS)
	b, err := ioutil.ReadFile(metaFile)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = json.Unmarshal(b, m)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if path != m.Root {
		m.Root = path
		m.Write()
	}
	return nil
}

func (m *Meta) Merge(m2 *Meta) *Meta {
	// 	m.mu.Lock()
	// 	defer m.mu.Unlock()
	// 	m2.mu.Lock()
	// 	defer m2.mu.Unlock()
	for k, i2 := range m2.MetaInfo {
		if i1, ok := m.Get(k); ok {
			if len(i1.Label) > len(i2.Label) {
				i2.Label = i1.Label
			}
			if len(i1.Tags) > len(i2.Tags) {
				i2.Tags = i1.Tags
			}
			i2.Star = i1.Star || i2.Star
			if len(i1.OldLoc) > len(i2.OldLoc) {
				i2.OldLoc = i1.OldLoc
			}
			for ck, cv := range i1.Context {
				i2.Context[ck] = cv
			}
		}
		m.Set(k, i2)
	}
	return m
}

func (m *Meta) Get(name string) (MetaInfo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.MetaInfo[name]
	if val.Context == nil {
		val.Context = make(map[string]interface{})
	}
	return val, ok
}

func (m *Meta) Set(name string, val MetaInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetaInfo[name] = val
}

func (m *Meta) Del(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.MetaInfo, name)
}

func (m *Meta) Write() error {
	metaFile := filepath.Join(m.Root, KFS)
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = os.WriteFile(metaFile, b, 0644)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if userEnabled {
		err = os.Chown(metaFile, user[0], user[1])
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

type MetaInfo struct {
	Label   string
	Tags    []string
	Star    bool
	OldLoc  string
	Context map[string]interface{}
}
