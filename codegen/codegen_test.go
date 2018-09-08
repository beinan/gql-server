package codegen

import "testing"

func TestGenerate(t *testing.T) {
	err := Generate(GenConfig{})
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}
