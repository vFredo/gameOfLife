package game

import (
	"encoding/json"
	"log"
)

type Preset struct {
	Name       string   `json:"name"`
	Width      uint     `json:"width"`
	Height     uint     `json:"height"`
	AliveCells [][]uint `json:"cells"`
}

// Encode struct Preset into a json string
func (p *Preset) EncodeToJson() []byte {
	// MarshalIndent to format the output
	encodedPreset, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Printf("Error occurred encoding struct Preset: %s", err.Error())
	}
	return encodedPreset
}

// Decode a json string to get the values insed the preset itself
func (p *Preset) DecodeFromJson(buffer []byte) {
	err := json.Unmarshal(buffer, &p)
	if err != nil {
		log.Printf("Error occurred decoding the json file: %s", err.Error())
	}
}
