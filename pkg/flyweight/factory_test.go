package flyweight

import (
	"errors"
	"math"
	"runtime"
	"testing"
)

func TestFactory_Get_SameValueAsBuilt(t *testing.T) {
	var builtValue *[]byte
	buildCounter := uint(0)
	builder := func() (*[]byte, error) {
		builtValue = new([]byte)
		*builtValue = make([]byte, 1000)
		if buildCounter < math.MaxUint {
			buildCounter++
		}
		return builtValue, nil
	}
	f := NewFactory[int, []byte]()

	// Ensure same value as the built one is returned all the time,
	// because memory, which `builtValue` variable points to, can not be collected by GC
	for range 1000 {
		actualValue, err := f.Get(1, builder)
		if err != nil {
			t.Fatalf("Error getting value: %v", err)
		}
		if actualValue == nil {
			t.Fatal("Expected not nil value if builder returned it")
		}
		if actualValue != builtValue {
			t.Fatalf("Wrong value: expected %v, got %v", builtValue, actualValue)
		}
		// Give a chance for GC
		actualValue = nil
		runtime.GC()
	}

	// Ensure all tests caused building of new value just once
	if buildCounter != 1 {
		t.Fatalf("Should build new value just once, but got %d time(s)", buildCounter)
	}

	// Just to avoid compiler optimization
	if builtValue == nil {
		t.Fatal()
	}
}

func TestFactory_Get_SameErrorAsBuilder(t *testing.T) {
	buildErr := errors.New("test")
	f := NewFactory[int, []byte]()

	actualValue, actualErr := f.Get(1, func() (*[]byte, error) {
		return new([]byte), buildErr
	})

	if actualValue != nil {
		t.Fatal("nil value should be returned when builder returns error")
	}
	if !errors.Is(actualErr, buildErr) {
		t.Fatalf("Expected error %v, got %v", buildErr, actualErr)
	}
}

// This test depends on GC - if it collects unused memory during loop execution -
// so this test can fail if GC decides to not collect unused memory (e.g. due to performance considerations).
func TestFactory_UnusedEntriesCleaned(t *testing.T) {
	buildCounter := uint(0)
	builder := func() (*[]byte, error) {
		b := new([]byte)
		*b = make([]byte, 1000)
		if buildCounter < math.MaxUint {
			buildCounter++
		}
		return b, nil
	}
	f := NewFactory[int, []byte]()

	for range 1000 {
		actualValue, err := f.Get(1, builder)
		if err != nil {
			t.Fatalf("Error getting value: %v", err)
		}
		if actualValue == nil {
			t.Fatal("Expected not nil value if builder returned it")
		}
		// Hope that during one of iterations GC collects unused memory
		actualValue = nil
		runtime.GC()
	}

	// Ensure stored value was removed from `f`, because respective memory was collected by GC.
	// As a side effect `builder` should be called more than once.
	if buildCounter <= 1 {
		t.Fatal("expected GC to collect unused memory causing removal of respective entry from factory")
	}
}
