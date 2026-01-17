package beprotojson

import (
	"testing"

	pb "github.com/tmc/nlm/gen/notebooklm/v1alpha1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUnmarshalWrapperBoolFromString(t *testing.T) {
	msg := &wrapperspb.BoolValue{}
	if err := Unmarshal([]byte(`["true"]`), msg); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if msg.Value != true {
		t.Fatalf("expected true, got %v", msg.Value)
	}
}

func TestUnmarshalWrapperInt32(t *testing.T) {
	msg := &wrapperspb.Int32Value{}
	if err := Unmarshal([]byte(`[123]`), msg); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if msg.Value != 123 {
		t.Fatalf("expected 123, got %d", msg.Value)
	}
}

func TestUnmarshalWrapperStringFromArray(t *testing.T) {
	msg := &wrapperspb.StringValue{}
	if err := Unmarshal([]byte(`[["hi"]]`), msg); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if msg.Value != "hi" {
		t.Fatalf("expected hi, got %q", msg.Value)
	}
}

func TestUnmarshalEnumFromString(t *testing.T) {
	msg := &pb.SourceSettings{}
	if err := Unmarshal([]byte(`[null,"SOURCE_STATUS_DISABLED"]`), msg); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if msg.Status != pb.SourceSettings_SOURCE_STATUS_DISABLED {
		t.Fatalf("expected disabled, got %v", msg.Status)
	}
}
