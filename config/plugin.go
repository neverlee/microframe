package config

import (
	"io/ioutil"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	td, err := time.ParseDuration(s)
	if err == nil {
		*d = Duration(td)
	}
	return err
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if len(b) > 2 {
		s = string(b[1 : len(b)-1])
	}

	td, err := time.ParseDuration(s)
	if err == nil {
		*d = Duration(td)
	}
	return err
}

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

type KeyIndex struct {
	Key   string
	Index int
}

var pluginki = 0
var mutexki sync.Mutex

func (ki *KeyIndex) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ki.Index = pluginki
	pluginki++
	return unmarshal(&ki.Key)
}

type RawYaml struct {
	unmarshaler func(interface{}) error
}

// func (r RawMessage) MarshalYAML() (interface{}, error) { }
func (r *RawYaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	r.unmarshaler = unmarshal
	return nil
}

func (r *RawYaml) Decode(i interface{}) error {
	if r.unmarshaler != nil {
		return r.unmarshaler(i)
	}
	return nil
}

type PluginConfig struct {
	Name   string  `yaml:"name"`
	Config RawYaml `yaml:"config"`
}

type BoardConf []*PluginConfig

func NewBoardConf(confPath string) (BoardConf, error) {
	var confBytes []byte
	var conf BoardConf
	var err error
	if confBytes, err = ioutil.ReadFile(confPath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(confBytes, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (b *BoardConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	mutexki.Lock()
	defer mutexki.Unlock()
	pluginki = 0
	var kimap map[KeyIndex]RawYaml
	if err := unmarshal(&kimap); err != nil {
		return err
	}
	*b = make(BoardConf, len(kimap))
	for ki, raw := range kimap {
		(*b)[ki.Index] = &PluginConfig{
			Name:   ki.Key,
			Config: raw,
		}
	}

	return nil
}
