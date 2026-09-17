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

// Tests for the generator merge logic.
//
// generator.go uses //go:build ignore (it is a go run tool, not a library),
// so these tests must be run with explicit file arguments:
//
//	go test -v generator.go generator_test.go
//
// Running "go test ./gen/..." will not find them.

//go:build ignore

package main

import (
	"testing"
)

// ---------------------------------------------------------------------------
// mergeConfigs tests
// ---------------------------------------------------------------------------

func TestMergeConfigs(t *testing.T) {
	tests := []struct {
		name     string
		base     YamlConfig
		override YamlConfig
		check    func(t *testing.T, got YamlConfig)
	}{
		{
			name: "basic field override",
			base: YamlConfig{
				Name: "Logging",
				Path: "Cisco-IOS-XR-um-logging-cfg:/logging",
			},
			override: YamlConfig{
				Version:        "25.4",
				ResDescription: "Updated description",
				DocCategory:    "Logging",
			},
			check: func(t *testing.T, got YamlConfig) {
				if got.Name != "Logging" {
					t.Errorf("Name: got %q, want %q", got.Name, "Logging")
				}
				if got.Path != "Cisco-IOS-XR-um-logging-cfg:/logging" {
					t.Errorf("Path: got %q, should be unchanged", got.Path)
				}
				if got.ResDescription != "Updated description" {
					t.Errorf("ResDescription: got %q, want %q", got.ResDescription, "Updated description")
				}
				if got.DocCategory != "Logging" {
					t.Errorf("DocCategory: got %q, want %q", got.DocCategory, "Logging")
				}
			},
		},
		{
			name:     "name override",
			base:     YamlConfig{Name: "Service Timestamps Old", Path: "old-module:/service/timestamps"},
			override: YamlConfig{Version: "25.4", Name: "Service Timestamps", Path: "new-module:/service/timestamps"},
			check: func(t *testing.T, got YamlConfig) {
				if got.Name != "Service Timestamps" {
					t.Errorf("Name: got %q, want %q", got.Name, "Service Timestamps")
				}
				if got.Path != "new-module:/service/timestamps" {
					t.Errorf("Path: got %q, want %q", got.Path, "new-module:/service/timestamps")
				}
			},
		},
		{
			name:     "path unchanged",
			base:     YamlConfig{Name: "Logging", Path: "Cisco-IOS-XR-um-logging-cfg:/logging"},
			override: YamlConfig{Version: "25.4", ResDescription: "New desc"},
			check: func(t *testing.T, got YamlConfig) {
				if got.Path != "Cisco-IOS-XR-um-logging-cfg:/logging" {
					t.Errorf("Path: got %q, want %q (unchanged)", got.Path, "Cisco-IOS-XR-um-logging-cfg:/logging")
				}
			},
		},
		{
			name: "legacy resource",
			base: YamlConfig{
				Name: "Old Feature",
				Path: "old-module:/feature",
				Attributes: []YamlConfigAttribute{
					{YangName: "attr1", TfName: "attr1", Type: "String"},
				},
			},
			override: YamlConfig{Version: "25.4", Legacy: true},
			check: func(t *testing.T, got YamlConfig) {
				if got.RemovedInVersion != "25.4" {
					t.Errorf("RemovedInVersion: got %q, want %q", got.RemovedInVersion, "25.4")
				}
				if !got.Legacy {
					t.Error("Legacy: got false, want true")
				}
				// Attributes should be unchanged — early return preserves base state
				if len(got.Attributes) != 1 {
					t.Errorf("Attributes: got %d, want 1 (preserved from base)", len(got.Attributes))
				}
			},
		},
		{
			name:     "no_delete propagates",
			base:     YamlConfig{Name: "Feature", NoDelete: false},
			override: YamlConfig{Version: "25.4", NoDelete: true},
			check: func(t *testing.T, got YamlConfig) {
				if !got.NoDelete {
					t.Error("NoDelete: got false, want true")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeConfigs(tc.base, tc.override)
			tc.check(t, got)
		})
	}
}

// ---------------------------------------------------------------------------
// mergeAttributes tests
// ---------------------------------------------------------------------------

func TestMergeAttributes(t *testing.T) {
	tests := []struct {
		name     string
		base     []YamlConfigAttribute
		override []YamlConfigAttribute
		check    func(t *testing.T, got []YamlConfigAttribute)
	}{
		{
			name: "new attr gets AddedInVersion",
			base: []YamlConfigAttribute{
				{YangName: "console", TfName: "console", Type: "String"},
			},
			override: []YamlConfigAttribute{
				{YangName: "new-leaf", TfName: "new_leaf", Type: "String"},
			},
			check: func(t *testing.T, got []YamlConfigAttribute) {
				if len(got) != 2 {
					t.Fatalf("len: got %d, want 2", len(got))
				}
				newAttr := got[1]
				if newAttr.YangName != "new-leaf" {
					t.Errorf("YangName: got %q, want %q", newAttr.YangName, "new-leaf")
				}
				if newAttr.AddedInVersion != "25.4" {
					t.Errorf("AddedInVersion: got %q, want %q", newAttr.AddedInVersion, "25.4")
				}
			},
		},
		{
			name: "legacy attr gets RemovedInVersion",
			base: []YamlConfigAttribute{
				{YangName: "source-interface-name", TfName: "name", Type: "String", Id: true},
			},
			override: []YamlConfigAttribute{
				{YangName: "source-interface-name", Legacy: true},
			},
			check: func(t *testing.T, got []YamlConfigAttribute) {
				if len(got) != 1 {
					t.Fatalf("len: got %d, want 1 (legacy keeps the attr)", len(got))
				}
				if got[0].RemovedInVersion != "25.4" {
					t.Errorf("RemovedInVersion: got %q, want %q", got[0].RemovedInVersion, "25.4")
				}
				if !got[0].Legacy {
					t.Error("Legacy: got false, want true")
				}
				// TfName should be unchanged when override doesn't specify one
				if got[0].TfName != "name" {
					t.Errorf("TfName: got %q, want %q (unchanged)", got[0].TfName, "name")
				}
			},
		},
		{
			// F7 fix: legacy block with explicit tf_name renames the base attribute
			// so the natural name is freed for a replacement attribute.
			name: "legacy attr renames tf_name",
			base: []YamlConfigAttribute{
				{YangName: "console", TfName: "console", Type: "String"},
			},
			override: []YamlConfigAttribute{
				{YangName: "console", TfName: "console_legacy", Legacy: true},
			},
			check: func(t *testing.T, got []YamlConfigAttribute) {
				if len(got) != 1 {
					t.Fatalf("len: got %d, want 1", len(got))
				}
				if got[0].TfName != "console_legacy" {
					t.Errorf("TfName: got %q, want %q", got[0].TfName, "console_legacy")
				}
				if got[0].RemovedInVersion != "25.4" {
					t.Errorf("RemovedInVersion: got %q, want %q", got[0].RemovedInVersion, "25.4")
				}
			},
		},
		{
			// A legacy entry in the override that has no matching base attr is a no-op.
			name: "legacy attr not in base is skipped",
			base: []YamlConfigAttribute{
				{YangName: "existing", TfName: "existing", Type: "String"},
			},
			override: []YamlConfigAttribute{
				{YangName: "ghost", Legacy: true},
			},
			check: func(t *testing.T, got []YamlConfigAttribute) {
				if len(got) != 1 {
					t.Errorf("len: got %d, want 1 (ghost attr should be skipped)", len(got))
				}
			},
		},
		{
			name: "existing attr fields updated",
			base: []YamlConfigAttribute{
				{YangName: "severity", TfName: "severity", Type: "String", Description: "old desc"},
			},
			override: []YamlConfigAttribute{
				{YangName: "severity", Description: "new desc", DefaultValue: "informational"},
			},
			check: func(t *testing.T, got []YamlConfigAttribute) {
				if len(got) != 1 {
					t.Fatalf("len: got %d, want 1", len(got))
				}
				if got[0].Description != "new desc" {
					t.Errorf("Description: got %q, want %q", got[0].Description, "new desc")
				}
				// DefaultValue introduced in an override becomes a versioned default (F17 VersionDefaults).
				// The base had no default (""), so this is a version-scoped change: stored in
				// VersionDefaults, DefaultValue is cleared to "".
				if got[0].DefaultValue != "" {
					t.Errorf("DefaultValue: got %q, want empty (versioned default stored in VersionDefaults)", got[0].DefaultValue)
				}
				if got[0].VersionDefaults["25.4"] != "informational" {
					t.Errorf("VersionDefaults[25.4]: got %q, want %q", got[0].VersionDefaults["25.4"], "informational")
				}
				// Version not stamped on existing attrs
				if got[0].AddedInVersion != "" {
					t.Errorf("AddedInVersion: got %q, want empty (existing attr)", got[0].AddedInVersion)
				}
			},
		},
		{
			// When yang_names differ but tf_names match, merge uses tf_name match.
			name: "match by tf_name",
			base: []YamlConfigAttribute{
				{YangName: "old-yang", TfName: "shared_name", Type: "String"},
			},
			override: []YamlConfigAttribute{
				{YangName: "new-yang", TfName: "shared_name", Type: "Int64"},
			},
			check: func(t *testing.T, got []YamlConfigAttribute) {
				if len(got) != 1 {
					t.Fatalf("len: got %d, want 1 (matched by tf_name)", len(got))
				}
				if got[0].Type != "Int64" {
					t.Errorf("Type: got %q, want %q", got[0].Type, "Int64")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeAttributes(tc.base, tc.override, "25.4")
			tc.check(t, got)
		})
	}
}

func TestMergeAttributes_NestedAttrsInheritAddedInVersion(t *testing.T) {
	// A completely new list attr in the override: its nested attrs should be
	// stamped with AddedInVersion too.
	base := []YamlConfigAttribute{}
	override := []YamlConfigAttribute{
		{
			YangName: "new-list",
			TfName:   "new_list",
			Type:     "List",
			Attributes: []YamlConfigAttribute{
				{YangName: "key", TfName: "key", Type: "String", Id: true},
				{YangName: "value", TfName: "value", Type: "String"},
			},
		},
	}

	got := mergeAttributes(base, override, "25.4")

	if len(got) != 1 {
		t.Fatalf("len: got %d, want 1", len(got))
	}
	if got[0].AddedInVersion != "25.4" {
		t.Errorf("outer AddedInVersion: got %q, want %q", got[0].AddedInVersion, "25.4")
	}
	for _, child := range got[0].Attributes {
		if child.AddedInVersion != "25.4" {
			t.Errorf("child %q AddedInVersion: got %q, want %q",
				child.YangName, child.AddedInVersion, "25.4")
		}
	}
}

func TestMergeAttributes_CompositeKeyPromotion(t *testing.T) {
	// Reproduces the logging source_interfaces scenario:
	// 24.4 key: source-interface-name (id:true)
	// 25.4: retire source-interface-name, add interface-name + vrf-name (both id:true)
	base := []YamlConfigAttribute{
		{YangName: "source-interface-name", TfName: "name", Type: "String", Id: true},
	}
	override := []YamlConfigAttribute{
		{YangName: "source-interface-name", Legacy: true},
		{YangName: "interface-name", TfName: "interface_name", Type: "String", Id: true},
		{YangName: "vrf-name", TfName: "vrf_name", Type: "String", Id: true},
	}

	got := mergeAttributes(base, override, "25.4")

	// Expect 3 attrs: retired name, new interface_name, new vrf_name
	if len(got) != 3 {
		t.Fatalf("len: got %d, want 3", len(got))
	}

	// Index 0: source-interface-name retired
	if got[0].RemovedInVersion != "25.4" {
		t.Errorf("name RemovedInVersion: got %q, want %q", got[0].RemovedInVersion, "25.4")
	}
	if !got[0].Id {
		t.Error("name Id: got false, want true (should retain Id)")
	}

	// Index 1: interface-name added
	if got[1].YangName != "interface-name" {
		t.Errorf("got[1] YangName: got %q, want %q", got[1].YangName, "interface-name")
	}
	if got[1].AddedInVersion != "25.4" {
		t.Errorf("interface-name AddedInVersion: got %q, want %q", got[1].AddedInVersion, "25.4")
	}
	if !got[1].Id {
		t.Error("interface-name Id: got false, want true")
	}

	// Index 2: vrf-name added
	if got[2].YangName != "vrf-name" {
		t.Errorf("got[2] YangName: got %q, want %q", got[2].YangName, "vrf-name")
	}
	if got[2].AddedInVersion != "25.4" {
		t.Errorf("vrf-name AddedInVersion: got %q, want %q", got[2].AddedInVersion, "25.4")
	}
}

func TestMergeAttributes_ReplacesYangName_OnKeyAttr(t *testing.T) {
	// Verifies that replaces_yang_name on an id:true attribute is fully processed:
	// Phase 1 (mergeAttributes): VersionYangNames populated with "_base" placeholder.
	// Phase 2 (fixAttributeBaseVersion): "_base" replaced with real base version,
	// MovedInVersion derived.
	base := []YamlConfigAttribute{
		{YangName: "name", TfName: "name", Type: "String", Id: true},
	}
	override := []YamlConfigAttribute{
		{YangName: "host", TfName: "name", Type: "String", Id: true, ReplacesYangName: "name"},
	}

	// Phase 1
	got := mergeAttributes(base, override, "25.4")
	if len(got) != 1 {
		t.Fatalf("len: got %d, want 1", len(got))
	}
	if got[0].YangName != "host" {
		t.Errorf("YangName: got %q, want %q", got[0].YangName, "host")
	}
	if !got[0].Id {
		t.Error("Id: got false, want true")
	}
	if got[0].TfName != "name" {
		t.Errorf("TfName: got %q, want %q", got[0].TfName, "name")
	}
	if got[0].AddedInVersion != "" {
		t.Errorf("AddedInVersion: got %q, want empty (not treated as new attr)", got[0].AddedInVersion)
	}
	if got[0].RemovedInVersion != "" {
		t.Errorf("RemovedInVersion: got %q, want empty", got[0].RemovedInVersion)
	}
	if got[0].VersionYangNames["_base"] != "name" {
		t.Errorf("VersionYangNames[_base]: got %q, want %q", got[0].VersionYangNames["_base"], "name")
	}
	if got[0].VersionYangNames["25.4"] != "host" {
		t.Errorf("VersionYangNames[25.4]: got %q, want %q", got[0].VersionYangNames["25.4"], "host")
	}
	if got[0].MovedInVersion != "" {
		t.Errorf("MovedInVersion: got %q, want empty (not set until fixAttributeBaseVersion)", got[0].MovedInVersion)
	}

	// Phase 2
	fixAttributeBaseVersion(&got[0], "24.4")
	if got[0].VersionYangNames["24.4"] != "name" {
		t.Errorf("VersionYangNames[24.4]: got %q, want %q", got[0].VersionYangNames["24.4"], "name")
	}
	if _, hasBase := got[0].VersionYangNames["_base"]; hasBase {
		t.Error("VersionYangNames[_base]: still present after fixAttributeBaseVersion, want removed")
	}
	if got[0].MovedInVersion != "25.4" {
		t.Errorf("MovedInVersion: got %q, want %q", got[0].MovedInVersion, "25.4")
	}
}

func TestMergeAttributes_ReplacesYangName_ThreeVersionChain(t *testing.T) {
	// Simulates an attribute that moves in 25.4 and again in 26.2.
	// Each delta points replaces_yang_name at the previous canonical yang_name.
	base := []YamlConfigAttribute{
		{YangName: "monitor", TfName: "monitor", Type: "String"},
	}
	// First move: 24.4 "monitor" → 25.4 "monitor/monitor-level"
	after25 := mergeAttributes(base, []YamlConfigAttribute{
		{YangName: "monitor/monitor-level", TfName: "monitor", Type: "String", ReplacesYangName: "monitor"},
	}, "25.4")

	// Second move: current canonical "monitor/monitor-level" → 26.2 "monitor/monitor-level/mode"
	after26 := mergeAttributes(after25, []YamlConfigAttribute{
		{YangName: "monitor/monitor-level/mode", TfName: "monitor", Type: "String", ReplacesYangName: "monitor/monitor-level"},
	}, "26.2")

	if len(after26) != 1 {
		t.Fatalf("len: got %d, want 1", len(after26))
	}
	attr := after26[0]
	if attr.YangName != "monitor/monitor-level/mode" {
		t.Errorf("YangName: got %q, want %q", attr.YangName, "monitor/monitor-level/mode")
	}
	if attr.VersionYangNames["_base"] != "monitor" {
		t.Errorf("VersionYangNames[_base]: got %q, want %q", attr.VersionYangNames["_base"], "monitor")
	}
	if attr.VersionYangNames["25.4"] != "monitor/monitor-level" {
		t.Errorf("VersionYangNames[25.4]: got %q, want %q", attr.VersionYangNames["25.4"], "monitor/monitor-level")
	}
	if attr.VersionYangNames["26.2"] != "monitor/monitor-level/mode" {
		t.Errorf("VersionYangNames[26.2]: got %q, want %q", attr.VersionYangNames["26.2"], "monitor/monitor-level/mode")
	}

	// After fixAttributeBaseVersion, "_base" → "24.4" and MovedInVersion is the earliest move.
	fixAttributeBaseVersion(&attr, "24.4")
	if attr.VersionYangNames["24.4"] != "monitor" {
		t.Errorf("VersionYangNames[24.4]: got %q, want %q", attr.VersionYangNames["24.4"], "monitor")
	}
	if _, hasBase := attr.VersionYangNames["_base"]; hasBase {
		t.Error("VersionYangNames[_base]: still present after fixAttributeBaseVersion, want removed")
	}
	// MovedInVersion is the earliest non-base version key ("25.4" < "26.2").
	if attr.MovedInVersion != "25.4" {
		t.Errorf("MovedInVersion: got %q, want %q", attr.MovedInVersion, "25.4")
	}
}

// ---------------------------------------------------------------------------
// VersionDeleteMode merge scenarios
// ---------------------------------------------------------------------------

func TestVersionDeleteMode(t *testing.T) {
	tests := []struct {
		name                   string
		baseDeleteParent       bool
		baseDeleteGrandparent  bool
		deltaReplacesYangName  string
		deltaDeleteParent      bool
		deltaDeleteGrandparent bool
		wantVersionDeleteMode  map[string]string // nil means expect nil (inherited/unchanged)
		checkFlags             bool              // also assert the OR-merged DeleteParent/DeleteGrandparent scalars
		wantDeleteParent       bool
		wantDeleteGrandparent  bool
	}{
		{
			// base has delete_parent, 25.4 has same yang_name (no replaces_yang_name) and no delete flag.
			// Expected: 25.4 inherits delete_parent — no VersionDeleteMode created.
			name:             "inherit delete_parent (no rename, no flag)",
			baseDeleteParent: true,
			checkFlags:       true,
			wantDeleteParent: true,
		},
		{
			// base has delete_grandparent, 25.4 has same yang_name (no replaces_yang_name) and no delete flag.
			// Expected: 25.4 inherits delete_grandparent — no VersionDeleteMode created.
			name:                  "inherit delete_grandparent (no rename, no flag)",
			baseDeleteGrandparent: true,
			checkFlags:            true,
			wantDeleteGrandparent: true,
		},
		{
			// base has delete_parent, 25.4 has replaces_yang_name and NO delete flag.
			// Expected: version-specific map — 24.4 deletes parent, 25.4 deletes leaf directly.
			name:                  "replacement drops delete_parent (rename, no flag)",
			baseDeleteParent:      true,
			deltaReplacesYangName: "maxfilesize",
			wantVersionDeleteMode: map[string]string{"_base": "parent", "25.4": ""},
		},
		{
			// base has delete_parent, 25.4 has replaces_yang_name WITH delete_parent.
			// Expected: both versions delete parent — no VersionDeleteMode, static path.
			name:                  "replacement keeps delete_parent (rename, with flag)",
			baseDeleteParent:      true,
			deltaReplacesYangName: "maxfilesize",
			deltaDeleteParent:     true,
			checkFlags:            true,
			wantDeleteParent:      true,
		},
		{
			// Same as above, delete_grandparent variant for base.
			name:                  "replacement drops delete_grandparent (rename, no flag)",
			baseDeleteGrandparent: true,
			deltaReplacesYangName: "maxfilesize",
			wantVersionDeleteMode: map[string]string{"_base": "grandparent", "25.4": ""},
		},
		{
			// base has delete_grandparent, 25.4 has replaces_yang_name WITH delete_grandparent.
			// Expected: no VersionDeleteMode, same mode throughout.
			name:                   "replacement keeps delete_grandparent (rename, with flag)",
			baseDeleteGrandparent:  true,
			deltaReplacesYangName:  "maxfilesize",
			deltaDeleteGrandparent: true,
			checkFlags:             true,
			wantDeleteGrandparent:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base := []YamlConfigAttribute{{
				YangName:          "maxfilesize",
				TfName:            "maxfilesize",
				DeleteParent:      tc.baseDeleteParent,
				DeleteGrandparent: tc.baseDeleteGrandparent,
			}}
			yangName := "maxfilesize"
			if tc.deltaReplacesYangName != "" {
				yangName = "path/maxfilesize"
			}
			delta := []YamlConfigAttribute{{
				YangName:          yangName,
				TfName:            "maxfilesize",
				ReplacesYangName:  tc.deltaReplacesYangName,
				DeleteParent:      tc.deltaDeleteParent,
				DeleteGrandparent: tc.deltaDeleteGrandparent,
			}}
			got := mergeAttributes(base, delta, "25.4")
			attr := got[0]

			if tc.wantVersionDeleteMode == nil {
				if attr.VersionDeleteMode != nil {
					t.Errorf("VersionDeleteMode should be nil, got %v", attr.VersionDeleteMode)
				}
			} else {
				if attr.VersionDeleteMode == nil {
					t.Fatal("VersionDeleteMode should be set")
				}
				for k, want := range tc.wantVersionDeleteMode {
					if got := attr.VersionDeleteMode[k]; got != want {
						t.Errorf("VersionDeleteMode[%q] = %q, want %q (full map: %v)", k, got, want, attr.VersionDeleteMode)
					}
				}
			}

			if tc.checkFlags {
				if attr.DeleteParent != tc.wantDeleteParent {
					t.Errorf("DeleteParent = %v, want %v", attr.DeleteParent, tc.wantDeleteParent)
				}
				if attr.DeleteGrandparent != tc.wantDeleteGrandparent {
					t.Errorf("DeleteGrandparent = %v, want %v", attr.DeleteGrandparent, tc.wantDeleteGrandparent)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// getVersionValue tests
// ---------------------------------------------------------------------------

func TestGetVersionValue(t *testing.T) {
	const (
		defVal   = "default"
		newVal   = "new"
		newerVal = "newer"
	)

	byVersion := map[string]string{
		"25.4": newVal,
	}

	twoThresholds := map[string]string{
		"25.2": newVal,
		"25.4": newerVal,
	}

	tests := []struct {
		name      string
		target    string
		byVersion map[string]string
		defVal    string
		want      string
	}{
		{
			name:      "empty map returns defaultValue",
			target:    "25.4",
			byVersion: map[string]string{},
			defVal:    defVal,
			want:      defVal,
		},
		{
			name:      "target below all thresholds returns defaultValue",
			target:    "24.4",
			byVersion: byVersion,
			defVal:    defVal,
			want:      defVal,
		},
		{
			name:      "exact match on threshold",
			target:    "25.4",
			byVersion: byVersion,
			defVal:    defVal,
			want:      newVal,
		},
		{
			name:      "target above threshold",
			target:    "25.6",
			byVersion: byVersion,
			defVal:    defVal,
			want:      newVal,
		},
		{
			name:      "target between two thresholds returns lower threshold value",
			target:    "25.3",
			byVersion: twoThresholds,
			defVal:    defVal,
			want:      newVal,
		},
		{
			name:      "exact match on upper threshold",
			target:    "25.4",
			byVersion: twoThresholds,
			defVal:    defVal,
			want:      newerVal,
		},
		{
			name:      "target above all thresholds returns highest threshold value",
			target:    "26.1",
			byVersion: twoThresholds,
			defVal:    defVal,
			want:      newerVal,
		},
		{
			name:      "patch component ignored: 24.4.2 treated as 24.4",
			target:    "24.4.2",
			byVersion: byVersion,
			defVal:    defVal,
			want:      defVal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getVersionValue(tc.target, tc.byVersion, tc.defVal)
			if got != tc.want {
				t.Errorf("getVersionValue(%q, %v, %q) = %q, want %q",
					tc.target, tc.byVersion, tc.defVal, got, tc.want)
			}
		})
	}
}

// TestGetDeletePathExpr_RenameWithoutModeChange_ThreeVersionChain reproduces the read-side
// stale-fallback bug fixed by getVersionValue: a mode change recorded at 25.4, followed by a
// 26.2 rename that restates the same mode (correctly producing no new VersionDeleteMode
// entry), used to bake the wrong, stale oldest mode into 26.2's generated path.
func TestGetDeletePathExpr_RenameWithoutModeChange_ThreeVersionChain(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:     "o1/m1/a",
		TfName:       "x",
		DeleteParent: true,
	}}
	delta25 := []YamlConfigAttribute{{
		YangName:          "o2/m2/b",
		TfName:            "x",
		ReplacesYangName:  "o1/m1/a",
		DeleteGrandparent: true,
	}}
	delta26 := []YamlConfigAttribute{{
		YangName:          "o3/m3/c",
		TfName:            "x",
		ReplacesYangName:  "o2/m2/b",
		DeleteGrandparent: true, // restated, same mode as 25.4
	}}

	merged := mergeAttributes(base, delta25, "25.4")
	merged = mergeAttributes(merged, delta26, "26.2")
	attr := merged[0]

	if attr.VersionDeleteMode == nil {
		t.Fatal("VersionDeleteMode should be set after the 25.4 delta")
	}
	if _, has26 := attr.VersionDeleteMode["26.2"]; has26 {
		t.Errorf("VersionDeleteMode should have no 26.2 entry (mode unchanged from 25.4), got %v", attr.VersionDeleteMode)
	}

	fixAttributeBaseVersion(&attr, "24.4")

	got := GetDeletePathExpr(attr, "version")
	want := `helpers.SelectYangPath(version, map[string]string{"24.4": "o1/m1", "25.4": "o2", "26.2": "o3"}, "o1/m1")`
	if got != want {
		t.Errorf("GetDeletePathExpr =\n  %s\nwant\n  %s", got, want)
	}
}
