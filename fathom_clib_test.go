package fathom

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

var tbDir string

func TestMain(m *testing.M) {
	flag.StringVar(&tbDir, "tbDir", "./tablebases", "Path to the directory containing Tablebase files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestNewFathom(t *testing.T) {
	tests := []struct {
		name    string
		tbDir   string
		wantErr bool
	}{
		{"wrong", "/non-existent/dir", true},
		{"okay", tbDir, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFathom(tt.tbDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFathom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != nil {
				got.Close()
			}
		})
	}
}

func TestWDL_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		this    WDL
		want    string
		wantErr bool
	}{
		{"Loss", Loss, `"Loss"`, false},
		{"BlessedLoss", BlessedLoss, `"Blessed Loss"`, false},
		{"Draw", Draw, `"Draw"`, false},
		{"CursedWin", CursedWin, `"Cursed Win"`, false},
		{"Win", Win, `"Win"`, false},
		{"<invalid>", -1, `"WDL???"`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWDL_String(t *testing.T) {
	tests := []struct {
		name string
		this WDL
		want string
	}{
		{"Loss:", Loss, `"Loss"`},
		{"BlessedLoss:", BlessedLoss, `"Blessed Loss"`},
		{"Draw:", Draw, `"Draw"`},
		{"CursedWin:", CursedWin, `"Cursed Win"`},
		{"Win:", Win, `"Win"`},
		{"<invalid>", -1, `"WDL???"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
