package app

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
)

// Config hold global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var Config = koanf.New(".")

// ConfigInit configure application runtime
func ConfigInit(path string) {
	// override configuration with YAML
	err := Config.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		log.Printf("KoanfInit error config load: %s", err)
	} else {
		log.Printf("KoanfInit config load ok from %s: %+v", path, Config.All())
	}
}

// ConfigWatch let the application watch for changes
func ConfigWatch(path string) {
	if err := file.Provider(path).Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("watch error: %v", err)
			return
		}

		log.Printf("config changed. Reloading %+v", event)
		if err := Config.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Printf("ERROR unable to load configuration from file: %s", path)
		}
		Config.Print()
	}); err != nil {
		log.Printf("unable to watch config changes from file: %s", path)
	}
}
