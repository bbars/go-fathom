package fathom

import (
	"encoding/json"

	"github.com/notnil/chess"
)

func (wdl WDL) MarshalJSON() ([]byte, error) {
	return json.Marshal(wdl.String())
}

func (m tbMove) MarshalJSON() ([]byte, error) {
	return m.Move().MarshalJSON()
}

func (ml tbMoveLong) MarshalJSON() ([]byte, error) {
	return ml.Move().MarshalJSON()
}

type _tbResultJson struct {
	Move chess.Move `json:"move"`
	WDL  WDL        `json:"wdl"`
	DTZ  int        `json:"dtz"`
}

func (r tbResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(_tbResultJson{
		Move: r.Move(),
		WDL:  r.WDL(),
		DTZ:  r.DTZ(),
	})
}

type _tbRootMoveJson struct {
	Move  chess.Move   `json:"move"`
	PV    []chess.Move `json:"pv"`
	Score int          `json:"score"`
	Rank  int          `json:"rank"`
}

func (rm tbRootMove) MarshalJSON() ([]byte, error) {
	return json.Marshal(_tbRootMoveJson{
		Move:  rm.Move(),
		PV:    rm.PV(),
		Score: rm.Score(),
		Rank:  rm.Rank(),
	})
}
