/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package health

import "context"

type Check func(ctx context.Context) error

type Checker struct {
	checks map[string]Check
}

func New() *Checker {
	return &Checker{checks: make(map[string]Check)}
}

func (c *Checker) Register(name string, check Check) {
	c.checks[name] = check
}

func (c *Checker) Check(ctx context.Context) (ready bool, results map[string]string) {
	results = make(map[string]string, len(c.checks))
	ready = true
	for name, check := range c.checks {
		if err := check(ctx); err != nil {
			ready = false
			results[name] = err.Error()
			continue
		}
		results[name] = "ok"
	}
	return ready, results
}
