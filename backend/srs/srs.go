package srs

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/lidchen/neuron_deck/backend/model"
)

// TODO: implement "AGAIN, HARD, NORMAL, EASY" logic
// dont use learning_steps start with 0

var learning_steps = []float32{0, 0.1667, 1.0}
var graduating_interval float32 = 24.0

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (r *RealClock) Now() time.Time { return time.Now() }

type MockClock struct {
	Current time.Time
}

func (m *MockClock) Now() time.Time { return m.Current }
func (m *MockClock) Advance(hours float64) {
	m.Current = m.Current.Add(
		time.Duration(hours * float64(time.Hour)),
	)
}

type SRSService struct {
	clock Clock
}

func NewSRSService(clock Clock) *SRSService {
	return &SRSService{clock: clock}
}

func (s *SRSService) Review(c *model.CardSrs, q int) *model.AppError {
	now := s.clock.Now()
	if c.Repetitions == 0 {
		// LEARNING PAHSE
		if err := learningPhase(q, c); err != nil {
			return err
		}
	} else {
		// REVIEW PHASE (SM-2)
		if err := sm2Phase(q, c); err != nil {
			log.Fatal(err.Message)
			return err
		}
	}

	c.LastReviewAt = now
	c.NextReviewAt = addHours(now, c.Interval)
	return nil
}

func learningPhase(q int, c *model.CardSrs) *model.AppError {
	switch q {
	case 1, 2:
		{
			// Failed: restart learning steps
			c.Interval = learning_steps[0]
		}
	case 3, 4, 5:
		{
			// Pass: advance to next learning step
			if n := getNextStep(c.Interval, learning_steps); n == nil {
				// Finish learning step, upgrade to graduating_interval
				c.Interval = graduating_interval
				c.Repetitions += 1
				c.EaseFactor = calEaseFactorUp(c.EaseFactor, q)
			} else {
				// Upgrade phase
				c.Interval = *n
			}
		}
	default:
		{
			return model.ErrInternal(fmt.Errorf("Expect q from 1-5, got %d", q))
		}
	}
	return nil
}

func sm2Phase(q int, c *model.CardSrs) *model.AppError {
	switch q {
	case 1, 2:
		{
			// Lapsed: drop back to relearning
			c.Interval = learning_steps[0]
			c.Repetitions = 0
			c.EaseFactor = calEaseFactorDown(c.EaseFactor)
		}
	case 3, 4, 5:
		{
			// Passed: grow
			c.Interval = float32(math.Round(float64(c.Interval * c.EaseFactor)))
			c.Repetitions += 1
			c.EaseFactor = calEaseFactorUp(c.EaseFactor, q)
		}
	default:
		return model.ErrBadRequest("BAD_REQUEST", fmt.Sprintf("Expect q from 1-5, got %d", q))
	}
	return nil
}

func calEaseFactorUp(current float32, q int) float32 {
	return max(1.3, current+(0.1-float32(5-q)*(0.08+float32(5-q)*0.02)))
}

func calEaseFactorDown(current float32) float32 {
	return max(1.3, current-0.2)
}

func getNextStep(current float32, steps []float32) *float32 {
	// Find the closest step to current, then advance one step.
	if len(steps) == 0 {
		return nil
	}

	closestIdx := 0
	minDiff := float32(math.Abs(float64(current - steps[0])))
	for i := 1; i < len(steps); i++ {
		diff := float32(math.Abs(float64(current - steps[i])))
		if diff < minDiff {
			minDiff = diff
			closestIdx = i
		}
	}

	nextIdx := closestIdx + 1
	if nextIdx >= len(steps) {
		return nil
	}

	return &steps[nextIdx]
}

func addHours(t time.Time, hour float32) time.Time {
	return t.Add(time.Duration(hour * float32(time.Hour)))
}
