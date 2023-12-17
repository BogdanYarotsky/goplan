package domain

import "time"

type TimeSlot struct {
	Id        int
	ContentId int
	Start     time.Time
	End       time.Time
}

var Slots = []*TimeSlot{
	{
		1, 777, time.Now(), time.Now(),
	},
}
