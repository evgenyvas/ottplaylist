package format

import (
	"ottplaylist/types"
	"strings"
)

func M3U(ch types.Channels) string {
	var pl strings.Builder
	pl.WriteString("#EXTM3U\n")
	for _, c := range ch {
		pl.WriteString("#EXTINF:0," + c.Name + "\n")
		pl.WriteString("#EXTGRP:" + c.Category + "\n")
		pl.WriteString(c.URL + "\n")
	}
	return pl.String()
}
