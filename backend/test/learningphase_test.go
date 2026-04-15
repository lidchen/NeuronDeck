package test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lidchen/neuron_deck/backend/model"
	"github.com/lidchen/neuron_deck/backend/srs"
)

func TestLearningPhase(t *testing.T) {
	type step struct {
		advanceHours float32
		q            int
	}
	type want struct {
		repetitions int
		interval    float32
		easeFactor  float32
	}

	tests := []struct {
		name  string
		steps []step
		want  want
	}{
		{
			name: "happy path through learning",
			steps: []step{
				{0, 4},      // new card, pass
				{0.1667, 4}, // pass 10min step
				{1.0, 4},    // pass 1hr step → graduate
			},
			want: want{
				repetitions: 1,
				interval:    24.0,
			},
		},
		{
			name: "fail multiple times then pass",
			steps: []step{
				{0, 1},      // fail → immediate
				{0, 1},      // fail → immediate
				{0, 3},      // pass → 10 min
				{0.1667, 4}, // pass → 1 hour
				{1.0, 4},    // graduate
			},
			want: want{
				repetitions: 1,
				interval:    24.0,
			},
		},
		{
			name: "fail never advances step index",
			steps: []step{
				{0, 4}, // advance to step 1
				{0, 1}, // fail at step 1
			},
			want: want{
				repetitions: 0,
				interval:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := &srs.MockClock{
				Current: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			}
			svc := srs.NewSRSService(clock)

			cSrs := model.CardSrs{
				CardId: 1, Interval: 0.0167, EaseFactor: 2.5, Repetitions: 0,
				NextReviewAt: clock.Now(), LastReviewAt: clock.Now(),
			}

			for _, s := range tt.steps {
				clock.Advance(float64(s.advanceHours))
				err := svc.Review(&cSrs, s.q)
				if err != nil {
					log.Fatal(err.Message)
				}
			}

			assert.Equal(t, tt.want.repetitions, cSrs.Repetitions)
			assert.InDelta(t, tt.want.interval, cSrs.Interval, 0.01)
		})
	}
}
