package dates

import "time"

type MonthTime struct {
	t time.Time
}

func (md *MonthTime) FirstDay() *DayTime {
	return &DayTime{md.t.AddDate(0, 0, -md.t.Day()+1)}
}

func (md *MonthTime) LastDay() *DayTime {
	return &DayTime{md.t}
}

func (md *MonthTime) DateTime() time.Time {
	return md.t
}

type DayTime struct {
	t time.Time
}

func (dt *DayTime) StartOfTheDay() time.Time {
	return StartOfTheDay(dt.t)
}

func (dt *DayTime) EndOfTheDay() time.Time {
	return EndOfTheDay(dt.t)
}

func (dt *DayTime) DateTime() time.Time {
	return dt.t
}

func ShiftMonth(shift int, dateTime ...time.Time) *MonthTime {
	var t time.Time
	if len(dateTime) == 0 {
		t = time.Now()
	} else {
		t = dateTime[0]
	}
	return &MonthTime{t.AddDate(0, shift+1, -t.Day())}
}

func Month(in ...time.Time) *MonthTime {
	return ShiftMonth(0, in...)
}

func LastMonth(in ...time.Time) *MonthTime {
	return ShiftMonth(-1, in...)
}

func NextMonth(in ...time.Time) *MonthTime {
	return ShiftMonth(1, in...)
}

func StartOfTheDay(dateTime time.Time) time.Time {
	return time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, dateTime.Location())
}

func EndOfTheDay(dateTime time.Time) time.Time {
	return time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 23, 59, 59, 999999999, dateTime.Location())
}
