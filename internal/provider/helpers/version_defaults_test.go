// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-iosxr/internal/provider/helpers"
)

func TestGetVersionDefault(t *testing.T) {
	defaults := map[string]string{
		"24.4": "informational",
		"25.4": "warnings",
	}
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"empty version", "", ""},
		{"below all thresholds", "24.1", ""},
		{"exact lower threshold", "24.4", "informational"},
		{"between thresholds", "24.11", "informational"},
		{"exact upper threshold", "25.4", "warnings"},
		{"above all thresholds", "25.11", "warnings"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetVersionDefault(tt.version, defaults)
			if got != tt.want {
				t.Errorf("GetVersionDefault(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestGetVersionDefaultEmpty(t *testing.T) {
	if got := helpers.GetVersionDefault("24.4", nil); got != "" {
		t.Errorf("expected empty string for nil map, got %q", got)
	}
	if got := helpers.GetVersionDefault("24.4", map[string]string{}); got != "" {
		t.Errorf("expected empty string for empty map, got %q", got)
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"162", 162},
		{"0", 0},
		{"-1", -1},
		{"notanumber", 0},
		{"", 0},
	}
	for _, tt := range tests {
		got := helpers.ParseInt64(tt.input)
		if got != tt.want {
			t.Errorf("ParseInt64(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"false", false},
		{"False", false},
		{"", false},
		{"yes", false},
	}
	for _, tt := range tests {
		got := helpers.ParseBool(tt.input)
		if got != tt.want {
			t.Errorf("ParseBool(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
