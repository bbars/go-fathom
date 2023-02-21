package fathom

import "errors"

// ErrNoTablebases no tablebase files are found within specified directory
var ErrNoTablebases = errors.New("go-fathom: no tablebase files are found")

// ErrCheckmate unable to dig, because the game is over: checkmate
var ErrCheckmate = errors.New("go-fathom: checkmate")

// ErrStalemate unable to dig, because the game is over: stalemate
var ErrStalemate = errors.New("go-fathom: stalemate")
