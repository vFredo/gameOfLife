package model

import (
	"encoding/json"
  "io/ioutil"
	"log"
)

type Preset struct {
	Name       string   `json:"name"`
	Width      uint     `json:"width"`
	Height     uint     `json:"height"`
	AliveCells [][]uint `json:"cells"`
}

type PresetManager struct {
	Presets []Preset
}

// Load every preset that is on the './presets/' folder
func (pm *PresetManager) FetchPresets() {
  files, err := ioutil.ReadDir("./presets/")
  if err != nil {
    log.Fatal("Can't read folder presets: ", err)
  }
  for _, file := range files {
    content, err := ioutil.ReadFile("./presets/" + file.Name())
    if err != nil {
      log.Fatal("Error when openning files: ", err)
    }

    var newPreset Preset
    err = newPreset.DecodeFromJson(content)
    if err != nil {
      log.Fatal("Error during decoding from json to preset: ", err)
    }
    pm.Presets = append(pm.Presets, newPreset)
  }
}

// Encode struct Preset into a json string
func (p *Preset) EncodeToJson() ([]byte, error) {
	encodedPreset, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error occurred encoding struct Preset: %s", err.Error())
	}
	return encodedPreset, nil
}

// Decode a json string into struct Preset
func (p *Preset) DecodeFromJson(buffer []byte) error {
	err := json.Unmarshal(buffer, &p)
	if err != nil {
		log.Fatalf("Error occurred decoding the json file: %s", err.Error())
	}
	return err
}
