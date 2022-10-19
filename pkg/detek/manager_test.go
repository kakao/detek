package detek

import (
	"context"
	"fmt"
	"testing"
)

var (
	TypeA  = TypeOf("")
	TypeB  = TypeOf(1)
	ValueA = "abc"
	ValueB = 1
)

func TestManager_Run(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		Detector  []Detector
		Collector []Collector
		store     *Store
	}
	type args struct {
		ctx  context.Context
		opts *MangerRunOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Report
		wantErr bool
	}{
		{
			name: "All Dectectors are passed",
			fields: fields{
				Collector: []Collector{
					FakeCollector{
						Name:      "col-1",
						Required:  []FD{},
						Producing: []FD{{Key: "typeA", Value: ValueA, ShouldProduce: true}},
					},
					FakeCollector{
						Name:      "col-2",
						Required:  []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						Producing: []FD{{Key: "typeB", Value: ValueB, ShouldProduce: true}},
					},
				},
				Detector: []Detector{
					FakeDetector{
						Name:        "det-1",
						Required:    []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						ShoudPassed: true},
					FakeDetector{
						Name:        "det-2",
						Required:    []FD{{Key: "typeB", Value: ValueB, ShouldConsume: true}},
						ShoudPassed: true,
					},
				},
				store: &Store{kv: make(map[string]Stored)},
			}, args: args{ctx: ctx},
			want: []Report{
				{MetaInfo: MetaInfo{ID: "det-1"}, Level: Normal},
				{MetaInfo: MetaInfo{ID: "det-2"}, Level: Normal},
			},
		},
		{
			name: "One of Dectectors are NOT passed",
			fields: fields{
				Collector: []Collector{
					FakeCollector{
						Name:      "col-1",
						Required:  []FD{},
						Producing: []FD{{Key: "typeA", Value: ValueA, ShouldProduce: true}},
					},
					FakeCollector{
						Name:      "col-2",
						Required:  []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						Producing: []FD{{Key: "typeB", Value: ValueB, ShouldProduce: true}},
					},
				},
				Detector: []Detector{
					FakeDetector{
						Name:        "det-1",
						Required:    []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						ShoudPassed: false},
					FakeDetector{
						Name:        "det-2",
						Required:    []FD{{Key: "typeB", Value: ValueB, ShouldConsume: true}},
						ShoudPassed: true,
					},
				},
				store: &Store{kv: make(map[string]Stored)},
			}, args: args{ctx: ctx},
			want: []Report{
				{MetaInfo: MetaInfo{ID: "det-1"}, Level: Error},
				{MetaInfo: MetaInfo{ID: "det-2"}, Level: Normal},
			},
		},
		{
			name: "One of Detectors return error",
			fields: fields{
				Collector: []Collector{
					FakeCollector{
						Name:      "col-1",
						Required:  []FD{},
						Producing: []FD{{Key: "typeA", Value: ValueA, ShouldProduce: true}},
					},
					FakeCollector{
						Name:      "col-2",
						Required:  []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						Producing: []FD{{Key: "typeB", Value: ValueB, ShouldProduce: true}},
					},
				},
				Detector: []Detector{
					FakeDetector{
						Name:     "det-1",
						Required: []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						IsError:  true,
					},
					FakeDetector{
						Name:        "det-2",
						Required:    []FD{{Key: "typeB", Value: ValueB, ShouldConsume: true}},
						ShoudPassed: true,
					},
				},
				store: &Store{kv: make(map[string]Stored)},
			}, args: args{ctx: ctx},
			want: []Report{
				{MetaInfo: MetaInfo{ID: "det-1"}, Level: Fatal},
				{MetaInfo: MetaInfo{ID: "det-2"}, Level: Normal},
			},
		},
		{
			name: "Panic on Detector",
			fields: fields{
				Collector: []Collector{
					FakeCollector{
						Name:      "col-1",
						Required:  []FD{},
						Producing: []FD{{Key: "typeA", Value: ValueA, ShouldProduce: true}},
					},
					FakeCollector{
						Name:      "col-2",
						Required:  []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						Producing: []FD{{Key: "typeB", Value: ValueB, ShouldProduce: true}},
					},
				},
				Detector: []Detector{
					FakeDetector{
						Name:     "det-1",
						Required: []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						IsError:  true,
					},
					FakeDetector{
						Name:     "det-2",
						Required: []FD{{Key: "typeB", Value: ValueB, ShouldConsume: true}},
						IsPanic:  true,
					},
				},
				store: &Store{kv: make(map[string]Stored)},
			}, args: args{ctx: ctx},
			want: []Report{
				{MetaInfo: MetaInfo{ID: "det-1"}, Level: Fatal},
				{MetaInfo: MetaInfo{ID: "det-2"}, Level: Fatal},
			},
		},
		{
			name: "Collector did not produce some of data",
			fields: fields{
				Collector: []Collector{
					FakeCollector{
						Name:      "col-1",
						Required:  []FD{},
						Producing: []FD{{Key: "typeA", Value: ValueA, ShouldProduce: true}},
					},
					FakeCollector{
						Name:      "col-2",
						Required:  []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						Producing: []FD{{Key: "typeB", Value: ValueB, ShouldProduce: false}},
					},
				},
				Detector: []Detector{
					FakeDetector{
						Name:        "det-1",
						Required:    []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						ShoudPassed: true},
					FakeDetector{
						Name:        "det-2",
						Required:    []FD{{Key: "typeB", Value: ValueB, ShouldConsume: true}},
						ShoudPassed: true,
					},
				},
				store: &Store{kv: make(map[string]Stored)},
			}, args: args{ctx: ctx},
			want: []Report{
				{MetaInfo: MetaInfo{ID: "det-1"}, Level: Normal},
				{MetaInfo: MetaInfo{ID: "det-2"}, Level: Unknown},
			},
		},
		{
			name: "Panic on Collector",
			fields: fields{
				Collector: []Collector{
					FakeCollector{
						Name:      "col-1",
						Required:  []FD{},
						Producing: []FD{{Key: "typeA", Value: ValueA, ShouldProduce: true}},
						IsPanic:   true,
					},
					FakeCollector{
						Name:      "col-2",
						Required:  []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						Producing: []FD{{Key: "typeB", Value: ValueB, ShouldProduce: true}},
					},
				},
				Detector: []Detector{
					FakeDetector{
						Name:        "det-1",
						Required:    []FD{{Key: "typeA", Value: ValueA, ShouldConsume: true}},
						ShoudPassed: true},
					FakeDetector{
						Name:        "det-2",
						Required:    []FD{{Key: "typeB", Value: ValueB, ShouldConsume: true}},
						ShoudPassed: true,
					},
				},
				store: &Store{kv: make(map[string]Stored)},
			}, args: args{ctx: ctx},
			want: []Report{
				{MetaInfo: MetaInfo{ID: "collector_reports"}, Level: Unknown},
				{MetaInfo: MetaInfo{ID: "det-1"}, Level: Unknown},
				{MetaInfo: MetaInfo{ID: "det-2"}, Level: Unknown},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				Detector:  tt.fields.Detector,
				Collector: tt.fields.Collector,
				store:     tt.fields.store,
			}
			got, err := m.Run(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := hasReport(tt.want, got.Reports); err != nil {
				t.Error(err, "expected reports not found")
			}
			if err := hasReport(got.Reports, tt.want); err != nil {
				t.Error(err, "unexpected report apeared")
			}
		})
	}
}

func hasReport(src, tgt []Report) error {
	for _, a := range src {
		isEqual := false
		for _, b := range tgt {
			if a.ID != b.ID {
				continue
			}
			if a.Level != b.Level {
				continue
			}
			// just check ID and Level
			isEqual = true
			break
		}
		if !isEqual {
			return fmt.Errorf("there's no matching report for %q", a.ID)
		}
	}
	return nil
}
