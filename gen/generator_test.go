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

func TestMergeConfigs_BasicFieldOverride(t *testing.T) {
	base := YamlConfig{
		Name: "Logging",
		Path: "Cisco-IOS-XR-um-logging-cfg:/logging",
	}
	override := YamlConfig{
		Version:    "25.4",
		ResDescription: "Updated description",
		DocCategory: "Logging",
	}

	got := mergeConfigs(base, override)

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
}

func TestMergeConfigs_NameOverride(t *testing.T) {
	base := YamlConfig{Name: "Service Timestamps Old", Path: "old-module:/service/timestamps"}
	override := YamlConfig{Version: "25.4", Name: "Service Timestamps", Path: "new-module:/service/timestamps"}

	got := mergeConfigs(base, override)

	if got.Name != "Service Timestamps" {
		t.Errorf("Name: got %q, want %q", got.Name, "Service Timestamps")
	}
	if got.Path != "new-module:/service/timestamps" {
		t.Errorf("Path: got %q, want %q", got.Path, "new-module:/service/timestamps")
	}
}

func TestMergeConfigs_PathUnchanged(t *testing.T) {
	path := "Cisco-IOS-XR-um-logging-cfg:/logging"
	base := YamlConfig{Name: "Logging", Path: path}
	override := YamlConfig{Version: "25.4", ResDescription: "New desc"}

	got := mergeConfigs(base, override)

	if got.Path != path {
		t.Errorf("Path: got %q, want %q (unchanged)", got.Path, path)
	}
}

func TestMergeConfigs_LegacyResource(t *testing.T) {
	base := YamlConfig{
		Name: "Old Feature",
		Path: "old-module:/feature",
		Attributes: []YamlConfigAttribute{
			{YangName: "attr1", TfName: "attr1", Type: "String"},
		},
	}
	override := YamlConfig{Version: "25.4", Legacy: true}

	got := mergeConfigs(base, override)

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
}

func TestMergeConfigs_NoDeletePropagates(t *testing.T) {
	base := YamlConfig{Name: "Feature", NoDelete: false}
	override := YamlConfig{Version: "25.4", NoDelete: true}

	got := mergeConfigs(base, override)

	if !got.NoDelete {
		t.Error("NoDelete: got false, want true")
	}
}

// ---------------------------------------------------------------------------
// mergeAttributes tests
// ---------------------------------------------------------------------------

func TestMergeAttributes_NewAttrGetsAddedInVersion(t *testing.T) {
	base := []YamlConfigAttribute{
		{YangName: "console", TfName: "console", Type: "String"},
	}
	override := []YamlConfigAttribute{
		{YangName: "new-leaf", TfName: "new_leaf", Type: "String"},
	}

	got := mergeAttributes(base, override, "25.4")

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
}

func TestMergeAttributes_LegacyAttrGetsRemovedInVersion(t *testing.T) {
	base := []YamlConfigAttribute{
		{YangName: "source-interface-name", TfName: "name", Type: "String", Id: true},
	}
	override := []YamlConfigAttribute{
		{YangName: "source-interface-name", Legacy: true},
	}

	got := mergeAttributes(base, override, "25.4")

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
}

func TestMergeAttributes_LegacyAttrRenamesTfName(t *testing.T) {
	// F7 fix: legacy block with explicit tf_name renames the base attribute
	// so the natural name is freed for a replacement attribute.
	base := []YamlConfigAttribute{
		{YangName: "console", TfName: "console", Type: "String"},
	}
	override := []YamlConfigAttribute{
		{YangName: "console", TfName: "console_legacy", Legacy: true},
	}

	got := mergeAttributes(base, override, "25.4")

	if len(got) != 1 {
		t.Fatalf("len: got %d, want 1", len(got))
	}
	if got[0].TfName != "console_legacy" {
		t.Errorf("TfName: got %q, want %q", got[0].TfName, "console_legacy")
	}
	if got[0].RemovedInVersion != "25.4" {
		t.Errorf("RemovedInVersion: got %q, want %q", got[0].RemovedInVersion, "25.4")
	}
}

func TestMergeAttributes_LegacyAttrNotInBase_IsSkipped(t *testing.T) {
	// A legacy entry in the override that has no matching base attr is a no-op.
	base := []YamlConfigAttribute{
		{YangName: "existing", TfName: "existing", Type: "String"},
	}
	override := []YamlConfigAttribute{
		{YangName: "ghost", Legacy: true},
	}

	got := mergeAttributes(base, override, "25.4")

	if len(got) != 1 {
		t.Errorf("len: got %d, want 1 (ghost attr should be skipped)", len(got))
	}
}

func TestMergeAttributes_ExistingAttrFieldsUpdated(t *testing.T) {
	base := []YamlConfigAttribute{
		{YangName: "severity", TfName: "severity", Type: "String", Description: "old desc"},
	}
	override := []YamlConfigAttribute{
		{YangName: "severity", Description: "new desc", DefaultValue: "informational"},
	}

	got := mergeAttributes(base, override, "25.4")

	if len(got) != 1 {
		t.Fatalf("len: got %d, want 1", len(got))
	}
	if got[0].Description != "new desc" {
		t.Errorf("Description: got %q, want %q", got[0].Description, "new desc")
	}
	// DefaultValue introduced in an override becomes a versioned default (F17 VersionDefaults).
	// The base had no default (""), so this is a version-scoped change: stored in VersionDefaults,
	// DefaultValue is cleared to "".
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
}

func TestMergeAttributes_MatchByTfName(t *testing.T) {
	// When yang_names differ but tf_names match, merge uses tf_name match.
	base := []YamlConfigAttribute{
		{YangName: "old-yang", TfName: "shared_name", Type: "String"},
	}
	override := []YamlConfigAttribute{
		{YangName: "new-yang", TfName: "shared_name", Type: "Int64"},
	}

	got := mergeAttributes(base, override, "25.4")

	if len(got) != 1 {
		t.Fatalf("len: got %d, want 1 (matched by tf_name)", len(got))
	}
	if got[0].Type != "Int64" {
		t.Errorf("Type: got %q, want %q", got[0].Type, "Int64")
	}
}

func TestMergeAttributes_NewAttrWithNestedGetAddedInVersion(t *testing.T) {
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

func TestMergeAttributes_CompositeKeyVersionedKeys(t *testing.T) {
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

// Scenario 1: base has delete_parent, 25.4 has same yang_name (no replaces_yang_name) and no delete flag.
// Expected: 25.4 inherits delete_parent — no VersionDeleteMode created.
func TestVersionDeleteMode_Scenario1_InheritDeleteParent(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:     "maxfilesize",
		TfName:       "maxfilesize",
		DeleteParent: true,
	}}
	delta := []YamlConfigAttribute{{
		YangName: "maxfilesize",
		TfName:   "maxfilesize",
		// no ReplacesYangName, no delete flag
	}}
	got := mergeAttributes(base, delta, "25.4")
	attr := got[0]
	if attr.VersionDeleteMode != nil {
		t.Errorf("Scenario 1: VersionDeleteMode should be nil (inherited), got %v", attr.VersionDeleteMode)
	}
	if !attr.DeleteParent {
		t.Error("Scenario 1: DeleteParent should remain true after inherit")
	}
}

// Scenario 2: base has delete_grandparent, 25.4 has same yang_name (no replaces_yang_name) and no delete flag.
// Expected: 25.4 inherits delete_grandparent — no VersionDeleteMode created.
func TestVersionDeleteMode_Scenario2_InheritDeleteGrandparent(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:          "maxfilesize",
		TfName:            "maxfilesize",
		DeleteGrandparent: true,
	}}
	delta := []YamlConfigAttribute{{
		YangName: "maxfilesize",
		TfName:   "maxfilesize",
		// no ReplacesYangName, no delete flag
	}}
	got := mergeAttributes(base, delta, "25.4")
	attr := got[0]
	if attr.VersionDeleteMode != nil {
		t.Errorf("Scenario 2: VersionDeleteMode should be nil (inherited), got %v", attr.VersionDeleteMode)
	}
	if !attr.DeleteGrandparent {
		t.Error("Scenario 2: DeleteGrandparent should remain true after inherit")
	}
}

// Scenario 3: base has delete_parent, 25.4 has replaces_yang_name and NO delete flag.
// Expected: version-specific map — 24.4 deletes parent, 25.4 deletes leaf directly.
func TestVersionDeleteMode_Scenario3_ReplacementDropsDeleteParent(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:     "maxfilesize",
		TfName:       "maxfilesize",
		DeleteParent: true,
	}}
	delta := []YamlConfigAttribute{{
		YangName:         "path/maxfilesize",
		TfName:           "maxfilesize",
		ReplacesYangName: "maxfilesize",
		// no delete flag — intentionally drops parent delete
	}}
	got := mergeAttributes(base, delta, "25.4")
	attr := got[0]
	if attr.VersionDeleteMode == nil {
		t.Fatal("Scenario 3: VersionDeleteMode should be set")
	}
	if v, ok := attr.VersionDeleteMode["_base"]; !ok || v != "parent" {
		t.Errorf("Scenario 3: expected _base=parent, got %v", attr.VersionDeleteMode)
	}
	if v, ok := attr.VersionDeleteMode["25.4"]; !ok || v != "" {
		t.Errorf("Scenario 3: expected 25.4='', got %v", attr.VersionDeleteMode)
	}
}

// Scenario 4: base has delete_parent, 25.4 has replaces_yang_name WITH delete_parent.
// Expected: both versions delete parent — no VersionDeleteMode, static path.
func TestVersionDeleteMode_Scenario4_ReplacementKeepsDeleteParent(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:     "maxfilesize",
		TfName:       "maxfilesize",
		DeleteParent: true,
	}}
	delta := []YamlConfigAttribute{{
		YangName:         "path/maxfilesize",
		TfName:           "maxfilesize",
		ReplacesYangName: "maxfilesize",
		DeleteParent:     true,
	}}
	got := mergeAttributes(base, delta, "25.4")
	attr := got[0]
	if attr.VersionDeleteMode != nil {
		t.Errorf("Scenario 4: VersionDeleteMode should be nil (same mode), got %v", attr.VersionDeleteMode)
	}
	if !attr.DeleteParent {
		t.Error("Scenario 4: DeleteParent should be true")
	}
}

// Scenario 5a: base has delete_parent, 25.4 has replaces_yang_name with NO delete flag.
// (Same as 3, delete_grandparent variant for base.)
func TestVersionDeleteMode_Scenario5a_ReplacementDropsDeleteGrandparent(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:          "maxfilesize",
		TfName:            "maxfilesize",
		DeleteGrandparent: true,
	}}
	delta := []YamlConfigAttribute{{
		YangName:         "path/maxfilesize",
		TfName:           "maxfilesize",
		ReplacesYangName: "maxfilesize",
		// no delete flag
	}}
	got := mergeAttributes(base, delta, "25.4")
	attr := got[0]
	if attr.VersionDeleteMode == nil {
		t.Fatal("Scenario 5a: VersionDeleteMode should be set")
	}
	if v := attr.VersionDeleteMode["_base"]; v != "grandparent" {
		t.Errorf("Scenario 5a: expected _base=grandparent, got %q", v)
	}
	if v := attr.VersionDeleteMode["25.4"]; v != "" {
		t.Errorf("Scenario 5a: expected 25.4='', got %q", v)
	}
}

// Scenario 5b: base has delete_grandparent, 25.4 has replaces_yang_name WITH delete_grandparent.
// Expected: no VersionDeleteMode, same mode throughout.
func TestVersionDeleteMode_Scenario5b_ReplacementKeepsDeleteGrandparent(t *testing.T) {
	base := []YamlConfigAttribute{{
		YangName:          "maxfilesize",
		TfName:            "maxfilesize",
		DeleteGrandparent: true,
	}}
	delta := []YamlConfigAttribute{{
		YangName:          "path/maxfilesize",
		TfName:            "maxfilesize",
		ReplacesYangName:  "maxfilesize",
		DeleteGrandparent: true,
	}}
	got := mergeAttributes(base, delta, "25.4")
	attr := got[0]
	if attr.VersionDeleteMode != nil {
		t.Errorf("Scenario 5b: VersionDeleteMode should be nil, got %v", attr.VersionDeleteMode)
	}
	if !attr.DeleteGrandparent {
		t.Error("Scenario 5b: DeleteGrandparent should be true")
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

// TestGetDeletePathExpr_ThreeVersionChain_RenameWithoutModeChange reproduces the read-side
// stale-fallback bug fixed by getVersionValue: a mode change recorded at 25.4, followed by a
// 26.2 rename that restates the same mode (correctly producing no new VersionDeleteMode
// entry), used to bake the wrong, stale oldest mode into 26.2's generated path.
func TestGetDeletePathExpr_ThreeVersionChain_RenameWithoutModeChange(t *testing.T) {
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
