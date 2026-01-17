package api

import (
	"encoding/json"
	"testing"

	pb "github.com/tmc/nlm/gen/notebooklm/v1alpha1"
)

func TestExtractSourceIDFormats(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "format1", raw: `[[[["id1"]]]]`, want: "id1"},
		{name: "format2", raw: `[[["id2"]]]`, want: "id2"},
		{name: "format3", raw: `[["id3"]]`, want: "id3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := extractSourceID(json.RawMessage(tt.raw))
			if err != nil {
				t.Fatalf("extractSourceID error: %v", err)
			}
			if id != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, id)
			}
		})
	}
}

func TestExtractSourceIDRegisterResponse(t *testing.T) {
	raw := `[[[[["sid"],"file",[null,null,null,null,0]]]]]`
	id, err := extractSourceIDFromRegisterResponse(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("extractSourceIDFromRegisterResponse error: %v", err)
	}
	if id != "sid" {
		t.Fatalf("expected sid, got %q", id)
	}
}

func TestParseArtifactFromResponse(t *testing.T) {
	client := &Client{}
	data := []interface{}{
		"artifact-1",
		float64(pb.ArtifactType_ARTIFACT_TYPE_NOTE),
		float64(pb.ArtifactState_ARTIFACT_STATE_READY),
		[]interface{}{"src1", "src2"},
	}

	artifact := client.parseArtifactFromResponse(data)
	if artifact == nil {
		t.Fatalf("expected artifact")
	}
	if artifact.ArtifactId != "artifact-1" {
		t.Fatalf("unexpected id: %q", artifact.ArtifactId)
	}
	if len(artifact.Sources) != 2 {
		t.Fatalf("expected 2 sources")
	}
}
