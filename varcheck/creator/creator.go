// Copyright 2016 Palantir Technologies, Inc.
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

package creator

import (
	"fmt"
	"regexp"

	"github.com/palantir/okgo/checker"
	"github.com/palantir/okgo/okgo"

	"github.com/palantir/godel-okgo-asset-varcheck/varcheck"
)

var lineRegexp = regexp.MustCompile(`.+: (.+):(\d+):(\d+): (.+)`)

func Varcheck() checker.Creator {
	return checker.NewCreator(
		varcheck.TypeName,
		varcheck.Priority,
		func(cfgYML []byte) (okgo.Checker, error) {
			return checker.NewAmalgomatedChecker(varcheck.TypeName, checker.ParamPriority(varcheck.Priority), checker.ParamLineParserWithWd(
				func(line, wd string) okgo.Issue {
					if match := lineRegexp.FindStringSubmatch(line); match != nil {
						// varcheck prepends output with package name: remove it
						line = fmt.Sprintf("%s:%s:%s: %s", match[1], match[2], match[3], match[4])
					}
					return okgo.NewIssueFromLine(line, wd)
				},
			)), nil
		},
	)
}
