// Copyright 2018 Synopsys, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hubapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const ContentTypeBdPolicyV5 = "application/vnd.blackducksoftware.policy-5+json"

type bdJsonPolicyV5 struct{}

func (bdJsonPolicyV5) GetMimeType() string {
	return ContentTypeBdPolicyV5
}

type PolicyRuleList struct {
	bdJsonPolicyV5
	ItemsListBase
	Items []PolicyRule `json:"items"`
}

type PolicyRule struct {
	bdJsonPolicyV5
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Enabled       bool             `json:"enabled"`
	Overridable   bool             `json:"overridable"`
	Severity      string           `json:"severity"`
	Expression    PolicyExpression `json:"expression"`
	CreatedAt     *time.Time       `json:"createdAt"`
	CreatedBy     string           `json:"createdBy"`
	CreatedByUser string           `json:"createdByUser"`
	UpdatedAt     *time.Time       `json:"updatedAt"`
	UpdatedBy     string           `json:"updatedBy"`
	UpdatedByUser string           `json:"updatedByUser"`
	Meta          Meta             `json:"_meta"`
}

type PolicyExpression struct {
	Operator    string       `json:"operator"`
	Expressions []Expression `json:"expressions"`
}

type Expression struct {
	Name       string              `json:"name"`
	Operation  string              `json:"operation"`
	Parameters ExpressionParameter `json:"parameters"`
}

type ExpressionParameter struct {
	Values []string                 `json:"values"`
	Data   []map[string]IntOrString `json:"data"`
}

type PolicyRuleRequest struct {
	bdJsonPolicyV5
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Enabled     bool             `json:"enabled"`
	Overridable bool             `json:"overridable"`
	Expression  PolicyExpression `json:"expression"`
	Severity    string           `json:"severity"`
}

type IntOrString struct {
	Type   Type   `json:"type"`
	IntVal int32  `json:"intVal"`
	StrVal string `json:"strVal"`
}

type Type int

const (
	Int    Type = iota // The IntOrString holds an int.
	String             // The IntOrString holds a string.
)

func (intstr *IntOrString) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		intstr.Type = String
		return json.Unmarshal(value, &intstr.StrVal)
	}
	intstr.Type = Int
	return json.Unmarshal(value, &intstr.IntVal)
}

func (intstr IntOrString) MarshalJSON() ([]byte, error) {
	switch intstr.Type {
	case Int:
		return json.Marshal(intstr.IntVal)
	case String:
		return json.Marshal(intstr.StrVal)
	default:
		return []byte{}, fmt.Errorf("impossible IntOrString.Type")
	}
}

func (pr *PolicyRule) IsEqual(obj *PolicyRule) bool {
	if !strings.EqualFold(pr.Name, obj.Name) {
		return false
	}

	if !strings.EqualFold(pr.Description, obj.Description) {
		return false
	}

	if pr.Overridable != obj.Overridable {
		return false
	}

	if !strings.EqualFold(pr.Severity, obj.Severity) {
		return false
	}

	if !reflect.DeepEqual(pr.Expression, obj.Expression) {
		return false
	}

	return true
}
