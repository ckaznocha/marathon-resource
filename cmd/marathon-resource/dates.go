package main

import (
	"errors"
	"sort"
	"time"
)

type dates []time.Time

func (d dates) Len() int           { return len(d) }
func (d dates) Less(i, j int) bool { return d[i].Before(d[j]) }
func (d dates) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func newerTimestamps(
	timestampStrings []string,
	currentTimestampString string,
) ([]string, error) {
	var (
		currentTimestampIndex int
		timestamps            = make(dates, len(timestampStrings))
	)
	for i, v := range timestampStrings {
		t, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return nil, err
		}
		timestamps[i] = t
	}

	currentTimestamp, err := time.Parse(
		time.RFC3339Nano,
		currentTimestampString,
	)
	if err != nil {
		return nil, err
	}

	if !sort.IsSorted(timestamps) {
		sort.Sort(timestamps)
	}

	currentTimestampIndex = sort.Search(
		len(timestamps),
		func(i int) bool {
			return timestamps[i].After(currentTimestamp) ||
				timestamps[i].Equal(currentTimestamp)
		},
	)

	if currentTimestampIndex >= len(timestamps) {
		return nil, errors.New("Version not found")
	}

	var newerTimestamps []string
	for _, v := range timestamps[currentTimestampIndex:] {
		newerTimestamps = append(newerTimestamps, v.Format(time.RFC3339Nano))
	}

	return newerTimestamps, nil

}
