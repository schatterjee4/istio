// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapterManager

import (
	"fmt"
	"testing"

	"istio.io/mixer/pkg/adapter"
)

type testAdapter struct {
	name string
}

func (t testAdapter) Name() string                                              { return t.name }
func (testAdapter) Close() error                                                { return nil }
func (testAdapter) Description() string                                         { return "mock adapter for testing" }
func (testAdapter) DefaultConfig() adapter.AspectConfig                         { return nil }
func (testAdapter) ValidateConfig(c adapter.AspectConfig) *adapter.ConfigErrors { return nil }

type denyAdapter struct{ testAdapter }

func (denyAdapter) NewDenyChecker(env adapter.Env, cfg adapter.AspectConfig) (adapter.DenyCheckerAspect, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestRegisterDenyChecker(t *testing.T) {
	reg := newRegistry()
	adapter := denyAdapter{testAdapter{name: "foo"}}

	if err := reg.RegisterDenyChecker(adapter); err != nil {
		t.Errorf("Failed to register deny adapter with err: %v", err)
	}

	impl, ok := reg.ByImpl(adapter.Name())
	if !ok {
		t.Errorf("No adapter by impl with name %s, expected adapter: %v", adapter.Name(), adapter)
	}

	if deny, ok := impl.(denyAdapter); !ok || deny != adapter {
		t.Errorf("reg.ByImpl(%s) expected adapter '%v', actual '%v'", adapter.Name(), adapter, impl)
	}
}

type listAdapter struct{ testAdapter }

func (listAdapter) NewListChecker(env adapter.Env, cfg adapter.AspectConfig) (adapter.ListCheckerAspect, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestRegisterListChecker(t *testing.T) {
	reg := newRegistry()
	adapter := listAdapter{testAdapter{name: "foo"}}

	if err := reg.RegisterListChecker(adapter); err != nil {
		t.Errorf("Failed to register check list adapter with err: %v", err)
	}

	impl, ok := reg.ByImpl(adapter.Name())
	if !ok {
		t.Errorf("No adapter by impl with name %s, expected adapter: %v", adapter.Name(), adapter)
	}

	if deny, ok := impl.(listAdapter); !ok || deny != adapter {
		t.Errorf("reg.ByImpl(%s) expected adapter '%v', actual '%v'", adapter.Name(), adapter, impl)
	}
}

type loggerAdapter struct{ testAdapter }

func (loggerAdapter) NewLogger(env adapter.Env, cfg adapter.AspectConfig) (adapter.LoggerAspect, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestRegisterLogger(t *testing.T) {
	reg := newRegistry()
	adapter := loggerAdapter{testAdapter{name: "foo"}}

	if err := reg.RegisterLogger(adapter); err != nil {
		t.Errorf("Failed to register logging adapter with err: %v", err)
	}

	impl, ok := reg.ByImpl(adapter.Name())
	if !ok {
		t.Errorf("No adapter by impl with name %s, expected adapter: %v", adapter.Name(), adapter)
	}

	if deny, ok := impl.(loggerAdapter); !ok || deny != adapter {
		t.Errorf("reg.ByImpl(%s) expected adapter '%v', actual '%v'", adapter.Name(), adapter, impl)
	}
}

type quotaAdapter struct{ testAdapter }

func (quotaAdapter) NewQuota(env adapter.Env, cfg adapter.AspectConfig) (adapter.QuotaAspect, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestRegisterQuota(t *testing.T) {
	reg := newRegistry()
	adapter := quotaAdapter{testAdapter{name: "foo"}}

	if err := reg.RegisterQuota(adapter); err != nil {
		t.Errorf("Failed to register quota adapter with err: %v", err)
	}

	impl, ok := reg.ByImpl(adapter.Name())
	if !ok {
		t.Errorf("No adapter by impl with name %s, expected adapter: %v", adapter.Name(), adapter)
	}

	if deny, ok := impl.(quotaAdapter); !ok || deny != adapter {
		t.Errorf("reg.ByImpl(%s) expected adapter '%v', actual '%v'", adapter.Name(), adapter, impl)
	}
}

func TestCollision(t *testing.T) {
	reg := newRegistry()
	name := "some name that they both have"

	a1 := denyAdapter{testAdapter{name}}
	if err := reg.RegisterDenyChecker(a1); err != nil {
		t.Errorf("Failed to insert first adapter with err: %s", err)
	}
	if a, ok := reg.ByImpl(name); !ok || a != a1 {
		t.Errorf("Failed to get first adapter by impl name; expected: '%v', actual: '%v'", a1, a)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected to recover from panic registering duplicate adapter, but recover was nil.")
		}
	}()

	a2 := listAdapter{testAdapter{name}}
	if err := reg.RegisterListChecker(a2); err != nil {
		t.Errorf("Expected a panic inserting duplicate adapter, got err instead: %s", err)
	}
	t.Error("Should not reach this statement due to panic.")
}
