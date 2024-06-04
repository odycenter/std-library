package dates_test

import (
	"github.com/stretchr/testify/assert"
	dates "std-library/date"
	"testing"
	"time"
)

func TestStartOfTheDay(t *testing.T) {
	d := time.Date(2023, 10, 12, 12, 34, 22, 123, time.Local)
	md := dates.Month(d)

	assert.Equal(t, time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local), md.FirstDay().StartOfTheDay())
	assert.Equal(t, time.Date(2023, 10, 31, 0, 0, 0, 0, time.Local), md.LastDay().StartOfTheDay())

	md = dates.LastMonth(d)
	assert.Equal(t, time.Date(2023, 9, 1, 0, 0, 0, 0, time.Local), md.FirstDay().StartOfTheDay())
	assert.Equal(t, time.Date(2023, 9, 30, 0, 0, 0, 0, time.Local), md.LastDay().StartOfTheDay())

	d = time.Date(2024, 1, 12, 12, 34, 22, 123, time.Local)
	md = dates.LastMonth(d)
	assert.Equal(t, time.Date(2023, 12, 1, 0, 0, 0, 0, time.Local), md.FirstDay().StartOfTheDay())
	assert.Equal(t, time.Date(2023, 12, 31, 0, 0, 0, 0, time.Local), md.LastDay().StartOfTheDay())

	md = dates.NextMonth(d)
	assert.Equal(t, time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local), md.FirstDay().StartOfTheDay())
	assert.Equal(t, time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local), md.LastDay().StartOfTheDay())

}

func TestEndOfTheDay(t *testing.T) {
	d := time.Date(2023, 2, 12, 12, 34, 56, 789, time.Local)
	md := dates.Month(d)
	assert.Equal(t, time.Date(2023, 2, 1, 23, 59, 59, 999999999, time.Local), md.FirstDay().EndOfTheDay())
	assert.Equal(t, time.Date(2023, 2, 28, 23, 59, 59, 999999999, time.Local), md.LastDay().EndOfTheDay())

	md = dates.LastMonth(d)
	assert.Equal(t, time.Date(2023, 1, 1, 23, 59, 59, 999999999, time.Local), md.FirstDay().EndOfTheDay())
	assert.Equal(t, time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.Local), md.LastDay().EndOfTheDay())

	d = time.Date(2024, 1, 12, 12, 34, 22, 123, time.Local)
	md = dates.LastMonth(d)
	assert.Equal(t, time.Date(2023, 12, 1, 23, 59, 59, 999999999, time.Local), md.FirstDay().EndOfTheDay())
	assert.Equal(t, time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.Local), md.LastDay().EndOfTheDay())

	md = dates.NextMonth(d)
	assert.Equal(t, time.Date(2024, 2, 1, 23, 59, 59, 999999999, time.Local), md.FirstDay().EndOfTheDay())
	assert.Equal(t, time.Date(2024, 2, 29, 23, 59, 59, 999999999, time.Local), md.LastDay().EndOfTheDay())
}

func TestDateTime(t *testing.T) {
	d := time.Date(2023, 10, 12, 12, 34, 56, 123, time.Local)
	md := dates.Month(d)
	assert.Equal(t, time.Date(2023, 10, 1, 12, 34, 56, 123, time.Local), md.FirstDay().DateTime())
	assert.Equal(t, time.Date(2023, 10, 31, 12, 34, 56, 123, time.Local), md.LastDay().DateTime())

	md = dates.LastMonth(d)
	assert.Equal(t, time.Date(2023, 9, 1, 12, 34, 56, 123, time.Local), md.FirstDay().DateTime())
	assert.Equal(t, time.Date(2023, 9, 30, 12, 34, 56, 123, time.Local), md.LastDay().DateTime())

	d = time.Date(2024, 1, 12, 12, 34, 60, 123, time.Local)
	md = dates.LastMonth(d)
	assert.Equal(t, time.Date(2023, 12, 1, 12, 35, 0, 123, time.Local), md.FirstDay().DateTime())
	assert.Equal(t, time.Date(2023, 12, 31, 12, 35, 0, 123, time.Local), md.LastDay().DateTime())

	md = dates.NextMonth(d)
	assert.Equal(t, time.Date(2024, 2, 29, 12, 35, 0, 123, time.Local), md.DateTime())
	assert.Equal(t, time.Date(2024, 2, 1, 12, 35, 0, 123, time.Local), md.FirstDay().DateTime())
	assert.Equal(t, time.Date(2024, 2, 29, 12, 35, 0, 123, time.Local), md.LastDay().DateTime())
	assert.Equal(t, md.LastDay().DateTime(), md.DateTime())
}
