package config

import (
	"sync"
	"time"
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

type boardConfig []*PluginConfig

func (b *boardConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	mutexki.Lock()
	defer mutexki.Unlock()
	pluginki = 0
	var kimap map[KeyIndex]RawYaml
	if err := unmarshal(&kimap); err != nil {
		return err
	}
	*b = make(boardConfig, len(kimap))
	for ki, raw := range kimap {
		(*b)[ki.Index] = &PluginConfig{
			Name:   ki.Key,
			Config: raw,
		}
	}

	return nil
}
