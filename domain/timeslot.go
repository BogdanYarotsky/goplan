package domain

import (
	"time"
)

type Id int
type SlotId Id
type ContentId Id
type MachineId Id

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (id MachineId) Validate() error {
	if id < 1 {
		return ValidationError{"machine id must be above 0"}
	}
	return nil
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (tr TimeRange) Validate() error {
	if tr.Start.After(tr.End) {
		return ValidationError{"start must be before end"}
	}

	return nil
}

type TimeSlot struct {
	Id        SlotId
	MachineId MachineId
	ContentId ContentId
	TimeRange TimeRange
}

var Slots = []TimeSlot{
	{
		1, 666, 777, TimeRange{time.Now(), time.Now()},
	},
	{
		2, 555, 888, TimeRange{time.Now(), time.Now()},
	},
}

type PlanRepository interface {
}

type PlanService struct {
}

func NewPlanService() *PlanService {
	return &PlanService{}
}

func (s *PlanService) PlanSlot(mid MachineId, cid ContentId, tr TimeRange) (TimeSlot, error) {
	return TimeSlot{}, nil
}

func (s *PlanService) GetSlots(mid MachineId, tr TimeRange) ([]TimeSlot, error) {
	if err := mid.Validate(); err != nil {
		return nil, err
	}

	if err := tr.Validate(); err != nil {
		return nil, err
	}

	return Slots, nil
}
