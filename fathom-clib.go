package fathom

// #cgo CFLAGS: -std=c99 -Imodule-fathom/src
// #include <module-fathom/src/tbprobe.c>
import "C"
import (
	// "unsafe"
	"fmt"
	
	"github.com/notnil/chess"
)

type cPos C.Pos

type cCastling int

const (
	castlingUnknown = cCastling(0)
	castling_K      = cCastling(C.TB_CASTLING_K)
	castling_Q      = cCastling(C.TB_CASTLING_Q)
	castling_k      = cCastling(C.TB_CASTLING_k)
	castling_q      = cCastling(C.TB_CASTLING_q)
)

// Win-Draw-Lose result status.
type WDL int

const (
	WDLUnknown  = WDL(-1)
	Loss        = WDL(C.TB_LOSS)
	BlessedLoss = WDL(C.TB_BLESSED_LOSS)
	Draw        = WDL(C.TB_DRAW)
	CursedWin   = WDL(C.TB_CURSED_WIN)
	Win         = WDL(C.TB_WIN)
)

func (this WDL) String() string {
	switch this {
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

type Fathom interface {
	ProbeWdl(chessPosition *chess.Position) (WDL, error)
	ProbeRoot(chessPosition *chess.Position) ([]TbResult, error)
	ProbeRootDtz(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error)
	ProbeRootWdl(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error)
}

type fathom struct {
	tbDir string
}

func NewFathom(tbDir string) (Fathom, error) {
	cTbDir := C.CString(tbDir)
	res := C.tb_init(cTbDir)
	if res != true {
		return nil, fmt.Errorf("go-fathom: unable to init tablebases")
	}
	return &fathom{
		tbDir: tbDir,
	}, nil
}

func (this *fathom) ProbeWdl(chessPosition *chess.Position) (WDL, error) {
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
	// fmt.Println("res", res)
	if res == C.TB_RESULT_FAILED {
		return WDLUnknown, fmt.Errorf("go-fathom: tb_probe_wdl failed")
	}
	return WDL(res), nil
}

func (this *fathom) ProbeRoot(chessPosition *chess.Position) ([]TbResult, error) {
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
	// fmt.Println("res", res)
	// fmt.Println("cResults", cResults)
	if res == C.TB_RESULT_FAILED {
		return nil, fmt.Errorf("go-fathom: probe_root failed")
	}
	results := make([]TbResult, 0, len(cResults))
	for _, cResult := range cResults {
		if cResult == C.TB_RESULT_FAILED {
			break
		}
		results = append(results, newTbResult(cResult))
	}
	return results, nil
}

func (this *fathom) ProbeRootDtz(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error) {
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

func (this *fathom) ProbeRootWdl(chessPosition *chess.Position, useRule50 bool) ([]TbRootMove, error) {
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
	fmt.Println("res", res)
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

func (this *fathom) Close() {
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

type TbMove interface {
	Move() chess.Move
}

type tbMove uint16

func (this tbMove) Move() chess.Move {
	return chess.NewMove(
		chess.Square((this >> 6) & 0x3F),
		chess.Square(this & 0x3F),
		tbPromoToChessPromo(int((this >> 12) & 0x7)),
		chess.MoveTag(0),
	)
}

/////////////////////////////////////////////////////

type TbResult interface {
	TbMove
	WDL() WDL
	DTZ() int
}

type tbResult struct {
	wdl  int
	from int
	to   int
	promotes int
	ep       int
	dtz      int
}

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

func (this tbResult) Move() chess.Move {
	return chess.NewMove(
		chess.Square(this.from),
		chess.Square(this.to),
		tbPromoToChessPromo(this.promotes),
		chess.MoveTag(0),
	)
}

func (this tbResult) WDL() WDL {
	return WDL(this.wdl)
}

func (this tbResult) DTZ() int {
	return this.dtz
}

/////////////////////////////////////////////////////

type TbRootMove interface {
	TbMove
	// Principal Variation
	PV() []chess.Move
	Score() int
	Rank() int
}

type tbRootMove struct {
	move tbMove
	pv []tbMove
	score int
	rank int
}

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

func (this tbRootMove) Move() chess.Move {
	return this.move.Move()
}

func (this tbRootMove) PV() []chess.Move {
	res := make([]chess.Move, len(this.pv))
	for i, tbMove := range this.pv {
		res[i] = tbMove.Move()
	}
	return res
}

func (this tbRootMove) Score() int {
	return this.score
}

func (this tbRootMove) Rank() int {
	return this.rank
}

/////////////////////////////////////////////////////

type position struct {
	*chess.Position
}

func (this *position) cCastling() cCastling {
	var res cCastling
	castleRights := this.CastleRights()
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

func (this *position) cPos() cPos {
	var white, black, kings, queens, rooks, bishops, knights, pawns uint64
	var rule50 uint8
	var ep uint8
	var turn bool
	var b uint64
	board := this.Board()
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
	rule50 = uint8(this.HalfMoveClock())
	if temp := this.EnPassantSquare(); temp != chess.NoSquare {
		ep = uint8(temp)
	}
	turn = this.Turn() != chess.Black
	return cPos{
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
