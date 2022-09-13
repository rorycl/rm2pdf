package penconfig

// colours to be dealt with by
// https://github.com/go-playground/colors or LocalColour parser

import (
	"strings"
	"testing"
)

// TestPenConfigParse tests basic parsing
func TestPenConfigParse(t *testing.T) {

	y := []byte(`
all:
  - pen:     fineliner
    weight:  narrow
    width:   0.95
    color:   black
    opacity: 0.9

"1":
  - pen:     fineliner
    weight:  narrow
    width:   0.8
    color:   blue
    opacity: 0.8`)

	lpc, err := LoadYaml(y)
	if err != nil {
		t.Errorf("load config error %s : %s", y, err)
	}
	t.Log(lpc)

}

// TestPenConfigParseFail1 tests basic parsing failure
func TestPenConfigParseFail1(t *testing.T) {

	y := []byte(`
"fail":
  - pen:     fineliner
    weight:  narrow
    width:   0.8
    color:   blue
    opacity: 0.8`)

	_, err := LoadYaml(y)
	if err == nil {
		t.Error("load config should error with invalid layer name")
	}
	if !strings.Contains(err.Error(), "layer fail") {
		t.Error("error should contain 'layer fail'")
	}
	expected := "layer fail"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("error should contain '%s'", expected)
	}
}

// TestPenConfigParseFail2 tests basic parsing failure
func TestPenConfigParseFail2(t *testing.T) {

	y := []byte(`
"all":
  - pen:     nonsense
    weight:  narrow
    width:   0.8
    color:   blue
    opacity: 0.8`)

	_, err := LoadYaml(y)
	if err == nil {
		t.Error("load config should error with invalid pen name")
	}
	expected := "pen type nonsense not in"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("error should contain '%s'", expected)
	}
}

// TestPenConfigParseOK2 tests basic parsing
func TestPenConfigParseOK2(t *testing.T) {

	y := []byte(`
"all":
  - pen:     fineliner
    weight:  narrow
    width:   0.8
    color:   "#963387"
    opacity: 0.8`)

	lpc, err := LoadYaml(y)
	if err != nil {
		t.Errorf("load config unexpectedly errored with %s", err)
	}
	t.Log(lpc)
}
