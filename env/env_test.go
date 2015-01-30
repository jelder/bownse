package env

import "testing"

var (
	testEnviron = []string{
		"KEY1=val1",
		"KEY2=val2",
		"BOOL1=true",
		"BOOL2=f",
		"FLOAT1=-3.14",
	}
	testEnv = MustLoadEnvArrayString(testEnviron)
)

func TestLoadEnvArrayString(t *testing.T) {
	if result := testEnv["KEY1"]; result != "val1" {
		t.Error(`Expected "val1", got`, result)
	}
}

func TestLoadEnvArrayString_safety(t *testing.T) {
	_, err := LoadEnvArrayString([]string{"OK=true", "Ok=false"})
	if err != nil {
		t.Error(`Expected error, got nil`)
	}
}

func TestGetBool(t *testing.T) {
	if result := testEnv.GetBool("BOOL1"); !result {
		t.Error(`Expected true, got`, result)
	}
	if result := testEnv.GetBool("BOOL2"); result {
		t.Error(`Expected false, got`, result)
	}
	if result := testEnv.GetBool("BOOL99"); result {
		t.Error(`Expected false, got`, result)
	}
}

func TestGetNumber(t *testing.T) {
	result := testEnv.GetNumber("FLOAT1", 0)
	if result != -3.14 {
		t.Error(`Expected 3.14, got`, result)
	}
}

func TestGetNumber_empty(t *testing.T) {
	result := testEnv.GetNumber("FLOAT2", 0)
	if result != 0 {
		t.Error(`Expected 0, got`, result)
	}
}
