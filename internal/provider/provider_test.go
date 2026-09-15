// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-iosxr/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"iosxr": providerserver.NewProtocol6WithError(New()),
	}
)

// iosxrVersionAtLeast returns true when currentVersion meets minVersion, or when
// currentVersion is empty (permissive default so tests aren't skipped without a version set).
func iosxrVersionAtLeast(currentVersion, minVersion string) bool {
	return helpers.VersionAtLeast(currentVersion, minVersion)
}

// selectVersionExample returns the example value appropriate for the IOSXR_VERSION in the
// environment. Keys in byVersion are version thresholds (e.g. "25.4"); the highest threshold
// satisfied by IOSXR_VERSION wins. Falls back to baseExample when no threshold matches or
// IOSXR_VERSION is unset.
func selectVersionExample(byVersion map[string]string, baseExample string) string {
	ver := os.Getenv("IOSXR_VERSION")
	if ver == "" || len(byVersion) == 0 {
		return baseExample
	}
	bestVer := ""
	for v := range byVersion {
		if iosxrVersionAtLeast(ver, v) && (bestVer == "" || iosxrVersionAtLeast(v, bestVer)) {
			bestVer = v
		}
	}
	if bestVer == "" {
		return baseExample
	}
	return byVersion[bestVer]
}

// selectVersionTestTags returns the test tag set appropriate for the IOSXR_VERSION in the
// environment. Keys in byVersion are version thresholds; the highest threshold satisfied by
// IOSXR_VERSION wins. Falls back to baseTags when no threshold matches or IOSXR_VERSION is unset.
func selectVersionTestTags(byVersion map[string][]string, baseTags []string) []string {
	ver := os.Getenv("IOSXR_VERSION")
	if ver == "" || len(byVersion) == 0 {
		return baseTags
	}
	bestVer := ""
	for v := range byVersion {
		if iosxrVersionAtLeast(ver, v) && (bestVer == "" || iosxrVersionAtLeast(v, bestVer)) {
			bestVer = v
		}
	}
	if bestVer == "" {
		return baseTags
	}
	return byVersion[bestVer]
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("IOSXR_USERNAME"); v == "" {
		t.Fatal("IOSXR_USERNAME env variable must be set for acceptance tests")
	}
	if v := os.Getenv("IOSXR_PASSWORD"); v == "" {
		t.Fatal("IOSXR_PASSWORD env variable must be set for acceptance tests")
	}
	if v := os.Getenv("IOSXR_HOST"); v == "" {
		t.Fatal("IOSXR_HOST env variable must be set for acceptance tests")
	}
}
