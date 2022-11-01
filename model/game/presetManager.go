package game

import (
	"errors"
	"log"
	"os"
)

const PRESET_FOLDER = "./presets/"

type PresetManager struct {
	currPresetIndex int
	Presets         []Preset
}

// Load all presets on the PRESET_FOLDER directory
func (pm *PresetManager) FetchPresets() {
	files, err := os.ReadDir(PRESET_FOLDER)
	if err != nil {
		log.Fatalf("Can't read folder presets: %s", err)
	}
	for _, file := range files {
		content, err := os.ReadFile(PRESET_FOLDER + file.Name())
		if err != nil {
			log.Fatalf("Error when openning file '%s': %s", file.Name(), err)
		}

		var newPreset Preset
		newPreset.DecodeFromJson(content)
		pm.Presets = append(pm.Presets, newPreset)
	}
}

// Create and save the preset in PRESET_FOLDER directory
func (pm *PresetManager) CreatePreset(name string, alive [][]uint, x uint, y uint) {
	var newPreset Preset

	newPreset.Name = name
	newPreset.Width = x
	newPreset.Height = y
	newPreset.AliveCells = alive

	encodedPreset := newPreset.EncodeToJson()
	pm.Presets = append(pm.Presets, newPreset)

	// Saving preset on PRESET_FOLDER
	file, err := os.OpenFile(PRESET_FOLDER+name+".json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error occurred while open/created file %s: %s", name, err)
	}
	defer file.Close()
	file.Write(encodedPreset)
}

// Get the preset that has the given name, if it doesn't exist, throw and error
func (pm *PresetManager) GetPreset(name string) (Preset, error) {
	for i := 0; i < len(pm.Presets); i++ {
		if pm.Presets[i].Name == name {
			pm.currPresetIndex = i
			return pm.Presets[i], nil
		}
	}
	return Preset{}, errors.New("Couldn't find preset with the name: " + name)
}

// Cycle between all available presets in the array Presets
func (pm *PresetManager) CyclePresets() (Preset, error) {
	if pm.currPresetIndex+1 < len(pm.Presets) {
		pm.currPresetIndex += 1
	} else {
		pm.currPresetIndex = 0
	}

	if len(pm.Presets) == 0 {
		return Preset{}, errors.New("there are no presets loaded in the PresetManager")
	}

	return pm.Presets[pm.currPresetIndex], nil
}
