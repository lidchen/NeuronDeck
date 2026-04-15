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
	current time.Time
}

func (m *MockClock) Now() time.Time { return m.current }
func (m *MockClock) Advance(hours float64) {
	m.current = m.current.Add(
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
			if n := getNextStep(c.Interval, learning_steps); n != nil {
				if *n == -1.0 {
					// Not found
					return model.ErrInternal(fmt.Errorf("current interval: %f is not found in learning step", c.Interval))
				}
				// Finish learning step, upgrade to graduating_interval
				c.Interval = graduating_interval
				c.Repetitions += 1
				c.EaseFactor = calEaseFactorUp(c.EaseFactor, q)
			} else {
				// Upgrade to review phase
				c.Interval = graduating_interval
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
	// assume current could be found in steps
	length := len(steps)
	for i, s := range steps {
		if current == s {
			i++
			if i > length {
				return nil
			} else {
				return &steps[i]
			}
		}
	}
	v := float32(-1)
	return &v
}

func addHours(t time.Time, hour float32) time.Time {
	return t.Add(time.Duration(hour * float32(time.Hour)))
}
