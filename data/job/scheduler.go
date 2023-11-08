package job

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ericmarcelinotju/gram/constant/enums"
	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	interval       enums.SchedulerInterval
	monthlyDate    int
	dailyHour      int
	dailyMinute    int
	minutelyMinute int
	scheduler      *gocron.Scheduler
	scheduleFunc   interface{}
}

func NewScheduler(hour, minute int) (*Scheduler, error) {

	s := &Scheduler{}

	err := s.SetDaily(hour, minute)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Scheduler) SetMinutely(minute int) error {
	s.interval = enums.SchedulerIntervalMinutely
	s.minutelyMinute = minute

	if s.IsRunning() {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) SetDaily(hour, minute int) error {
	s.interval = enums.SchedulerIntervalDaily
	s.dailyHour = hour
	s.dailyMinute = minute

	if s.IsRunning() {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) SetMonthly(monthlyDate int) error {
	s.interval = enums.SchedulerIntervalMonthly
	s.monthlyDate = monthlyDate

	if s.IsRunning() {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) SetScheduleFunc(scheduleFunc interface{}) error {
	typ := reflect.TypeOf(scheduleFunc)
	if typ.Kind() != reflect.Func {
		return errors.New("not a function")
	}

	s.scheduleFunc = scheduleFunc

	if s.IsRunning() {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) Start() error {
	s.Stop()

	scheduler := gocron.NewScheduler(time.Local)
	s.scheduler = scheduler

	if s.scheduleFunc == nil {
		return errors.New("schedule function not set yet")
	}
	if s.interval == enums.DefaultSchedulerInterval {
		return errors.New("schedule time not set yet")
	}

	if s.interval == enums.SchedulerIntervalMonthly {
		_, err := s.scheduler.Every(1).Month(s.monthlyDate).At("00:00").Do(s.scheduleFunc)
		if err != nil {
			return err
		}
	}
	if s.interval == enums.SchedulerIntervalDaily {
		timeStr := fmt.Sprintf("%02d:%02d", s.dailyHour, s.dailyMinute)
		_, err := s.scheduler.Every(1).Day().At(timeStr).Do(s.scheduleFunc)
		if err != nil {
			return err
		}
	}
	if s.interval == enums.SchedulerIntervalMinutely {
		timeStr := fmt.Sprintf("%dm", s.minutelyMinute)
		_, err := s.scheduler.Every(timeStr).Do(s.scheduleFunc)
		if err != nil {
			return err
		}
	}
	s.scheduler.StartAsync()
	return nil
}

func (s *Scheduler) Stop() {
	if s.scheduler != nil {
		s.scheduler.Stop()
	}
}

func (s *Scheduler) IsRunning() bool {
	return s.scheduler != nil && s.scheduler.IsRunning()
}
