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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bgpASFormatModel mirrors the generated BGPASFormat struct for test purposes.
type bgpASFormatModel struct {
	Device     types.String `tfsdk:"device"`
	Id         types.String `tfsdk:"id"`
	DeleteMode types.String `tfsdk:"delete_mode"`
	AsFormat   types.String `tfsdk:"as_format"`
}

// bgpASFormatConstraints mirrors the GetEnumConstraints() output for bgp_as_format
// after the 25.4 test override: 24.4={asdot,asplain}, 25.4={asdot,asdot+}.
var bgpASFormatConstraints = []helpers.FieldEnumConstraint{
	{
		FieldPath: "as_format",
		VersionEnums: map[string][]string{
			"24.4": {"asdot", "asplain"},
			"25.4": {"asdot", "asdot+"},
		},
	},
}

func TestValidateVersionEnums_24_4_allowsAsdot(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asdot")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("24.4", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error for asdot on 24.4, got: %v", diags)
	}
}

func TestValidateVersionEnums_24_4_allowsAsplain(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asplain")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("24.4", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error for asplain on 24.4, got: %v", diags)
	}
}

func TestValidateVersionEnums_24_4_rejectsAsdotPlus(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asdot+")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("24.4", plan, bgpASFormatConstraints, &diags)
	if !diags.HasError() {
		t.Error("expected error for asdot+ on 24.4 (not in 24.4 set), got none")
	}
}

func TestValidateVersionEnums_25_4_allowsAsdot(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asdot")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("25.4", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error for asdot on 25.4, got: %v", diags)
	}
}

func TestValidateVersionEnums_25_4_allowsAsdotPlus(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asdot+")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("25.4", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error for asdot+ on 25.4, got: %v", diags)
	}
}

func TestValidateVersionEnums_25_4_rejectsAsplain(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asplain")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("25.4", plan, bgpASFormatConstraints, &diags)
	if !diags.HasError() {
		t.Error("expected error for asplain on 25.4 (removed in 25.4), got none")
	}
}

func TestValidateVersionEnums_nullFieldSkipped(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringNull()}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("25.4", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error for null field, got: %v", diags)
	}
}

func TestValidateVersionEnums_emptyVersionSkipped(t *testing.T) {
	plan := bgpASFormatModel{AsFormat: types.StringValue("asplain")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error when version is empty, got: %v", diags)
	}
}

func TestValidateVersionEnums_versionBelowAllThresholdsSkipped(t *testing.T) {
	// Version 24.3 is below both 24.4 and 25.4 thresholds — no restriction applies.
	plan := bgpASFormatModel{AsFormat: types.StringValue("asplain")}
	var diags diag.Diagnostics
	helpers.ValidateVersionEnums("24.3", plan, bgpASFormatConstraints, &diags)
	if diags.HasError() {
		t.Errorf("expected no error for version below all thresholds, got: %v", diags)
	}
}
