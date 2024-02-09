package plugin

import (
	"encoding/json"
	"log"
	"os"
)

const (
	pluginDir           = "plugin"
	pluginStateFile     = "plugin_sate.json"
	pluginStateFilePath = pluginDir + "/" + pluginStateFile
)

type State struct {
	PluginName string
	ModuleName string
	IsEnable   bool
}

func ReadPluginState() ([]State, error) {
	_, err := os.Stat(pluginStateFilePath)
	if err != nil {
		return []State{}, nil
	}

	//const pluginStateFileLockPath = pluginStateFilePath + ".lock"
	content, err := os.ReadFile(pluginStateFilePath)
	if err != nil {
		log.Println("failed to load plugin state")
		return nil, err
	}

	var states []State
	err = json.Unmarshal(content, &states)
	if err != nil {
		log.Println("failed to parse plugin state")
		return nil, err
	}

	return states, nil
}

func WritePluginState(states []State) error {
	// TODO
	return nil
}
