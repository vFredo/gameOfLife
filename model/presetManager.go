package model

import (
	"errors"
	"io/ioutil"
  "os"
	"log"
	"math"
)

const PRESET_FOLDER = "./presets/"

type PresetManager struct {
	Presets []Preset
}

// Load every preset that is on the PRESET_FOLDER
func (pm *PresetManager) FetchPresets() {
	files, err := ioutil.ReadDir(PRESET_FOLDER)
	if err != nil {
		log.Fatalf("Can't read folder presets: %s", err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(PRESET_FOLDER + file.Name())
		if err != nil {
			log.Fatalf("Error when openning file '%s': %s", file.Name(), err)
		}

		var newPreset Preset
		newPreset.DecodeFromJson(content)
		pm.Presets = append(pm.Presets, newPreset)
	}
}

// TODO: Need to test this
// Create and save the preset in PRESET_FOLDER
func (pm *PresetManager) CreatePreset(name string, board []uint8, x int, y int) {
	max_width, max_height := math.Inf(-1), math.Inf(-1)
	min_width, min_height := math.Inf(1), math.Inf(1)
	var newPreset Preset
	newPreset.Name = name

	for i := x - 1; i > 0; i-- {
		for j := y - 1; j > 0; j-- {
			pos := (i * y) + j
			if (board[pos] & 0x01) == 0x01 {
				max_width = math.Max(max_width, float64(i))
				min_width = math.Min(min_width, float64(i))

				max_height = math.Max(max_height, float64(j))
				min_height = math.Min(min_height, float64(j))
				newPreset.AliveCells = append(newPreset.AliveCells, []uint{uint(i), uint(j)})
			}
		}
	}

	newPreset.Width = uint(max_width - min_width)
	newPreset.Height = uint(max_height - min_height)

	for i := 0; i < len(newPreset.AliveCells); i++ {
		newPreset.AliveCells[i][0] -= newPreset.Width
		newPreset.AliveCells[i][1] -= newPreset.Height
	}

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

func (pm *PresetManager) GetPreset(name string) (Preset, error) {
	for i := 0; i < len(pm.Presets); i++ {
		if pm.Presets[i].Name == name {
			return pm.Presets[i], nil
		}
	}
	return Preset{}, errors.New("Couldn't find preset with the name: " + name)
}
