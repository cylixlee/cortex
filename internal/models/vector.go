package models

import (
	pgvector "github.com/pgvector/pgvector-go"
)

type Vector = pgvector.Vector

func Float64ToVector(floats []float64) Vector {
	f32 := make([]float32, len(floats))
	for i, f := range floats {
		f32[i] = float32(f)
	}
	return pgvector.NewVector(f32)
}
