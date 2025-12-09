package lib

import "fmt"

func DeferErr(f func() error) {
	err := f()
	if err != nil {
		fmt.Println("error on deferred function:", err.Error())
	}
}

func Must[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}

	return res
}
