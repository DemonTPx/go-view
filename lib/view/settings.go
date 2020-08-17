package view

import (
	"encoding/json"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type WindowSettings struct {
	X uint32
	Y uint32
	W uint32
	H uint32
}

type Settings struct {
	Window WindowSettings
}

var DefaultSettings = Settings{
	Window: WindowSettings{
		X: 0,
		Y: 0,
		W: 1200,
		H: 900,
	},
}

const SettingsFilename = "settings.json"

func LoadSettings() Settings {
	settingPath := sdl.GetPrefPath("demontpx", "go-view")
	file, err := os.Open(filepath.Join(settingPath, SettingsFilename))
	if err != nil {
		log.Printf("could not open settings file %s%s\n", settingPath, SettingsFilename)
		return DefaultSettings
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("failed to read settings file: %s", err)
		return DefaultSettings
	}
	var settings Settings
	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		log.Printf("failed to unmarshal settings: %s", err)
		return DefaultSettings
	}
	return settings
}

func SaveSettings(settings Settings) {
	settingPath := sdl.GetPrefPath("demontpx", "go-view")
	file, err := os.OpenFile(filepath.Join(settingPath, SettingsFilename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Printf("failed to open settings file for writing: %s", err)
		return
	}
	defer file.Close()

	bytes, err := json.Marshal(settings)
	if err != nil {
		log.Printf("failed to marshal settings: %s", err)
		return
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("failed to write settings file: %s", err)
		return
	}
}
