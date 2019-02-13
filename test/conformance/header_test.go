// +build e2e

/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conformance

import (
	"net"
	"regexp"
	"strings"
	"testing"

	"github.com/knative/pkg/test/logging"
	"github.com/knative/serving/test"
)

// TestMustHaveHeadersSet verified that all headers declared as "MUST" in the runtime
// contract are present from the point of view of the user container.
func TestMustHaveHeadersSet(t *testing.T) {
	logger := logging.GetContextLogger(t.Name())
	clients := setup(t)

	ri, err := fetchRuntimeInfo(clients, logger, &test.Options{})
	if err != nil {
		t.Fatalf("Error fetching runtime info: %v", err)
	}

	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map. Therefore we
	// check against the Host field instead of the map.
	if ri.Request.Host == "" {
		// We just check that the host header exists and is non-empty as the request
		// may be made internally or externally which will result in a different host.
		t.Error("Header host was not present on request")
	}

	// TODO(#3112): Validate Forwarded header once it is enabled.
}

// TestMustHaveHeadersSet verified that all headers declared as "SHOULD" in the runtime
// contract are present from the point of view of the user container.
func TestShouldHaveHeadersSet(t *testing.T) {
	logger := logging.GetContextLogger(t.Name())
	clients := setup(t)

	expectedHeaders := map[string]*regexp.Regexp{
		// We expect the protocol to be http for our test image.
		"x-forwarded-proto": regexp.MustCompile("https?"),
		// We expect the value to be a list of at least one comma separated IP addresses (IPv4 or IPv6).
		"x-forwarded-for": nil, // Non-regex based validation performed for this header

		// Trace Headers
		// See https://github.com/openzipkin/b3-propagation#overall-process
		// We use the multiple header variant for tracing. We do not validate the single header variant.
		// We expect the value to be a 64-bit hex string
		"x-b3-spanid": regexp.MustCompile("[0-9a-f]{16}"),
		// We expect the value to be a 64-bit or 128-bit hex string
		"x-b3-traceid": regexp.MustCompile("[0-9a-f]{16}|[0-9a-f]{32}"),

		// "x-b3-parentspanid" and "x-b3-sampled" are often present for tracing, but are not
		// required for tracing so we do not validate them.
	}

	ri, err := fetchRuntimeInfo(clients, logger, &test.Options{})
	if err != nil {
		t.Fatalf("Error fetching runtime info: %v", err)
	}

	headers := ri.Request.Headers

	for header, regex := range expectedHeaders {
		hv, ok := headers[header]
		if !ok {
			t.Errorf("Header %s was not present on request", header)
			continue
		}

		switch {
		case strings.EqualFold(header, "x-forwarded-for"):
			for _, ip := range strings.Split(hv, ",") {
				if net.ParseIP(strings.TrimSpace(ip)) == nil {
					t.Errorf("Header %s has invalid IP: %s", header, ip)
				}
			}
		default:
			if !regex.MatchString(hv) {
				t.Errorf("%s = %s; want: %s", header, hv, regex.String())
			}
		}
	}
}
