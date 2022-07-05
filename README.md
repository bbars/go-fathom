# go-fathom
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/bbars/go-fathom)

Syzygy tablebase probe implementation for Go.

This package wraps a fork of an original Fathom tool written in C: https://github.com/jdart1/Fathom.

Also it uses well-known external package https://github.com/notnil/chess as a public communication layer.

Interface methods allow to probe DTZ and WDL tables (\*.rtbw and \*.rtbz Syzygy tablebases).

## Tablebase files

If you don't have tablebase files yet, follow [this link](https://github.com/niklasf/shakmaty-syzygy/tree/master/tables) and explore suggested TXT. These TXT files contain URLs to tablebase files you should download.

## Example: play `chess.Game`

We are going to call `fathom.ProbeRoot(...)` and apply the suggested move in a loop until the game is over.

```golang
package main

import (
	"fmt"

	gofathom "github.com/bbars/go-fathom"
	"github.com/notnil/chess"
)

func main() {
	fen := "4k3/7q/8/8/8/8/B61/4K3 b - - 0 1"
	
	fathom, err := gofathom.NewFathom("./path/to/tables/chess/")
	if err != nil {
		panic(err)
	}
	
	fenOption, err := chess.FEN(fen)
	if err != nil {
		panic(err)
	}
	game := chess.NewGame(fenOption)
	game.AddTagPair("FEN", fen)
	fmt.Printf("FEN: %q\n\n", game.Position())
	
	// generate moves until game is over
	for game.Outcome() == chess.NoOutcome {
		tbMove, _, err := fathom.ProbeRoot(game.Position())
		if err != nil {
			panic(err)
		}
		// select a best move
		move := tbMove.Move()
		fmt.Printf("%s move %s -> ", game.Position().Turn(), &move)
		err = game.Move(&move)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", game.Position())
	}
	// print outcome and game PGN
	fmt.Println(game.Position().Board().Draw())
	fmt.Printf("Game completed. %s by %s.\n\n", game.Outcome(), game.Method())
	fmt.Println("PGN:\n")
	fmt.Println(game.String())
}
```

## Example: call all available funcs

Here we are going to invoke all available funcs to analyze the state.

```golang
package main

import (
	"fmt"

	gofathom "github.com/bbars/go-fathom"
	"github.com/notnil/chess"
)

func main() {
	fen := "4k3/7q/8/8/8/8/B61/4K3 b - - 0 1"
	
	fathom, err := gofathom.NewFathom("./path/to/tables/chess/")
	if err != nil {
		panic(err)
	}
	
	pos := &chess.Position{}
	err = pos.UnmarshalText([]byte(fen))
	if err != nil {
		panic(err)
	}
	
	fmt.Println("# fathom.ProbeWDL")
	res1, _ := fathom.ProbeWDL(pos)
	fmt.Printf("%#v = %s\n", res1, res1)
	
	fmt.Println("\n# fathom.ProbeRoot")
	tbMove, res2, _ := fathom.ProbeRoot(pos)
	move := tbMove.Move()
	fmt.Printf("move: %#v\tstring: %q\n", move, &move)
	fmt.Println("move#\tWDL\tDTZ")
	for _, tbRes := range res2 {
		move := tbRes.Move()
		fmt.Printf("%s\t%s\t%v\n", move.String(), tbRes.WDL(), tbRes.DTZ())
	}
	
	fmt.Println("\n# fathom.ProbeRootDTZ")
	res3, _ := fathom.ProbeRootDTZ(pos, true)
	fmt.Println("move#\tRank\tScore\tPV")
	for _, tbRes := range res3 {
		move := tbRes.Move()
		fmt.Printf("%s\t%v\t%v\t%v\n", move.String(), tbRes.Rank(), tbRes.Score(), tbRes.PV())
	}
	
	fmt.Println("\n# fathom.ProbeRootWDL")
	fmt.Println("move#\tRank\tScore\tPV")
	res4, _ := fathom.ProbeRootWDL(pos, true)
	for _, tbRes := range res4 {
		move := tbRes.Move()
		fmt.Printf("%s\t%v\t%v\t%v\n", move.String(), tbRes.Rank(), tbRes.Score(), tbRes.PV())
	}
}
```
