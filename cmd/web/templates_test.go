// testing humanDate() from the template file
package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2004, 11, 21, 10, 0, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "21 Nov 2004 at 10:00" {
		t.Errorf("want %q; got %q", "21 Nov 2004 at 10:00", hd)
	}
}
