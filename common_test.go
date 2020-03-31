package ofxgo

import (
	"testing"
)

func TestStatusValid(t *testing.T) {
	s := Status{
		Code:     0,
		Severity: "INFO",
		Message:  "Success",
	}
	if ok, err := s.Valid(); !ok {
		t.Fatalf("Status unexpectedly invalid: %s\n", err)
	}

	s.Severity = "INVALID"
	if ok, err := s.Valid(); ok || err == nil {
		t.Fatalf("Status unexpectedly valid invalid Severity\n")
	}

	s.Severity = "WARN"
	if ok, err := s.Valid(); ok || err == nil {
		t.Fatalf("Status unexpectedly valid with wrong Severity for Code 0\n")
	}

	s.Code = 9
	if ok, err := s.Valid(); ok || err == nil {
		t.Fatalf("Status unexpectedly valid with invalid Code\n")
	}
}

func TestStatusCodeMeaning(t *testing.T) {
	s := Status{
		Code:     15500,
		Severity: "ERROR",
	}
	meaning, err := s.CodeMeaning()
	if err != nil {
		t.Fatalf("Status.CodeMeaning unexpectedly failed: %s\n", err)
	}
	if meaning != "Signon invalid" {
		t.Fatalf("Unexpected meaning for Code 15500: \"%s\"\n", meaning)
	}

	s.Code = 999
	if meaning, err := s.CodeMeaning(); len(meaning) != 0 || err == nil {
		t.Fatalf("Status.CodeMeaning unexpectedly succeeded with invalid Code\n")
	}
}

func TestStatusCodeConditions(t *testing.T) {
	s := Status{
		Code:     2006,
		Severity: "ERROR",
	}
	if conditions, err := s.CodeConditions(); len(conditions) == 0 || err != nil {
		t.Fatalf("Status.CodeConditions unexpectedly failed: %s\n", err)
	}

	s.Code = 999
	if conditions, err := s.CodeConditions(); len(conditions) != 0 || err == nil {
		t.Fatalf("Status.CodeConditions unexpectedly succeeded with invalid Code\n")
	}
}
