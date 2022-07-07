#!/bin/sh

git submodule update --recursive --init
2>/dev/null rm -rf fathom-copy/*
cp -R fathom-module/* fathom-copy/
