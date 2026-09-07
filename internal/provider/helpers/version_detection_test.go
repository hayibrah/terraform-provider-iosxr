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

package helpers

import (
	"encoding/base64"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
		wantOk bool
	}{
		{name: "three_dotted_decimal", input: "24.4.2", want: "24.4", wantOk: true},
		{name: "two_dotted_decimal", input: "25.4", want: "25.4", wantOk: true},
		{name: "legacy_format", input: "2442", want: "24.4", wantOk: true},
		{name: "empty_string", input: "", want: "", wantOk: false},
		{name: "malformed_word", input: "garbage", want: "", wantOk: false},
		{name: "malformed_dotted_alpha", input: "abc.def", want: "", wantOk: false},
		{name: "gnmi_version", input: "0.10.0", want: "0.10", wantOk: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := NormalizeVersion(tc.input)
			if got != tc.want || ok != tc.wantOk {
				t.Errorf("NormalizeVersion(%q) = %q, %v; want %q, %v",
					tc.input, got, ok, tc.want, tc.wantOk)
			}
		})
	}
}

func TestValidateSupportedVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "three_dotted_decimal", input: "24.4.2", want: true},
		{name: "two_dotted_decimal", input: "25.4", want: true},
		{name: "legacy_format", input: "2442", want: true},
		{name: "empty_string", input: "", want: false},
		{name: "malformed_word", input: "garbage", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateSupportedVersion(tc.input)
			if got != tc.want {
				t.Errorf("ValidateSupportedVersion(%q) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestExtractVersionFromCLIString(t *testing.T) {
	const cliConfig = `!! Building configuration...
!! IOS XR Configuration 24.4.2
!! Last configuration change at Sun Sep  6 21:52:11 2026 by cisco
!
hostname xrv9k-24-1
!
end
`
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "version_header", input: cliConfig, want: "24.4.2"},
		{name: "no_version_header", input: "hostname R1\ninterface Loopback0\n!", want: ""},
		{name: "empty_string", input: "", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractVersionFromCLIString(tc.input)
			if got != tc.want {
				t.Errorf("extractVersionFromCLIString(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestExtractVersionString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "version_header", input: "!! IOS XR Configuration 25.4.1", want: "25.4.1"},
		{name: "gnmi_version", input: "0.10.0", want: "0.10.0"},
		{name: "no_version", input: "hostname R1", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractVersionString(tc.input)
			if got != tc.want {
				t.Errorf("extractVersionString(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestDecodeBase64CLI(t *testing.T) {
	plaintext := "!! IOS XR Configuration 24.4.2\nhostname R1\n"
	encoded := base64.StdEncoding.EncodeToString([]byte(plaintext))

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "valid_base64", input: encoded, want: plaintext},
		{name: "invalid_base64", input: "not-valid-base64!!!", want: ""},
		{name: "empty_string", input: "", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := decodeBase64CLI(tc.input)
			if got != tc.want {
				t.Errorf("decodeBase64CLI(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}
