package time

import (
	"fmt"
	"time"
)

const (
	Minute   = 60
	Hour     = 60 * Minute
	Day      = 24 * Hour
	Week     = 7 * Day
	Month    = 30 * Day
	Year     = 12 * Month
	LongTime = 37 * Year
)

func FormatResponseDiffDuration(diff time.Duration) string {
	diff /= time.Second

	lbl := "ago"

	after := diff > 0

	if after {
		lbl = "from now"
		diff += 1
	} else {
		diff *= -1
	}

	switch {
	case diff <= 0:
		return "now"
	case diff <= 2:
		return fmt.Sprintf("1 second %s", lbl)
	case diff < 1*Minute:
		return fmt.Sprintf("%d seconds %s", diff, lbl)

	case diff < 2*Minute:
		return fmt.Sprintf("1 minute %s", lbl)
	case diff < 1*Hour:
		return fmt.Sprintf("%d minutes %s", diff/Minute, lbl)

	case diff < 2*Hour:
		return fmt.Sprintf("1 hour %s", lbl)
	case diff < 1*Day:
		return fmt.Sprintf("%d hours %s", diff/Hour, lbl)

	case diff < 2*Day:
		return fmt.Sprintf("1 day %s", lbl)
	case diff < 1*Week:
		return fmt.Sprintf("%d days %s", diff/Day, lbl)

	case diff < 2*Week:
		return fmt.Sprintf("1 week %s", lbl)
	case diff < 1*Month:
		return fmt.Sprintf("%d weeks %s", diff/Week, lbl)

	case diff < 2*Month:
		return fmt.Sprintf("1 month %s", lbl)
	case diff < 1*Year:
		return fmt.Sprintf("%d months %s", diff/Month, lbl)

	case diff < 18*Month:
		return fmt.Sprintf("1 year %s", lbl)
	case diff < 2*Year:
		return fmt.Sprintf("2 years %s", lbl)
	case diff < LongTime:
		return fmt.Sprintf("%d years %s", diff/Year, lbl)
	}

	if after {
		return "a while from now"
	}
	return "long ago"
}

func FormatResponseDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
