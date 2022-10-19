package detek

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetekContext_Get(t *testing.T) {
	ctx := context.Background()
	stored := Stored{Value: "data", Type: TypeOf("data")}
	AllowedOpts := detekConfigOpts{ConsumingPlan: DependencyMeta{"tmp": {Type: stored.Type}}}
	DisallowedOpts := detekConfigOpts{}
	var expected string
	var unexpected int

	type fields struct {
		ctx   context.Context
		opt   detekConfigOpts
		store *Store
	}
	type args struct {
		key string
		i   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Stored
		wantErr bool
	}{
		{
			name: "Normal",
			fields: fields{
				ctx:   ctx,
				store: &Store{kv: map[string]Stored{"tmp": stored}},
				opt:   AllowedOpts,
			},
			args: args{
				key: "tmp",
				i:   &expected,
			},
			want: &stored,
		},
		{
			name: "Invalid Type",
			fields: fields{
				ctx:   ctx,
				store: &Store{kv: map[string]Stored{"tmp": stored}},
				opt:   AllowedOpts,
			},
			args: args{
				key: "tmp",
				i:   &unexpected,
			},
			wantErr: true,
		},
		{
			name: "Nil Type",
			fields: fields{
				ctx:   ctx,
				store: &Store{kv: map[string]Stored{"tmp": stored}},
				opt:   AllowedOpts,
			},
			args: args{
				key: "tmp",
				i:   nil,
			},
			want: &stored,
		},
		{
			name: "Not Allowed",
			fields: fields{
				ctx:   ctx,
				store: &Store{kv: map[string]Stored{"tmp": stored}},
				opt:   DisallowedOpts,
			},
			args: args{
				key: "tmp",
				i:   &expected,
			},
			wantErr: true,
		},
		{
			name: "Not exists key",
			fields: fields{
				ctx:   ctx,
				store: &Store{kv: map[string]Stored{"tmp": stored}},
				opt:   AllowedOpts,
			},
			args: args{
				key: "NOKEY",
				i:   &expected,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DetekContext{
				ctx:   tt.fields.ctx,
				opt:   tt.fields.opt,
				store: tt.fields.store,
			}
			got, err := c.Get(tt.args.key, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetekContext.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetekContext.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetekContext_Set(t *testing.T) {
	const (
		STRING_KEY  = "some_key"
		STRING_DATA = "some_data"
	)
	newContext := func() DetekContext {
		return DetekContext{
			ctx: context.TODO(),
			opt: detekConfigOpts{
				Meta: MetaInfo{ID: "test"},
				ProducingPlan: DependencyMeta{
					STRING_KEY: DependencyInfo{TypeOf(STRING_DATA)},
				},
				ConsumingPlan: DependencyMeta{
					STRING_KEY: DependencyInfo{TypeOf(STRING_DATA)},
				},
			},
			store: &Store{kv: make(map[string]Stored), mu: sync.RWMutex{}},
		}
	}
	t.Run("Basic Usage", func(t *testing.T) {
		c := newContext()
		assert.NoError(t, c.Set(STRING_KEY, STRING_DATA))

		var tmp string
		s, err := c.Get(STRING_KEY, &tmp)
		assert.NoError(t, err)
		assert.Equal(t, reflect.TypeOf(STRING_DATA), s.Type)
	})
	t.Run("Basic Usage w/ Generics", func(t *testing.T) {
		c := newContext()
		assert.NoError(t, c.Set(STRING_KEY, STRING_DATA))

		data, err := Typing[string](c.Get(STRING_KEY, nil))
		assert.NoError(t, err)
		assert.Equal(t, STRING_DATA, data)
	})
	t.Run("Wrong Type", func(t *testing.T) {
		c := newContext()
		assert.Error(t, c.Set(STRING_KEY, true))
	})
	t.Run("Wrong Type to Get", func(t *testing.T) {
		c := newContext()
		assert.NoError(t, c.Set(STRING_KEY, STRING_DATA))

		var tmp bool
		_, err := c.Get(STRING_KEY, &tmp)
		assert.Error(t, err)

		var tmp2 string
		_, err = c.Get(STRING_KEY, tmp2)
		assert.Error(t, err)
	})
	t.Run("Wrong Type to Get w/ Generics", func(t *testing.T) {
		c := newContext()
		assert.NoError(t, c.Set(STRING_KEY, STRING_DATA))

		data, err := Typing[bool](c.Get(STRING_KEY, nil))
		assert.Error(t, err)
		assert.Empty(t, data)
	})
}

func TestTyping(t *testing.T) {
	t.Run("Basic Usage", func(t *testing.T) {
		v, err := Typing[string](&Stored{
			Value:      "test_string",
			Type:       TypeOf(""),
			ProducedBy: &MetaInfo{ID: "test_case"},
		}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "test_string", v)
	})
	t.Run("Invalid Type", func(t *testing.T) {
		v, err := Typing[bool](&Stored{
			Value:      "test_string",
			Type:       TypeOf(""),
			ProducedBy: &MetaInfo{ID: "test_case"},
		}, nil)
		assert.Error(t, err)
		assert.Empty(t, v)
	})
	t.Run("Should return error, when error", func(t *testing.T) {
		_, err := Typing[string](nil, fmt.Errorf("some error"))
		assert.Error(t, err)
	})
}
