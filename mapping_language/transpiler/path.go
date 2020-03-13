// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transpiler

import (
	"fmt"

	"github.com/GoogleCloudPlatform/healthcare_data_harmonization/mapping_language/parser" /* copybara-comment: parser */
)

type pathSpec struct {
	arg, field, index string
}

// VisitPath returns a pathSpec for the given PathContext.
func (t *transpiler) VisitPath(ctx *parser.PathContext) interface{} {
	p := ctx.PathHead().Accept(t).(pathSpec)
	for i := range ctx.AllPathSegment() {
		p.field += ctx.PathSegment(i).Accept(t).(string)
	}

	if ctx.ARRAYMOD() != nil && ctx.ARRAYMOD().GetText() != "" {
		p.field += ctx.ARRAYMOD().GetText()
	}

	if ctx.OWMOD() != nil && ctx.OWMOD().GetText() != "" {
		p.field += ctx.OWMOD().GetText()
	}

	// Only one of p.arg and p.index can be filled.
	if (p.arg == "") == (p.index == "") {
		t.fail(ctx, fmt.Errorf("invalid path - expected arg xor index but got both or neither (arg %s and index %s)", p.arg, p.index))
	}

	return p
}

// VisitPath returns a partially filled pathSpec for the given PathHeadContext.
// Either the arg or index field will be filled, as appropriate.
func (t *transpiler) VisitPathHead(ctx *parser.PathHeadContext) interface{} {
	if ctx.TOKEN() != nil && ctx.TOKEN().GetText() != "" {
		return pathSpec{
			arg: ctx.TOKEN().GetText(),
		}
	}

	// ROOT_INPUT is a special case path segment since normally they cannot contain $.
	if ctx.ROOT_INPUT() != nil && ctx.ROOT_INPUT().GetText() != "" {
		return pathSpec{
			arg: ctx.ROOT_INPUT().GetText(),
		}
	}

	// ROOT is a special case path segment since it is a keyword and does not get tokenized as a TOKEN
	// TODO: Remove after sunset.
	if ctx.ROOT() != nil && ctx.ROOT().GetText() != "" {
		return pathSpec{
			arg: ctx.ROOT().GetText(),
		}
	}

	if ctx.INDEX() != nil && ctx.INDEX().GetText() != "" {
		return pathSpec{
			index: ctx.INDEX().GetText(),
		}
	}

	if ctx.WILDCARD() != nil && ctx.WILDCARD().GetText() != "" {
		return pathSpec{
			index: ctx.WILDCARD().GetText(),
		}
	}

	t.fail(ctx, fmt.Errorf("invalid path head - no token, index, or wildcard"))
	return nil
}

// VisitPathSegment returns a string of the PathSegmentContext contents.
func (t *transpiler) VisitPathSegment(ctx *parser.PathSegmentContext) interface{} {
	return ctx.GetText()
}