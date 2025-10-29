package docs

import "testing"

func TestDocTemplateNotEmpty(t *testing.T) {
	if len(docTemplate) == 0 {
		t.Fatal("docTemplate should not be empty")
	}
}
