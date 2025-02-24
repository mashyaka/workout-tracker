package database

import (
	"time"
)

type StatisticsItem struct {
	Unit          string
	Counter       int
	Distance      float64
	TotalDistance float64
	Duration      time.Duration
	TotalDuration time.Duration
	Speed         float64
	FirstPoint    *MapPoint
	LastPoint     *MapPoint
	IsBest        bool
	IsWorst       bool
}

func (si *StatisticsItem) SpeedKPH() float64 {
	return 3.6 * si.Speed
}

func (si *StatisticsItem) createNext(fp *MapPoint) StatisticsItem {
	return StatisticsItem{
		Unit:          si.Unit,
		Counter:       si.Counter + 1,
		TotalDistance: si.TotalDistance,
		TotalDuration: si.TotalDuration,
		FirstPoint:    fp,
	}
}

func (si *StatisticsItem) canHave(unit string, fp *MapPoint) bool {
	switch unit {
	case "km":
		return si.canHaveDistance(fp.Distance)
	case "min":
		return si.canHaveDuration(fp.Duration)
	}

	return true
}

func (si *StatisticsItem) canHaveDistance(distance float64) bool {
	return int(si.TotalDistance+distance) < si.Counter*1000
}

func (si *StatisticsItem) canHaveDuration(duration time.Duration) bool {
	return si.Duration+duration < time.Minute
}

func (si *StatisticsItem) CalcultateSpeed() {
	si.Speed = si.Distance / si.Duration.Seconds()
}

func calculateBestAndWorst(items []StatisticsItem) {
	if len(items) == 0 {
		return
	}

	worst := 0
	best := 0

	for i := range items {
		if items[i].Speed < items[worst].Speed {
			worst = i
		}

		if items[i].Speed > items[best].Speed {
			best = i
		}
	}

	items[worst].IsWorst = true
	items[best].IsBest = true
}

func (w *Workout) StatisticsPerKilometer() []StatisticsItem {
	return w.statisticsWithUnit("km")
}

func (w *Workout) StatisticsPerMinute() []StatisticsItem {
	return w.statisticsWithUnit("min")
}

func (w *Workout) statisticsWithUnit(unit string) []StatisticsItem {
	var items []StatisticsItem

	nextItem := StatisticsItem{
		Unit:       unit,
		Counter:    1,
		FirstPoint: &w.Data.Points[0],
	}

	for i, p := range w.Data.Points {
		if !nextItem.canHave(unit, &w.Data.Points[i]) {
			nextItem.LastPoint = &w.Data.Points[i]
			nextItem.CalcultateSpeed()
			items = append(items, nextItem)
			nextItem = nextItem.createNext(&w.Data.Points[i])
		}

		nextItem.Distance += p.Distance
		nextItem.TotalDistance += p.Distance

		// m/s -> km/h, cut-off is speed less than 1 km/h
		if p.AverageSpeed()*3.6 >= 1.0 {
			nextItem.Duration += p.Duration
			nextItem.TotalDuration += p.Duration
		}
	}

	nextItem.LastPoint = &w.Data.Points[len(w.Data.Points)-1]

	if nextItem.FirstPoint != nil {
		nextItem.CalcultateSpeed()
		items = append(items, nextItem)
	}

	calculateBestAndWorst(items)

	return items
}
