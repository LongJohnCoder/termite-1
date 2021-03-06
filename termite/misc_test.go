package termite

import (
	"log"
	"testing"
)

var _ = log.Println

func TestHasDirPrefix(t *testing.T) {
	if !HasDirPrefix("a/b", "a") {
		t.Errorf("HasDirPrefix(a/b, a) fail")
	}
	if HasDirPrefix("a/b", "ab") {
		t.Errorf("HasDirPrefix(a/b, ab) succeed")
	}
}

func TestEscapeRegexp(t *testing.T) {
	s := EscapeRegexp("a+b")
	if s != "a\\+b" {
		t.Error("mismatch", s)
	}
}

func TestDetectFiles(t *testing.T) {
	fs := DetectFiles("/src/foo", "gcc /src/foo/bar.cc -I/src/foo/baz")
	result := map[string]int{}
	for _, f := range fs {
		result[f] = 1
	}
	if len(result) != 2 {
		t.Error("length", result)
	}
	if result["/src/foo/bar.cc"] != 1 || result["/src/foo/baz"] != 1 {
		t.Error("not found", result)
	}
}

func TestParseCommand(t *testing.T) {
	fail := []string{
		"echo hoi;",
		"echo \"${hoi}\"",
		"a && b",
		"a || b",
		"echo a*b",
		"echo 'x' \\ >> temp.sed",
	}
	for _, s := range fail {
		result := ParseCommand(s)
		if result != nil {
			t.Errorf("should fail: cmd=%#v, result=%#v", s, result)
		}
	}

	type Succ struct {
		cmd string
		res []string
	}

	succ := []Succ{
		{"echo \"a'b\"", []string{"echo", "a'b"}},
		{"\"a'b\"", []string{"a'b"}},
		{"a\\ b", []string{"a b"}},
		{"a'x y'b", []string{"ax yb"}},
		{"echo \"a[]<>*&;;\"", []string{"echo", "a[]<>*&;;"}},
		{"a   b", []string{"a", "b"}},
		{"a\\$b", []string{"a$b"}},
	}
	for _, entry := range succ {
		r := ParseCommand(entry.cmd)
		if len(r) != len(entry.res) {
			t.Error("len mismatch", r, entry)
		} else {
			for i := range r {
				if r[i] != entry.res[i] {
					t.Errorf("component mismatch for %v comp %d got %v want %v",
						entry.cmd, i, r[i], entry.res[i])
				}
			}
		}
	}
}

func TestMakeUnescape(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"abc\ndef", "abc\ndef"},
		{"abc\\\ndef", "abcdef"},
		{"abc\\\\\ndef", "abc\\\\\ndef"},
	}
	for i, c := range cases {
		got := MakeUnescape(c.in)
		if c.out != got {
			t.Errorf("%d: MakeUnescape(%q) = %q != %q", i, c.in, got, c.out)
		}
	}
}

func TestIntToExponent(t *testing.T) {
	e := IntToExponent(1)
	if e != 0 {
		t.Error("1", e)
	}
	e = IntToExponent(2)
	if e != 1 {
		t.Error("2", e)
	}
	e = IntToExponent(3)
	if e != 2 {
		t.Error("3", e)
	}
	e = IntToExponent(4)
	if e != 2 {
		t.Error("4", e)
	}
}
