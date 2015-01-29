package main

import "testing"

func TestCheckEnvUrl_junk(t *testing.T) {
	result := CheckEnvUrl([]string{`PATH`, `bin:bin/bash`})
	if result != "" {
		t.Error(`Expected "", got`, result)
	}
}

func TestCheckEnvUrl_good_key_bad_url(t *testing.T) {

	result := CheckEnvUrl([]string{`CHAT_URL`, `bin:bin/bash`})
	if result != "" {
		t.Error(`Expected "", got`, result)
	}
}

func TestCheckEnvUrl_bad_key_good_url(t *testing.T) {
	result := CheckEnvUrl([]string{`FOO`, `http://example.com/blah`})
	if result != "" {
		t.Error(`Expected "", got`, result)
	}
}

func TestCheckEnvUrl_correct(t *testing.T) {

	result := CheckEnvUrl([]string{`FOO_URL`, `http://example.com/blah`})
	if result != `http://example.com/blah` {
		t.Error(`Expected "http://example.com/blah", got`, result)
	}

}
