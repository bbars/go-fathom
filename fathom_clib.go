// Package fathom is a Syzygy tablebase probe implementation.
//
// This package wraps a fork of an original Fathom tool
// written in C: https://github.com/jdart1/Fathom.
//
// Also it uses well-known external package https://github.com/notnil/chess
// as a public communication layer.
//
// Methods allow to probe DTZ and WDL tables (*.rtbw and *.rtbz Syzygy tablebases).
package fathom

// #cgo CFLAGS: -std=c99 -Ifathom-copy/src
// #include <fathom-copy/src/tbprobe.c>
import "C"
import (
	// "unsafe"
	"fmt"

	"github.com/notnil/chess"
)

type cPosType C.Pos

type cCastling int

//goland:noinspection GoSnakeCaseUsage
const (
	castlingUnknown = cCastling(0)
	castling_K      = cCastling(C.TB_CASTLING_K)
	castling_Q      = cCastling(C.TB_CASTLING_Q)
	castling_k      = cCastling(C.TB_CASTLING_k)
	castling_q      = cCastling(C.TB_CASTLING_q)
)

// WDL win-draw-lose result status.
type WDL int

const (
	WDLUnknown  = WDL(-1)
	Loss        = WDL(C.TB_LOSS)
	BlessedLoss = WDL(C.TB_BLESSED_LOSS)
	Draw        = WDL(C.TB_DRAW)
	CursedWin   = WDL(C.TB_CURSED_WIN)
	Win         = WDL(C.TB_WIN)
)

func (wdl WDL) String() string {
	switch wdl {
	case Loss:
		return "Loss"
	case BlessedLoss:
		return "Blessed Loss"
	case Draw:
		return "Draw"
	case CursedWin:
		return "Cursed Win"
	case Win:
		return "Win"
	}
	return "WDL???"
}

/////////////////////////////////////////////////////

type fathom struct {
	tbDir string
}

var _ Fathom = &fathom{}

// NewFathom creates a new Fathom reader.
// If no tablebase files are found, then error is returned.
//
// Possible errors: [ErrNoTablebases]
//
// If you don't have tablebase files yet,
// follow this link (https://github.com/niklasf/shakmaty-syzygy/tree/master/tables)
// and explore suggested TXT.
// These TXT files contain URLs to tablebase files you should download.
func NewFathom(tbDir string) (Fathom, error) {
	cTbDir := C.CString(tbDir)
	res := C.tb_init(cTbDir)
	if res != true {
		return nil, fmt.Errorf("go-fathom: unable to init tablebases")
	}
	if C.TB_LARGEST == 0 {
		C.tb_free()
		return nil, ErrNoTablebases
	}
	return &fathom{
		tbDir: tbDir,
	}, nil
}

func (f *fathom) ProbeWDL(chessPosition *chess.Position) (WDL, error) {
	pos := &position{chessPosition}
	cPos := pos.cPos()

	res := C.tb_probe_wdl(
		cPos.white,
		cPos.black,
		cPos.kings,
		cPos.queens,
		cPos.rooks,
		cPos.bishops,
		cPos.knights,
		cPos.pawns,
		C.unsigned(cPos.rule50),     // 0, // C.unsigned(), // not supported by Fathom
		C.unsigned(pos.cCastling()), // 0, // C.unsigned(), // not supported by Fathom
		C.unsigned(cPos.ep),
		cPos.turn,
	)
	if res == C.TB_RESULT_FAILED {
		return WDLUnknown, fmt.Errorf("go-fathom: tb_probe_wdl failed")
	}
	return WDL(res), nil
}

func (f *fathom) ProbeRoot(chessPosition *chess.Position) (TbMove, []TbResult, error) {
	pos := &position{chessPosition}
	cPos := pos.cPos()

	var cResults [C.TB_MAX_MOVES]C.unsigned
	res := C.tb_probe_root(
		cPos.white,
		cPos.black,
		cPos.kings,
		cPos.queens,
		cPos.rooks,
		cPos.bishops,
		cPos.knights,
		cPos.pawns,
		C.unsigned(cPos.rule50),
		C.unsigned(pos.cCastling()),
		C.unsigned(cPos.ep),
		cPos.turn,
		&cResults[0],
	)
	switch res {
	case C.TB_RESULT_FAILED:
		return nil, nil, fmt.Errorf("go-fathom: probe_root failed")
	case C.TB_RESULT_CHECKMATE:
		return nil, nil, ErrCheckmate
	case C.TB_RESULT_STALEMATE:
		return nil, nil, ErrStalemate
	}
	results := make([]TbResult, 0, len(cResults))
	for _, cResult := range cResults {
		if cResult == C.TB_RESULT_FAILED {
			break
		}
		results = append(results, newTbResult(cResult))
	}
	return tbMoveLong(res), results, nil
}

func (f *fathom) ProbeRootDTZ(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error) {
	pos := &position{chessPosition}
	cPos := pos.cPos()

	var cResults C.struct_TbRootMoves = C.struct_TbRootMoves{
		0,
		[C.TB_MAX_MOVES]C.struct_TbRootMove{},
	}
	res := C.tb_probe_root_dtz(
		cPos.white,
		cPos.black,
		cPos.kings,
		cPos.queens,
		cPos.rooks,
		cPos.bishops,
		cPos.knights,
		cPos.pawns,
		C.unsigned(cPos.rule50),
		C.unsigned(pos.cCastling()),
		C.unsigned(cPos.ep),
		cPos.turn,
		false, // hasRepeated
		C.bool(useRule50),
		&cResults,
	)
	// fmt.Println("res", res)
	// fmt.Println("cResults", cResults)
	results := make([]TbRootMove, cResults.size)
	for i := 0; i < len(results); i++ {
		results[i] = newRootMove(cResults.moves[i])
	}

	if res == 0 {
		return results, fmt.Errorf("go-fathom: not all probes were successful")
	}
	return results, nil
}

func (f *fathom) ProbeRootWDL(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error) {
	pos := &position{chessPosition}
	cPos := pos.cPos()

	var cResults C.struct_TbRootMoves = C.struct_TbRootMoves{
		0,
		[C.TB_MAX_MOVES]C.struct_TbRootMove{},
	}
	res := C.tb_probe_root_wdl(
		cPos.white,
		cPos.black,
		cPos.kings,
		cPos.queens,
		cPos.rooks,
		cPos.bishops,
		cPos.knights,
		cPos.pawns,
		C.unsigned(cPos.rule50),
		C.unsigned(pos.cCastling()),
		C.unsigned(cPos.ep),
		cPos.turn,
		C.bool(useRule50),
		&cResults,
	)
	// fmt.Println("res", res)
	// fmt.Println("cResults", cResults)
	results := make([]TbRootMove, cResults.size)
	for i := 0; i < len(results); i++ {
		results[i] = newRootMove(cResults.moves[i])
	}

	if res == 0 {
		return results, fmt.Errorf("go-fathom: not all probes were successful")
	}
	return results, nil
}

func (f *fathom) Close() {
	C.tb_free()
}

/////////////////////////////////////////////////////

func tbPromoToChessPromo(promo int) chess.PieceType {
	switch promo {
	case C.TB_PROMOTES_NONE:
		return chess.NoPieceType
	case C.TB_PROMOTES_QUEEN:
		return chess.Queen
	case C.TB_PROMOTES_ROOK:
		return chess.Rook
	case C.TB_PROMOTES_BISHOP:
		return chess.Bishop
	case C.TB_PROMOTES_KNIGHT:
		return chess.Knight
	}
	return chess.NoPieceType
}

/////////////////////////////////////////////////////

type tbMove uint16

var _ TbMove = tbMove(0)

func (m tbMove) Move() chess.Move {
	return chess.NewMove(
		chess.Square((m>>6)&0x3F),
		chess.Square(m&0x3F),
		tbPromoToChessPromo(int((m>>12)&0x7)),
		chess.MoveTag(0),
	)
}

type tbMoveLong uint64

var _ TbMove = tbMoveLong(0)

func (ml tbMoveLong) Move() chess.Move {
	return chess.NewMove(
		chess.Square(int((ml&C.TB_RESULT_FROM_MASK)>>C.TB_RESULT_FROM_SHIFT)),
		chess.Square(int((ml&C.TB_RESULT_TO_MASK)>>C.TB_RESULT_TO_SHIFT)),
		tbPromoToChessPromo(int((ml&C.TB_RESULT_PROMOTES_MASK)>>C.TB_RESULT_PROMOTES_SHIFT)),
		chess.MoveTag(0),
	)
}

/////////////////////////////////////////////////////

type tbResult struct {
	wdl      int
	from     int
	to       int
	promotes int
	ep       int
	dtz      int
}

var _ TbResult = tbResult{}

func newTbResult(cResult C.unsigned) TbResult {
	res := tbResult{}
	v := uint64(cResult)
	res.wdl = int((v & C.TB_RESULT_WDL_MASK) >> C.TB_RESULT_WDL_SHIFT)
	res.from = int((v & C.TB_RESULT_FROM_MASK) >> C.TB_RESULT_FROM_SHIFT)
	res.to = int((v & C.TB_RESULT_TO_MASK) >> C.TB_RESULT_TO_SHIFT)
	res.promotes = int((v & C.TB_RESULT_PROMOTES_MASK) >> C.TB_RESULT_PROMOTES_SHIFT)
	res.ep = int((v & C.TB_RESULT_EP_MASK) >> C.TB_RESULT_EP_SHIFT)
	res.dtz = int((v & C.TB_RESULT_DTZ_MASK) >> C.TB_RESULT_DTZ_SHIFT)
	return res
}

func (r tbResult) Move() chess.Move {
	return chess.NewMove(
		chess.Square(r.from),
		chess.Square(r.to),
		tbPromoToChessPromo(r.promotes),
		chess.MoveTag(0),
	)
}

func (r tbResult) WDL() WDL {
	return WDL(r.wdl)
}

func (r tbResult) DTZ() int {
	return r.dtz
}

/////////////////////////////////////////////////////

type tbRootMove struct {
	move  tbMove
	pv    []tbMove
	score int
	rank  int
}

var _ TbRootMove = tbRootMove{}

func newRootMove(cResult C.struct_TbRootMove) TbRootMove {
	res := tbRootMove{}
	res.move = tbMove(cResult.move)
	res.pv = make([]tbMove, cResult.pvSize)
	for i := 0; i < len(res.pv); i++ {
		res.pv[i] = tbMove(cResult.pv[i])
	}
	res.score = int(cResult.tbScore)
	res.rank = int(cResult.tbRank)
	return res
}

func (rm tbRootMove) Move() chess.Move {
	return rm.move.Move()
}

func (rm tbRootMove) PV() []chess.Move {
	res := make([]chess.Move, len(rm.pv))
	for i, tbMove := range rm.pv {
		res[i] = tbMove.Move()
	}
	return res
}

func (rm tbRootMove) Score() int {
	return rm.score
}

func (rm tbRootMove) Rank() int {
	return rm.rank
}

/////////////////////////////////////////////////////

type position struct {
	*chess.Position
}

func (pos *position) cCastling() cCastling {
	var res cCastling
	castleRights := pos.CastleRights()
	if castleRights.CanCastle(chess.White, chess.KingSide) {
		res |= castling_K
	}
	if castleRights.CanCastle(chess.White, chess.QueenSide) {
		res |= castling_Q
	}
	if castleRights.CanCastle(chess.Black, chess.KingSide) {
		res |= castling_k
	}
	if castleRights.CanCastle(chess.Black, chess.QueenSide) {
		res |= castling_q
	}
	return res
}

func (pos *position) cPos() cPosType {
	var white, black, kings, queens, rooks, bishops, knights, pawns uint64
	var rule50 uint8
	var ep uint8
	var turn bool
	var b uint64
	board := pos.Board()
	for sq := 0; sq < 64; sq++ {
		piece := board.Piece(chess.Square(sq))
		b = 1 << sq
		switch piece {
		case chess.WhiteKing:
			white |= b
			kings |= b
		case chess.WhiteQueen:
			white |= b
			queens |= b
		case chess.WhiteRook:
			white |= b
			rooks |= b
		case chess.WhiteBishop:
			white |= b
			bishops |= b
		case chess.WhiteKnight:
			white |= b
			knights |= b
		case chess.WhitePawn:
			white |= b
			pawns |= b
		case chess.BlackKing:
			black |= b
			kings |= b
		case chess.BlackQueen:
			black |= b
			queens |= b
		case chess.BlackRook:
			black |= b
			rooks |= b
		case chess.BlackBishop:
			black |= b
			bishops |= b
		case chess.BlackKnight:
			black |= b
			knights |= b
		case chess.BlackPawn:
			black |= b
			pawns |= b
		}
	}
	rule50 = uint8(pos.HalfMoveClock())
	if temp := pos.EnPassantSquare(); temp != chess.NoSquare {
		ep = uint8(temp)
	}
	turn = pos.Turn() != chess.Black
	return cPosType{
		white:   C.uint64_t(white),
		black:   C.uint64_t(black),
		kings:   C.uint64_t(kings),
		queens:  C.uint64_t(queens),
		rooks:   C.uint64_t(rooks),
		bishops: C.uint64_t(bishops),
		knights: C.uint64_t(knights),
		pawns:   C.uint64_t(pawns),
		rule50:  C.uint8_t(rule50), // not yet supported by fathom?
		ep:      C.uint8_t(ep),
		turn:    C.bool(turn),
	}
}
