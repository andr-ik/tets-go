package service

import "time"

type DateIntervalService struct {
	current  time.Time
	start    time.Time
	duration time.Duration
	end      time.Time
}

func NewDateIntervalService(start time.Time, duration time.Duration, end time.Time) DateIntervalService {
	return DateIntervalService{
		current:  start.Add(-duration),
		start:    start,
		duration: duration,
		end:      end,
	}
}

func (d *DateIntervalService) HasNext() bool {
	return d.current.Add(d.duration).Before(d.end) || d.current.Add(d.duration).Equal(d.end)
}

func (d *DateIntervalService) Next() time.Time {
	d.current = d.current.Add(d.duration)

	return d.current
}

func (d *DateIntervalService) Clear() {
	d.current = d.start.Add(-d.duration)
}
