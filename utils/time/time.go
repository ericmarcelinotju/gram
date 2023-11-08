package time

import "time"

const DEFAULT_RESPONSE_TIME_FORMAT = "02 January 2006 15:04:05"

func FormatResponseTime(input *time.Time) *string {
	if input != nil {
		res := input.Format(DEFAULT_RESPONSE_TIME_FORMAT)
		return &res
	}
	return nil
}

func ParseTime(layout, value string) (result time.Time, err error) {
	result, err = time.Parse(layout, value)
	if err != nil {
		return
	}
	result = result.Add(-(7 * time.Hour))
	return
}
