package preprocessor

import "testing"

func TestNearestColor(t *testing.T) {
	palette := []color{
		{0x00, 0x00, 0x00},
		{0x55, 0x55, 0x55},
		{0xAA, 0xAA, 0xAA},
		{0xFF, 0xFF, 0xFF},
	}

	if !nearestColor([]float64{0, 0, 0}, palette).equals(palette[0]) {
		t.Errorf("Expected nearest color to be 0x00,0x00,0x00")
	}
	if !nearestColor([]float64{10, 10, 10}, palette).equals(palette[0]) {
		t.Errorf("Expected nearest color to be 0x00,0x00,0x00")
	}
	if !nearestColor([]float64{0x55, 0x55, 0x55}, palette).equals(palette[1]) {
		t.Errorf("Expected nearest color to be 0x55,0x55,0x55")
	}
	if !nearestColor([]float64{0x55, 0x55, 0x56}, palette).equals(palette[1]) {
		t.Errorf("Expected nearest color to be 0x55,0x55,0x55")
	}
	if !nearestColor([]float64{0x55, 0x56, 0x55}, palette).equals(palette[1]) {
		t.Errorf("Expected nearest color to be 0x55,0x55,0x55")
	}
	if !nearestColor([]float64{0x56, 0x55, 0x55}, palette).equals(palette[1]) {
		t.Errorf("Expected nearest color to be 0x55,0x55,0x55")
	}
	if !nearestColor([]float64{0xFF, 0xFF, 0xFF}, palette).equals(palette[3]) {
		t.Errorf("Expected nearest color to be 0xFF,0xFF,0xFF")
	}
	if !nearestColor([]float64{0xFE, 0xFF, 0xFF}, palette).equals(palette[3]) {
		t.Errorf("Expected nearest color to be 0xFF,0xFF,0xFF")
	}
	if !nearestColor([]float64{0xFF, 0xFE, 0xFF}, palette).equals(palette[3]) {
		t.Errorf("Expected nearest color to be 0xFF,0xFF,0xFF")
	}
	if !nearestColor([]float64{0xFF, 0xFF, 0xFE}, palette).equals(palette[3]) {
		t.Errorf("Expected nearest color to be 0xFF,0xFF,0xFF")
	}
}
