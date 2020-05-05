// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package schedule

import (
	"errors"
	"fmt"
	"time"

	"github.com/moov-io/base"

	"github.com/robfig/cron/v3"
)

// CutoffTimes is a time.Ticker which fires on banking days to trigger processing
// events (like end-of-day, or same-day ACH).
type CutoffTimes struct {
	C chan time.Time

	sched *cron.Cron
}

func ForCutoffTimes(tz string, timestamps []string) (*CutoffTimes, error) {
	ct := &CutoffTimes{
		C:     make(chan time.Time),
		sched: cron.New(),
	}
	if err := ct.registerCutoffs(tz, timestamps); err != nil {
		return nil, err
	}
	ct.sched.Start()
	return ct, nil
}

func (ct *CutoffTimes) Stop() {
	if ct == nil {
		return
	}
	if ct.C != nil {
		close(ct.C)
	}
	if ct.sched != nil {
		ct.sched.Stop()
	}
}

func (ct *CutoffTimes) maybeTick() {
	now := base.Now()
	if !now.IsWeekend() && now.IsBankingDay() {
		ct.C <- now.Time.In(time.Local)
	}
}

func (ct *CutoffTimes) registerCutoffs(tz string, timestamps []string) error {
	if len(timestamps) == 0 {
		return errors.New("missing cutoff times")
	}
	for i := range timestamps {
		if err := ct.register(tz, timestamps[i]); err != nil {
			return fmt.Errorf("timestamp=%s error=%v", timestamps[i], err)
		}
	}
	return nil
}

func (ct *CutoffTimes) register(tz string, timestamp string) error {
	when, err := time.Parse("15:04", timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' error=%v", timestamp, err)
	}

	schedule := fmt.Sprintf("CRON_TZ=%s %d %d * * *", tz, when.Minute(), when.Hour())
	ct.sched.AddFunc(schedule, func() { ct.maybeTick() })

	return nil
}
