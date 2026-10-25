// Package contract checks that the implemented routes match docs/openapi.yaml.
package contract

import (
	"os"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/AknurOE/sen-support-platform/internal/adapter/httpapi"
)

func TestRoutesMatchOpenAPI(t *testing.T) {
	raw, err := os.ReadFile("../../../docs/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	spec := map[string]bool{}
	for p, ops := range doc.Paths {
		for m, op := range ops {
			o, ok := op.(map[string]any)
			if !ok || m == "parameters" {
				continue
			}
			if o["x-status"] == "planned" {
				continue
			}
			spec[strings.ToUpper(m)+" /api/v1"+p] = true
		}
	}
	impl := map[string]bool{}
	for _, r := range httpapi.Routes() {
		impl[r] = true
	}
	var missing, extra []string
	for k := range spec {
		if !impl[k] {
			missing = append(missing, k)
		}
	}
	for k := range impl {
		if !spec[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	if len(missing)+len(extra) > 0 {
		t.Fatalf("contract drift\nin openapi but not implemented: %v\nimplemented but not in openapi: %v", missing, extra)
	}
	t.Logf("%d operations match the contract", len(spec))
}
