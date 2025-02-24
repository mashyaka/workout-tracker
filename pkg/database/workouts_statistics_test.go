package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStatisticsItem_CalcultateSpeed(t *testing.T) {
	si := &StatisticsItem{
		Distance: 1000,
		Duration: 60 * time.Second,
	}

	assert.Zero(t, si.Speed)
	si.CalcultateSpeed()

	assert.InDelta(t, 16.67, si.Speed, 0.01)
}

func TestStatisticsItem_CanHave(t *testing.T) {
	si := &StatisticsItem{
		TotalDistance: 1800,
		Counter:       2,
	}

	assert.True(t, si.canHaveDistance(100))
	assert.True(t, si.canHaveDistance(199))
	assert.False(t, si.canHaveDistance(200))
	assert.False(t, si.canHaveDistance(203))
}

func TestStatisticsItem_StatisticsPerKilometer(t *testing.T) {
	w := defaultWorkout(t)
	assert.NotNil(t, w)

	stats := w.StatisticsPerKilometer()
	assert.Len(t, stats, 4)

	assert.Equal(t, 1, stats[0].Counter)
	assert.Equal(t, 2, stats[1].Counter)
	assert.Equal(t, 4, stats[3].Counter)

	assert.InDelta(t, 1000, stats[0].TotalDistance, 20)
	assert.InDelta(t, 2000, stats[1].TotalDistance, 20)
	assert.InDelta(t, 3100, stats[3].TotalDistance, 20)

	assert.Equal(t, 242*time.Second, stats[0].Duration)
	assert.Equal(t, 336*time.Second, stats[1].Duration)
	assert.Equal(t, 31*time.Second, stats[3].Duration)
}

func TestStatisticsItem_StatisticsPerMinute(t *testing.T) {
	w := defaultWorkout(t)
	assert.NotNil(t, w)

	stats := w.StatisticsPerMinute()
	assert.Len(t, stats, 16)

	assert.Equal(t, 1, stats[0].Counter)
	assert.Equal(t, 3, stats[2].Counter)
	assert.Equal(t, 8, stats[7].Counter)
	assert.Equal(t, 16, stats[15].Counter)

	assert.InDelta(t, 180, stats[0].TotalDistance, 20)
	assert.InDelta(t, 690, stats[2].TotalDistance, 20)
	assert.InDelta(t, 1765, stats[7].TotalDistance, 20)
	assert.InDelta(t, 3100, stats[15].TotalDistance, 20)

	assert.Equal(t, 56*time.Second, stats[0].Duration)
	assert.Equal(t, 57*time.Second, stats[1].Duration)
	assert.Equal(t, 106*time.Second, stats[7].Duration)
	assert.Equal(t, 31*time.Second, stats[15].Duration)

	assert.Equal(t, 56*time.Second, stats[0].TotalDuration)
	assert.Equal(t, 113*time.Second, stats[1].TotalDuration)
	assert.Equal(t, 509*time.Second, stats[7].TotalDuration)
	assert.Equal(t, 948*time.Second, stats[15].TotalDuration)
}
