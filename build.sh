#!/usr/bin/env bash

TARGET="go-chan-scraper"

if [ -f "$TARGET" ] ; then
    rm "$TARGET"
fi

go build
