package env

import "testing"

var (
	testEnviron = []string{
		"KEY1=val1",
		"KEY2=val2",
		"BOOL1=true",
		"BOOL2=f",
		"FLOAT1=3.14",
		"INT1=-1",
		"UINT1=69",
	}
	testEnv = LoadEnvArrayString(testEnviron)
)

func TestLoadEnvArrayString(t *testing.T) {
	if result := testEnv["KEY1"]; result != "val1" {
		t.Error(`Expected "val1", got`, result)
	}
}

func TestLoadEnvArrayString_safety(t *testing.T) {
	// TODO
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

func TestGetFloat(t *testing.T) {
	result, err := testEnv.GetFloat("FLOAT1", 64)
	if err != nil {
		t.Error(err)
	}
	if result != 3.14 {
		t.Error(`Expected 3.14, got`, result)
	}
}

func TestGetFloat_empty(t *testing.T) {
	result, err := testEnv.GetFloat("FLOAT2", 64)
	if err != nil {
		t.Error(err)
	}
	if result != 0 {
		t.Error(`Expected 0, got`, result)
	}
}

func TestGetInt(t *testing.T) {
	result, err := testEnv.GetInt("INT1", 10, 64)
	if err != nil {
		t.Error(err)
	}
	if result != -1 {
		t.Error(`Expected -1, got`, result)
	}
}

func TestGetInt_empty(t *testing.T) {
	result, err := testEnv.GetInt("INT2", 10, 64)
	if err != nil {
		t.Error(err)
	}
	if result != 0 {
		t.Error(`Expected 0, got`, result)
	}
}
