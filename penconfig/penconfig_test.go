package penconfig

// colours to be dealt with by
// https://github.com/go-playground/colors or LocalColour parser

import (
	"io/ioutil"
	"math"
	"os"
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
	p, ok := lpc.GetPen(1, "fineliner", "narrow")
	if !ok {
		t.Errorf("could not get pen 1/fineliner/narrow")
	}
	t.Log("got pen ", p)

	n := p.GetWidth("narrow")
	if n != 0.8 {
		t.Errorf("width expected 0.8, got %f", n)
	}
	n = p.GetWidth("broad")
	if math.Round(n*1000) != math.Round(0.8*2.125/1.875*1000) {
		t.Errorf("width expected %f, got %f", 0.8*2.125/1.875, n)
	}

	c := p.GetColour()
	if c.R != 0 || c.G != 0 || c.B != 255 {
		t.Errorf("c %v != 0/0/255", c)
	}

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

	lpc, err := LoadYaml(y)
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
	p, ok := lpc.GetPen(1, "nonsense", "problem")
	if ok {
		t.Errorf("got pen 1/nonsense/problem")
	}
	t.Logf("empty pen %+v : %t", p, p == (&PenConfig{}))

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

// TestNewConfigFromFile tests loading a pen configuration from a yaml
// file
func TestNewConfigFromFile(t *testing.T) {

	y := []byte(`
"all":
  - pen:     fineliner
    weight:  narrow
    width:   0.8
    color:   "#963387"
    opacity: 0.8`)

	tmpFile, err := ioutil.TempFile("", "settings_")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(y)
	if err != nil {
		t.Fatal(err)
	}

	lpc, err := NewPenConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := lpc["all"]; !ok {
		t.Error("lpc structure not ok ", lpc)
	}
}

// TestNewConfigFromFileFail tests failing to loading a pen
// configuration from a yaml file
func TestNewConfigFromFileFail(t *testing.T) {

	tmpFile, err := ioutil.TempFile("", "settings_")
	if err != nil {
		t.Error(err)
	}
	tName := tmpFile.Name()
	os.Remove(tmpFile.Name())

	_, err = NewPenConfigFromFile(tName)
	if err == nil {
		t.Fatalf("loading pen config should fail")
	}
}
