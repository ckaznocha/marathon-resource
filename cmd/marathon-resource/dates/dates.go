package dates

import (
	"errors"
	"sort"
	"time"
)

//Dates is a sortable slice of times
type Dates []time.Time

func (d Dates) Len() int           { return len(d) }
func (d Dates) Less(i, j int) bool { return d[i].Before(d[j]) }
func (d Dates) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

//NewerTimestamps returns all timestamps in a list newer than a given timestamp
func NewerTimestamps(
	timestampStrings []string,
	currentTimestampString string,
) ([]string, error) {
	var (
		currentTimestampIndex int
		timestamps            = make(Dates, len(timestampStrings))
	)
	for i, v := range timestampStrings {
		t, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return nil, err
		}
		timestamps[i] = t
	}

	if len(currentTimestampString) == 0 {
		currentTimestampString = time.Unix(0, 0).Format(time.RFC3339Nano)
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
