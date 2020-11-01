package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestBuildPath(t *testing.T) {
	path, err := BuildPath("{\"id\":1}", "/v1/cards/:id", "/change/xy7-54")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if path != "/v1/cards/xy7-54" {
		t.Errorf("Error path[%s]", path)
	}
}

func TestHttpReq(t *testing.T) {
	data, err := HttpReq("GET", "https://api.pokemontcg.io/v1/cards/xy7-54", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	out.WriteTo(os.Stdout)
}

func TestBuildTmpl(t *testing.T) {
	compiled, err := BuildTmpl([]byte(data), tmpl)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf(compiled)
}

const tmpl = `
{{- $types := list "Colorless" "Darkness" "Dragon" "Fairy" "Fighting" "Fire" "Grass" "Lightning" "Metal" "Psychic" "Water" }}
{{- $newTypes := list }}
{{- range .card.types }}
	{{- $rand := randNumeric 3 }}
	{{- $max := len $types }}
	{{- $type := index $types (mod $rand $max) }}
	{{- $newTypes = append $newTypes $type }}
{{- end }}

{"message":"Change types for pokemon {{.card.name}}[{{.card.id}}] from {{.card.types}} to {{$newTypes}}"}
`
const data = `
{
  "card": {
    "id": "xy7-54",
    "name": "Gardevoir",
    "nationalPokedexNumber": 282,
    "imageUrl": "https://images.pokemontcg.io/xy7/54.png",
    "imageUrlHiRes": "https://images.pokemontcg.io/xy7/54_hires.png",
    "types": [
      "Fairy"
    ],
    "supertype": "Pokémon",
    "subtype": "Stage 2",
    "evolvesFrom": "Kirlia",
    "ability": {
      "name": "Bright Heal",
      "text": "Once during your turn (before your attack), you may heal 20 damage from each of your Pokémon.",
      "type": "Ability"
    },
    "hp": "130",
    "retreatCost": [
      "Colorless",
      "Colorless"
    ],
    "convertedRetreatCost": 2,
    "number": "54",
    "artist": "TOKIYA",
    "rarity": "Rare Holo",
    "series": "XY",
    "set": "Ancient Origins",
    "setCode": "xy7",
    "attacks": [
      {
        "cost": [
          "Colorless",
          "Colorless",
          "Colorless"
        ],
        "name": "Telekinesis",
        "text": "This attack does 50 damage to 1 of your opponent's Pokémon. This attack's damage isn't affected by Weakness or Resistance.",
        "damage": "",
        "convertedEnergyCost": 3
      }
    ],
    "resistances": [
      {
        "type": "Darkness",
        "value": "-20"
      }
    ],
    "weaknesses": [
      {
        "type": "Metal",
        "value": "×2"
      }
    ]
  }
}
`
