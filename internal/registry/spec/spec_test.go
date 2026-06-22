package spec

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestValidateAcceptsValidToolSpec(t *testing.T) {
	tool := validToolSpec()
	tool.Metadata = map[string]string{
		"owner": "",
		"team":  "platform",
	}

	if err := Validate(tool); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidateRequiresTopLevelFields(t *testing.T) {
	tool := validToolSpec()
	tool.ID = ""
	tool.Name = ""
	tool.Version = ""
	tool.Status = ""

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "id", "name", "version", "status")
}

func TestValidateRejectsInvalidName(t *testing.T) {
	tool := validToolSpec()
	tool.Name = "Bad.Name"

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "name")
}

func TestValidateRejectsInvalidVersion(t *testing.T) {
	tool := validToolSpec()
	tool.Version = "1.0 beta"

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "version")
}

func TestValidateRejectsInvalidStatus(t *testing.T) {
	tool := validToolSpec()
	tool.Status = "archived"

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "status")
}

func TestValidateRejectsEmptyActions(t *testing.T) {
	tool := validToolSpec()
	tool.Actions = nil

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "actions")
}

func TestValidateRejectsDuplicateActionNames(t *testing.T) {
	tool := validToolSpec()
	tool.Actions = []ToolAction{
		{Name: "search", Description: "Search records."},
		{Name: "search", Description: "Search records again."},
	}

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "actions[1].name")

	if !strings.Contains(err.Error(), "duplicates actions[0].name") {
		t.Fatalf("Validate() error = %q, want duplicate action detail", err.Error())
	}
}

func TestValidateRequiresActionFields(t *testing.T) {
	tool := validToolSpec()
	tool.Actions = []ToolAction{{DisplayName: "Search"}}

	err := mustValidationError(t, Validate(tool))
	wantFields(t, err, "actions[0].name", "actions[0].description")
}

func TestValidateMetadataBoundaryCases(t *testing.T) {
	t.Run("nil metadata is valid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = nil

		if err := Validate(tool); err != nil {
			t.Fatalf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("empty metadata is valid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = map[string]string{}

		if err := Validate(tool); err != nil {
			t.Fatalf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("max metadata entries is valid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = metadataEntries(MaxMetadataEntries)

		if err := Validate(tool); err != nil {
			t.Fatalf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("too many metadata entries is invalid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = metadataEntries(MaxMetadataEntries + 1)

		err := mustValidationError(t, Validate(tool))
		wantFields(t, err, "metadata")
	})

	t.Run("blank metadata key is invalid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = map[string]string{" ": "platform"}

		err := mustValidationError(t, Validate(tool))
		wantFields(t, err, `metadata[" "]`)
	})

	t.Run("metadata key with surrounding whitespace is invalid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = map[string]string{" owner": "platform"}

		err := mustValidationError(t, Validate(tool))
		wantFields(t, err, `metadata[" owner"]`)
	})

	t.Run("metadata key over max length is invalid", func(t *testing.T) {
		tool := validToolSpec()
		tool.Metadata = map[string]string{strings.Repeat("a", MaxMetadataKeyLength+1): "platform"}

		err := mustValidationError(t, Validate(tool))
		wantFieldPrefix(t, err, "metadata[")
	})
}

func validToolSpec() ToolSpec {
	return ToolSpec{
		ID:          "tool-search",
		Name:        "search",
		DisplayName: "Search",
		Description: "Search known records.",
		Version:     "1.0.0",
		Status:      ToolStatusDraft,
		Tags:        []string{"records"},
		Actions: []ToolAction{
			{Name: "search", DisplayName: "Search", Description: "Search records."},
		},
	}
}

func metadataEntries(count int) map[string]string {
	metadata := make(map[string]string, count)
	for i := 0; i < count; i++ {
		metadata[fmt.Sprintf("key_%02d", i)] = "value"
	}

	return metadata
}

func mustValidationError(t *testing.T, err error) *ValidationError {
	t.Helper()

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("Validate() error = %v, want *ValidationError", err)
	}
	if len(validationErr.Fields) == 0 {
		t.Fatalf("Validate() fields are empty, want at least one field error")
	}

	return validationErr
}

func wantFields(t *testing.T, err *ValidationError, fields ...string) {
	t.Helper()

	got := make(map[string]bool, len(err.Fields))
	for _, field := range err.Fields {
		got[field.Field] = true
	}

	for _, field := range fields {
		if !got[field] {
			t.Fatalf("Validate() fields = %#v, want field %q", err.Fields, field)
		}
	}
}

func wantFieldPrefix(t *testing.T, err *ValidationError, prefix string) {
	t.Helper()

	for _, field := range err.Fields {
		if strings.HasPrefix(field.Field, prefix) {
			return
		}
	}

	t.Fatalf("Validate() fields = %#v, want field prefix %q", err.Fields, prefix)
}
