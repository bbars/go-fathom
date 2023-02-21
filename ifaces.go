package fathom

import (
	"github.com/notnil/chess"
)

type Fathom interface {
	// ProbeWDL probes the win-draw-loss (WDL) table.
	//
	// Returns one of: [Loss], [BlessedLoss], [Draw], [CursedWin] or [Win].
	// When error occurs, WDL results in [WDLUnknown].
	//
	// NOTES:
	// - Engines should use this function during search.
	ProbeWDL(chessPosition *chess.Position) (WDL, error)

	// ProbeRoot probes the Distance-To-Zero (DTZ) table.
	//
	// The suggested move is guaranteed to preserve the WDL value.
	// Possible errors: [ErrCheckmate], [ErrStalemate] (and other).
	//
	// NOTES:
	// - Engines can use this function to probe at the root. This function should
	//   not be used during search.
	// - DTZ tablebases can suggest unnatural moves, especially for losing
	//   positions. Engines may prefer to perform traditional search combined with WDL
	//   move filtering using the alternative results array.
	// - This function is NOT thread safe. For engines this function should only
	//   be called once at the root per search.
	ProbeRoot(chessPosition *chess.Position) (TbMove, []TbResult, error)

	// ProbeRootDTZ uses the DTZ tables to rank and score all root moves.
	ProbeRootDTZ(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error)

	// ProbeRootWDL uses the WDL tables to rank and score all root moves.
	//
	// NOTES:
	// - This is a fallback for the case that some or all DTZ tables are missing.
	ProbeRootWDL(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error)

	// Close frees any resources allocated by an instance.
	Close()
}

type TbMove interface {
	Move() chess.Move
}

// TbResult is a result value comprising:
// suggested move, the WDL value and the DTZ value.
type TbResult interface {
	TbMove

	WDL() WDL
	DTZ() int
}

// TbRootMove suggests a move, a rank, a score, and a predicted principal variation.
type TbRootMove interface {
	TbMove

	// PV stands for Principal Variation
	PV() []chess.Move
	Score() int
	Rank() int
}
