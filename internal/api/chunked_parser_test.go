package api

import "testing"

const testProjectID = "123e4567-e89b-12d3-a456-426614174000"

func TestParseStandardJSON(t *testing.T) {
	projectsEscaped := `[[\\\"Alpha\\\",null,\\\"` + testProjectID + `\\\",\\\"📄\\\"]]`
	jsonChunk := `["wrb.fr","wXbhsf","` + projectsEscaped + `"]`
	raw := ")]}'\n" + jsonChunk + "\n25"
	parser := NewChunkedResponseParser(raw)
	parser.rawChunks = parser.extractChunks()

	projects, err := parser.parseStandardJSON()
	if err != nil {
		t.Fatalf("parseStandardJSON: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	got := projects[0]
	if got.Title != "Alpha" || got.ProjectId != testProjectID || got.Emoji != "📄" {
		t.Fatalf("unexpected project: %#v", got)
	}
}

func TestParseAsObject(t *testing.T) {
	parser := NewChunkedResponseParser("")
	data := `{"` + testProjectID + `":{"title":"Alpha","emoji":"🔥"}}`

	projects, err := parser.parseAsObject(data)
	if err != nil {
		t.Fatalf("parseAsObject: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	got := projects[0]
	if got.ProjectId != testProjectID || got.Title != "Alpha" || got.Emoji != "🔥" {
		t.Fatalf("unexpected project: %#v", got)
	}
}

func TestParseDirectScan(t *testing.T) {
	parser := NewChunkedResponseParser("")
	parser.cleanedData = `before "My Project" ` + testProjectID + ` "🚀"`

	projects, err := parser.parseDirectScan()
	if err != nil {
		t.Fatalf("parseDirectScan: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	got := projects[0]
	if got.ProjectId != testProjectID || got.Title != "My Project" || got.Emoji != "🚀" {
		t.Fatalf("unexpected project: %#v", got)
	}
}

func TestParseJSONArrayFallback(t *testing.T) {
	parser := NewChunkedResponseParser(")]}'\n{\"a\":1}\n[[\"b\"]]\n")

	result, err := parser.ParseJSONArray()
	if err != nil {
		t.Fatalf("ParseJSONArray: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
}

func TestSanitizeResponse(t *testing.T) {
	parser := NewChunkedResponseParser("")
	raw := ")]}'\n12\n[[\"a\"\n]\n]\n25\n"

	if got := parser.SanitizeResponse(raw); got != `[[\"a\"]]` && got != "[[\"a\"]]" {
		t.Fatalf("unexpected sanitized response: %q", got)
	}
}

func TestBalancedBrackets(t *testing.T) {
	if !balancedBrackets(`{"a":[1,2]}`) {
		t.Fatalf("expected balanced brackets")
	}
	if balancedBrackets(`{"a":[1,2}`) {
		t.Fatalf("expected unbalanced brackets")
	}
}
