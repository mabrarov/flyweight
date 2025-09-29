package main

import (
	"fmt"

	"github.com/mabrarov/flyweight/pkg/flyweight"
)

func main() {
	f := flyweight.NewFactory[int, float64]()
	for range 10 {
		one, _ := f.Get(1, func() (*float64, error) {
			return toPtr(1.0), nil
		})
		two, _ := f.Get(2, func() (*float64, error) {
			return toPtr(2.0), nil
		})
		three, _ := f.Get(3, func() (*float64, error) {
			return toPtr(3.0), nil
		})
		fmt.Printf("%p: %f, %p: %f, %p, %f\n", one, *one, two, *two, three, *three)
	}
}

func toPtr[T any](v T) *T {
	return &v
}
