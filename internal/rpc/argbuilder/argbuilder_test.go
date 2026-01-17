package argbuilder

import (
	"reflect"
	"testing"

	pb "github.com/tmc/nlm/gen/notebooklm/v1alpha1"
)

func TestEncodeArgs(t *testing.T) {
	enc := NewArgumentEncoder()
	msg := &pb.Project{
		Title:     "My Project",
		ProjectId: "proj-1",
	}

	args, err := enc.EncodeArgs(msg, `[%title%, %project_id%, null, [%project_id%], [1,2], 42]`)
	if err != nil {
		t.Fatalf("EncodeArgs error: %v", err)
	}

	expected := []interface{}{
		"My Project",
		"proj-1",
		nil,
		"proj-1",
		[]interface{}{1, 2},
		42,
	}

	if !reflect.DeepEqual(args, expected) {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestSplitFormatNested(t *testing.T) {
	enc := NewArgumentEncoder()
	parts := enc.splitFormat(`%title%, [1,2,[3,4]], null`)
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}
}

func TestParseLiteral(t *testing.T) {
	enc := NewArgumentEncoder()
	if got := enc.parseLiteral("123"); got != 123 {
		t.Fatalf("expected 123, got %#v", got)
	}
	if got := enc.parseLiteral(`"hello"`); got != "hello" {
		t.Fatalf("expected hello, got %#v", got)
	}
}

func TestSnakeToCamel(t *testing.T) {
	if got := snakeToCamel("project_id"); got != "projectId" {
		t.Fatalf("expected projectId, got %q", got)
	}
}
