package caedral

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func fixturePath(name string) string {
	return filepath.Join("tests", "fixtures", "notre-contract", name)
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(fixturePath(name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func TestNotreAutoTelemetryRequestParity(t *testing.T) {
	raw := loadFixture(t, "notre-auto-telemetry.json")
	var req ChatCompletionRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Notre == nil || req.Notre.Mode != "auto" || req.Notre.Telemetry == nil || !*req.Notre.Telemetry {
		t.Fatalf("unexpected notre: %+v", req.Notre)
	}
	out, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var want map[string]any
	var got map[string]any
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("json drift:\nwant %#v\ngot %#v", want, got)
	}
}

func TestNotreOmittedBackwardCompatible(t *testing.T) {
	raw := loadFixture(t, "notre-omitted.json")
	var req ChatCompletionRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Notre != nil {
		t.Fatalf("expected nil notre")
	}
}

func TestResponseTelemetryMetadataV1(t *testing.T) {
	raw := loadFixture(t, "notre-response-telemetry.json")
	var resp ChatCompletion
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Notre == nil {
		t.Fatal("expected notre metadata")
	}
	if !resp.Notre.Enabled || resp.Notre.Mode != "auto" || resp.Notre.Intervened || resp.Notre.FallbackUsed {
		t.Fatalf("unexpected metadata: %+v", resp.Notre)
	}
}
