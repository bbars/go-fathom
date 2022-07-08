package fathom

import (
	"encoding/json"
	
	"github.com/notnil/chess"
)

func (this WDL) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.String())
}

func (this tbMove) MarshalJSON() ([]byte, error) {
	return this.Move().MarshalJSON()
}

func (this tbMoveLong) MarshalJSON() ([]byte, error) {
	return this.Move().MarshalJSON()
}

type _tbResultJson struct {
	Move chess.Move `json:"move"`
	WDL WDL `json:"wdl"`
	DTZ int `json:"dtz"`
}

func (this tbResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(_tbResultJson{
		Move: this.Move(),
		WDL: this.WDL(),
		DTZ: this.DTZ(),
	})
}

type _tbRootMoveJson struct {
	Move chess.Move `json:"move"`
	PV []chess.Move `json:"pv"`
	Score int `json:"score"`
	Rank int `json:"rank"`
}

func (this tbRootMove) MarshalJSON() ([]byte, error) {
	return json.Marshal(_tbRootMoveJson{
		Move: this.Move(),
		PV: this.PV(),
		Score: this.Score(),
		Rank: this.Rank(),
	})
}
