package fathom

import (
	"flag"
	"os"
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

var tbDir string

func TestMain(m *testing.M) {
	flag.StringVar(&tbDir, "tbDir", "./tablebases_test", "Path to the directory containing Tablebase files")
	flag.Parse()
	os.Exit(m.Run())
}

func mustParseFen(t *testing.T, fen string) *chess.Position {
	pos := &chess.Position{}
	err := pos.UnmarshalText([]byte(fen))
	if err != nil {
		t.Fatal("broken test case: unable to parse FEN")
	}
	return pos
}

func mustNewFathom(t *testing.T) Fathom {
	f, err := NewFathom(tbDir)
	if err != nil {
		t.Fatal("broken test case: unable to instantiate Fathom")
	}
	return f
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
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.JSONEq(t, tt.want, string(got))
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
			got := tt.this.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_fathom_ProbeWDL(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		want    WDL
		wantErr bool
	}{
		{
			name:    "win",
			fen:     "8/4K3/8/8/8/7R/3k4/8 w - - 0 1",
			want:    Win,
			wantErr: false,
		},
		{
			name:    "loss",
			fen:     "8/4K3/8/8/8/7r/3k4/8 w - - 0 1",
			want:    Loss,
			wantErr: false,
		},
		{
			name:    "draw",
			fen:     "8/4K3/8/8/8/8/3k4/8 w - - 0 1",
			want:    Draw,
			wantErr: false,
		},
	}

	f := mustNewFathom(t)
	defer f.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := mustParseFen(t, tt.fen)

			got, err := f.ProbeWDL(pos)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
