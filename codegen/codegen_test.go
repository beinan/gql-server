package codegen

import (
	"os"
	"testing"
)

func TestGenerateModel(t *testing.T) {
	err := GenerateModel(GenConfig{
		SchemaPath: "../example/schema",
	}, os.Stdout)
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}

func TestGenerateResolver(t *testing.T) {
	err := GenerateResolver(GenConfig{
		SchemaPath: "../example/schema",
	}, os.Stdout)
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}

func TestGenerateGqlResolver(t *testing.T) {
	err := GenerateGqlResolver(GenConfig{
		SchemaPath: "../example/schema",
	}, os.Stdout)
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}

func TestLoadSchema(t *testing.T) {
	loadSchema("../example/schema")
}
