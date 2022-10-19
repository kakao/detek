package detek

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kakao/detek/pkg/log"
	"github.com/pkg/errors"
)

type DetekContext struct {
	ctx   context.Context
	opt   detekConfigOpts
	store *Store
}

type detekConfigOpts struct {
	ConsumingPlan DependencyMeta
	ProducingPlan DependencyMeta
	Meta          MetaInfo
}

func newDetekContext(ctx context.Context, name string, store *Store, opt detekConfigOpts) (*DetekContext, context.CancelFunc, error) {
	if store == nil {
		return nil, nil, fmt.Errorf("store should not be nil")
	}
	ctx = log.SetContext(ctx, name)
	_, cancel := context.WithCancel(ctx)
	return &DetekContext{
		ctx:   ctx,
		opt:   opt,
		store: store,
	}, cancel, nil
}

// Typing a result of (c *DetekContext).Get
// Just a Syntactic Sugar.
func Typing[T any](s *Stored, err error) (T, error) {
	var result T
	if err != nil {
		return result, err
	}
	if v, ok := s.Value.(T); !ok {
		return result, fmt.Errorf("type not match (given)%q != (actual)%q", TypeOf(result), TypeOf(s.Value))
	} else {
		return v, nil
	}
}

func (c *DetekContext) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

// Get will return data which was "Set" with specifed "key"
// if "v" is nil, it will just return "Stored" data. (if exists)
// if "v" is NOT a pointer type of the original data, it will return error, saying that type is incorrect.
// if "v" is a pointer type of the original data, it will attempts to copy the original data to a given pointer.
//
// example 1:
//
//	var data string
//	c.Set("some_key", "some_string_data")
//	c.Get("some_key", &data)
//	fmt.Println(data) // will print "some_string_data"
//
// example 2:
//
//	c.Set("some_key", "some_string_data")
//	data, err := Typing[string](c.Get("some_key", nil))
//	fmt.Println(data) // will print "some_string_data"
func (c *DetekContext) Get(key string, v any) (*Stored, error) {
	// fetch value from store
	val, stored, err := c.store.Get(key)
	if err != nil {
		log.Error(c.ctx, "store error:[%s] no requested key in store", key)
		return nil, errors.Wrapf(err, "fail to get %q", key)
	}
	// if i == nil, return directly
	if v == nil {
		return stored, nil
	}

	var srcVal, tgtVal reflect.Value
	var srcType, tgtType reflect.Type

	// Check if tgtval is pointer or not
	tgtVal = reflect.ValueOf(v)
	if tgtVal.Kind() != reflect.Ptr || tgtVal.IsNil() {
		log.Error(c.ctx, "store error:[%s] try insert to nil interface", key)
		return nil, fmt.Errorf("%v is not a valid type for Get, use pointer type", reflect.TypeOf(v))
	}
	// unwrap pointer
	tgtVal = tgtVal.Elem()
	tgtType = reflect.TypeOf(v).Elem()

	// verify with pre submitted plan
	plan, ok := c.opt.ConsumingPlan[key]
	if !ok {
		log.Error(c.ctx, "not allowed to access %q", key)
		return nil, fmt.Errorf("this case not allowed to access %q", key)
	}
	if plan.Type != tgtType {
		log.Error(c.ctx, "store error:[%s] try insert to not matched type %q != %q", key, tgtType, plan.Type)
		return nil, fmt.Errorf("type %q does not match with what this case planed for %q", tgtType, plan.Type)
	}

	// validation
	srcType = reflect.TypeOf(val)
	srcVal = reflect.ValueOf(val)
	if srcType != tgtType {
		log.Error(c.ctx, "store error:[%s] try insert to not matched type. this may be a bug in detek %q != %q", key, tgtType, srcType)
		return nil, fmt.Errorf("type %q does not match with what producer %q produced %q", tgtType, stored.ProducedBy.ID, srcType)
	}

	// copy to i interface{}
	log.Info(c.ctx, "store: get [%s](%s)", key, tgtType)
	tgtVal.Set(srcVal)
	return stored, nil
}

func (c *DetekContext) Set(key string, val interface{}) error {
	if val == nil {
		log.Error(c.ctx, "store error: [%s] try setting nil value", key)
		return fmt.Errorf("nil value can not be set")
	}
	plan, ok := c.opt.ProducingPlan[key]
	if !ok {
		log.Error(c.ctx, "store error: unexpected to set %q", key)
		return fmt.Errorf("it does not have a plan to produce %q", key)
	}
	stored := Stored{
		Value:      val,
		ProducedBy: &c.opt.Meta,
	}
	stored.Type = reflect.TypeOf(val)
	if plan.Type != stored.Type {
		log.Error(c.ctx, "store error: [%s] type not match %q != %q", key, stored.Type, plan.Type)
		return fmt.Errorf("type %q is not match with planed type %q", stored.Type, plan.Type)
	}
	log.Info(c.ctx, "store: key %q being set with %q", key, stored.Type)
	return c.store.Set(key, &stored)
}
