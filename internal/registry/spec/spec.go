package spec

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type ToolID string

type ToolStatus string

const (
	ToolStatusDraft      ToolStatus = "draft"
	ToolStatusActive     ToolStatus = "active"
	ToolStatusDeprecated ToolStatus = "deprecated"
	ToolStatusDisabled   ToolStatus = "disabled"
)

const (
	MaxMetadataEntries   = 32
	MaxMetadataKeyLength = 64
)

var machineNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,62}$`)

type ToolSpec struct {
	ID          ToolID
	Name        string
	DisplayName string
	Description string
	Version     string
	Status      ToolStatus
	Tags        []string
	Actions     []ToolAction
	Metadata    map[string]string
}

type ToolAction struct {
	Name        string
	DisplayName string
	Description string
}

type FieldError struct {
	Field   string
	Message string
}

type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	if e == nil || len(e.Fields) == 0 {
		return "invalid tool spec"
	}

	parts := make([]string, 0, len(e.Fields))
	for _, field := range e.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", field.Field, field.Message))
	}

	return "invalid tool spec: " + strings.Join(parts, "; ")
}

func Validate(tool ToolSpec) error {
	var fields []FieldError

	if strings.TrimSpace(string(tool.ID)) == "" {
		fields = append(fields, FieldError{Field: "id", Message: "is required"})
	}

	if strings.TrimSpace(tool.Name) == "" {
		fields = append(fields, FieldError{Field: "name", Message: "is required"})
	} else if !isMachineName(tool.Name) {
		fields = append(fields, FieldError{Field: "name", Message: "must start with a lowercase letter and contain only lowercase letters, digits, underscores, or hyphens"})
	}

	if strings.TrimSpace(tool.Version) == "" {
		fields = append(fields, FieldError{Field: "version", Message: "is required"})
	} else if !isOpaqueVersion(tool.Version) {
		fields = append(fields, FieldError{Field: "version", Message: "must not contain whitespace or control characters"})
	}

	if tool.Status == "" {
		fields = append(fields, FieldError{Field: "status", Message: "is required"})
	} else if !isValidStatus(tool.Status) {
		fields = append(fields, FieldError{Field: "status", Message: "must be one of draft, active, deprecated, disabled"})
	}

	fields = append(fields, validateActions(tool.Actions)...)
	fields = append(fields, validateMetadata(tool.Metadata)...)

	if len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}

	return nil
}

func isValidStatus(status ToolStatus) bool {
	switch status {
	case ToolStatusDraft, ToolStatusActive, ToolStatusDeprecated, ToolStatusDisabled:
		return true
	default:
		return false
	}
}

func isMachineName(name string) bool {
	return machineNamePattern.MatchString(name)
}

func isOpaqueVersion(version string) bool {
	for _, r := range version {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return false
		}
	}

	return true
}

func validateActions(actions []ToolAction) []FieldError {
	if len(actions) == 0 {
		return []FieldError{{Field: "actions", Message: "must contain at least one action"}}
	}

	var fields []FieldError
	seen := make(map[string]int, len(actions))
	for i, action := range actions {
		nameField := fmt.Sprintf("actions[%d].name", i)
		descriptionField := fmt.Sprintf("actions[%d].description", i)

		if strings.TrimSpace(action.Name) == "" {
			fields = append(fields, FieldError{Field: nameField, Message: "is required"})
		} else if !isMachineName(action.Name) {
			fields = append(fields, FieldError{Field: nameField, Message: "must start with a lowercase letter and contain only lowercase letters, digits, underscores, or hyphens"})
		} else if firstIndex, ok := seen[action.Name]; ok {
			fields = append(fields, FieldError{Field: nameField, Message: fmt.Sprintf("duplicates actions[%d].name", firstIndex)})
		} else {
			seen[action.Name] = i
		}

		if strings.TrimSpace(action.Description) == "" {
			fields = append(fields, FieldError{Field: descriptionField, Message: "is required"})
		}
	}

	return fields
}

func validateMetadata(metadata map[string]string) []FieldError {
	if len(metadata) == 0 {
		return nil
	}

	var fields []FieldError
	if len(metadata) > MaxMetadataEntries {
		fields = append(fields, FieldError{Field: "metadata", Message: fmt.Sprintf("must contain at most %d entries", MaxMetadataEntries)})
	}

	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		field := fmt.Sprintf("metadata[%q]", key)
		if strings.TrimSpace(key) == "" {
			fields = append(fields, FieldError{Field: field, Message: "key is required"})
			continue
		}
		if strings.TrimSpace(key) != key {
			fields = append(fields, FieldError{Field: field, Message: "key must not have surrounding whitespace"})
		}
		if len(key) > MaxMetadataKeyLength {
			fields = append(fields, FieldError{Field: field, Message: fmt.Sprintf("key must be at most %d bytes", MaxMetadataKeyLength)})
		}
	}

	return fields
}
