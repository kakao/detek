package detek

import (
	"context"
	"fmt"
	"time"

	"github.com/kakao/detek/pkg/log"
	"github.com/pkg/errors"
)

type Manager struct {
	Detector  []Detector
	Collector []Collector
	store     *Store
}

func NewManager(collector []Collector, detector []Detector) *Manager {
	kv := make(map[string]Stored)
	return &Manager{
		store: &Store{
			kv: kv,
		},
		Collector: collector,
		Detector:  detector,
	}
}

func (m *Manager) ShowPlan() ([]Collector, []Detector) {
	return m.Collector, m.Detector
}

// TODO(@scotty.scott): Labels, etc....
type MangerRunOptions struct{}

/*
Work Flow (for now)
1. Do Collector Things Synchronously (for now)
2. Check Any Error is Returned (if error, creating reports in here)
3. Do Detector Things Synchronously (for now)
4. Aggregate Reports
5. Return
*/
func (m *Manager) Run(ctx context.Context, opts *MangerRunOptions) (*ReportList, error) {
	log.Info(ctx, "Starting Collector....")
	result := ReportList{
		StartedAt: time.Now(),
	}

	// Collecting
	type CollectingProblem struct {
		ID    string `json:"collector_id"`
		Error string `json:"fail_reason"`
	}
	problems := []CollectingProblem{}
	for _, p := range m.Collector {
		// Preparing
		producer := p
		meta := producer.GetMeta()
		dctx, _, err := newDetekContext(ctx, meta.ID, m.store, detekConfigOpts{
			ConsumingPlan: meta.Required,
			ProducingPlan: meta.Producing,
			Meta:          p.GetMeta().MetaInfo,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "fail to generate detek context for %q", meta.ID)
		}

		// Validating
		for k, v := range meta.Required {
			var val interface{}
			val, _, err = m.store.Get(k)
			if err != nil {
				err = errors.Wrap(err, k)
				break
			} else if v.Type.Kind() != TypeOf(val).Kind() {
				err = fmt.Errorf("expect %q type for key %q, but got %q", v.Type, k, TypeOf(val))
				err = errors.Wrap(err, k)
				break
			}
		}
		if err != nil {
			// Wrapping dependency error
			err = errors.Wrap(err, "will not run this collector, since required data is not provided")
		} else {
			// Run
			err = func() (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panic occured %v", r)
					}
					if err != nil {
						log.Error(dctx.Context(), "Error: %v", err)
					}
					log.Info(dctx.Context(), "done")
				}()
				err = producer.Do(*dctx)
				return
			}()
		}
		if err != nil {
			problems = append(problems, CollectingProblem{
				ID:    producer.GetMeta().ID,
				Error: fmt.Sprintf("%v", err),
			})
		}
	}
	log.Info(ctx, "All Collector are doing there jobs well")
	collectingReport := Report{
		CreatedAt: time.Now(),
		MetaInfo: MetaInfo{
			ID:          "collector_reports",
			Description: "failed Collectors will represented in here",
			Labels:      []string{"detek", "collector", "data"},
		},
		Level:        Normal,
		CurrentState: NormalStatus,
		ReportSpec:   ReportSpec{HasPassed: true},
	}

	log.Info(ctx, "Starting detectors.....")
	reports := []Report{}
	if len(problems) != 0 {
		collectingReport.Level = Unknown
		collectingReport.Problem = JSONableData{
			Description: "list of failed collectors",
			Data:        problems,
		}
		collectingReport.CurrentState = Description{
			Explanation: "some of Collectors are failed",
			Solution:    "check returned error from Collectors",
		}
		reports = append(reports, collectingReport)
	}

	// Detecting
	for _, c := range m.Detector {
		// Preparing
		consumer := c
		meta := consumer.GetMeta()
		dctx, _, err := newDetekContext(ctx, meta.ID, m.store, detekConfigOpts{
			ConsumingPlan: meta.Required,
			Meta:          meta.MetaInfo,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "fail to generate detek context for %q", meta.ID)
		}

		// Validating
		var report *Report
		for k, v := range meta.Required {
			var val interface{}
			val, _, err = m.store.Get(k)
			if err != nil {
				err = errors.Wrap(err, k)
				break
			} else if v.Type.Kind() != TypeOf(val).Kind() {
				err = fmt.Errorf("expect %q type for key %q, but got %q", v.Type, k, TypeOf(val))
				err = errors.Wrap(err, k)
				break
			}
		}

		// Run
		if err != nil {
			report = &Report{
				Level:        Unknown,
				CurrentState: NoDepStatus,
				ReportSpec: ReportSpec{
					Problem: JSONableData{
						Description: "reason",
						Data:        fmt.Sprintf("%v", err),
					},
				},
			}
		} else {
			report, err = func() (report *Report, err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panic occured %v", r)
					}
					if err != nil {
						log.Error(dctx.Context(), "Error: %v", err)
					}
					log.Info(dctx.Context(), "done")
				}()
				var spec *ReportSpec
				spec, err = consumer.Do(*dctx)
				report = &Report{ReportSpec: *spec}
				return
			}()
			if report == nil && err == nil {
				err = errors.New("No report from test")
			}
			if err != nil {
				report = &Report{
					Level:        Fatal,
					CurrentState: ErrOnDetectorStatus,
					ReportSpec: ReportSpec{
						Problem: JSONableData{
							Description: "Detector is failed with following error",
							Data:        fmt.Sprintf("%v", err),
						},
					},
				}
			} else {
				report.Level = Normal
				report.CurrentState = NormalStatus
				if !report.HasPassed {
					report.Level = meta.Level
					report.CurrentState = meta.IfHappened
				}
			}
		}

		report.MetaInfo = meta.MetaInfo
		report.CreatedAt = time.Now()

		log.Info(ctx, "%v", report)
		reports = append(reports, *report)
	}

	// Writing report
	result.FinishedAt = time.Now()
	result.Reports = reports
	return &result, nil
}
