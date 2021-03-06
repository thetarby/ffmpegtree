package ffmpegtree

import (
	"fmt"
	"time"
)

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

var x = 0

func randStr() string {
	x++
	return fmt.Sprintf("var_%v", x)
}
