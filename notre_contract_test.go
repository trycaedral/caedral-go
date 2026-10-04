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
	// V1 base shape: no economy fields.
	if resp.Notre.InputSaved != nil || resp.Notre.Result != "" {
		t.Fatalf("unexpected economy on V1 fixture: %+v", resp.Notre)
	}
}

func TestResponseTelemetryMetadataV2(t *testing.T) {
	raw := loadFixture(t, "notre-response-telemetry-v2.json")
	var resp ChatCompletion
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Notre == nil {
		t.Fatal("expected notre metadata")
	}
	if !resp.Notre.Enabled || !resp.Notre.Intervened || resp.Notre.FallbackUsed {
		t.Fatalf("unexpected metadata: %+v", resp.Notre)
	}
	if resp.Notre.InputBefore == nil || *resp.Notre.InputBefore != 1200 {
		t.Fatalf("input_before: %+v", resp.Notre)
	}
	if resp.Notre.InputSent == nil || *resp.Notre.InputSent != 310 {
		t.Fatalf("input_sent: %+v", resp.Notre)
	}
	if resp.Notre.InputSaved == nil || *resp.Notre.InputSaved != 890 {
		t.Fatalf("input_saved: %+v", resp.Notre)
	}
	if resp.Notre.Result != "optimized" {
		t.Fatalf("result: %+v", resp.Notre)
	}
	if resp.Notre.ValueUSD == nil || *resp.Notre.ValueUSD < 0.0026 || *resp.Notre.ValueUSD > 0.0028 {
		t.Fatalf("value_usd: %+v", resp.Notre)
	}
}

func TestResponseTelemetryMetadataV3(t *testing.T) {
	raw := loadFixture(t, "notre-response-telemetry-v3.json")
	var resp ChatCompletion
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Notre == nil {
		t.Fatal("expected notre metadata")
	}
	if resp.Notre.Shape != "chat" {
		t.Fatalf("shape: %+v", resp.Notre)
	}
	if resp.Notre.ContractVersion == nil || *resp.Notre.ContractVersion != 3 {
		t.Fatalf("contract_version: %+v", resp.Notre)
	}
	if resp.Notre.SavedBreakdown == nil {
		t.Fatal("expected saved_breakdown")
	}
	if resp.Notre.SavedBreakdown.CacheHitTokens == nil || *resp.Notre.SavedBreakdown.CacheHitTokens != 640 {
		t.Fatalf("cache_hit_tokens: %+v", resp.Notre.SavedBreakdown)
	}
	if resp.Notre.SavedBreakdown.DedupTokens == nil || *resp.Notre.SavedBreakdown.DedupTokens != 20 {
		t.Fatalf("dedup_tokens: %+v", resp.Notre.SavedBreakdown)
	}
	if resp.Notre.SavedBreakdown.PrefilterTokens == nil || *resp.Notre.SavedBreakdown.PrefilterTokens != 0 {
		t.Fatalf("prefilter_tokens: %+v", resp.Notre.SavedBreakdown)
	}
	if resp.Notre.InputSaved == nil || resp.Notre.InputBefore == nil || resp.Notre.InputSent == nil {
		t.Fatalf("economy fields: %+v", resp.Notre)
	}
	if *resp.Notre.InputSaved != *resp.Notre.InputBefore-*resp.Notre.InputSent {
		t.Fatalf("saved invariant: %+v", resp.Notre)
	}
}
