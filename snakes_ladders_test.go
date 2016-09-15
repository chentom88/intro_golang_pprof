package main

import (
	"testing"
)

func BenchmarkBad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunFromFile("input.txt")
	}
}

func BenchmarkGood(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunFromFileGood("input.txt")
	}
}
