package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
)

const PRESET_FILE = "./presets/"

type Preset struct {
	Name       string   `json:"name"`
	Width      uint     `json:"width"`
	Height     uint     `json:"height"`
	AliveCells [][]uint `json:"cells"`
}

type PresetManager struct {
	Presets []Preset
}

// Load every preset that is on the PRESET_FILE folder
func (pm *PresetManager) FetchPresets() {
	files, err := ioutil.ReadDir(PRESET_FILE)
	if err != nil {
		log.Fatalf("Can't read folder presets: %s", err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(PRESET_FILE + file.Name())
		if err != nil {
			log.Fatalf("Error when openning file '%s': %s", file.Name(), err)
		}

		var newPreset Preset
		err = newPreset.DecodeFromJson(content)
		if err != nil {
			log.Fatalf("Error during decoding from json to preset: %s", err)
		}
		pm.Presets = append(pm.Presets, newPreset)
	}
}

// TODO: Need to test this
// Create and save the preset in PRESET_FILE
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

	newPreset.EncodeToJson()
	pm.Presets = append(pm.Presets, newPreset)
}

// Encode struct Preset into a json string
func (p *Preset) EncodeToJson() {
	// MarshalIndent to format the output
	encodedPreset, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Fatalf("Error occurred encoding struct Preset: %s", err.Error())
	}

	file, err := os.OpenFile(PRESET_FILE+p.Name+".json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error occurred while open/created file %s: %s", p.Name, err)
	}
	defer file.Close()

	file.Write(encodedPreset)
}

// Decode a json string into struct Preset
func (p *Preset) DecodeFromJson(buffer []byte) error {
	err := json.Unmarshal(buffer, &p)
	if err != nil {
		log.Fatalf("Error occurred decoding the json file: %s", err.Error())
	}
	return err
}
