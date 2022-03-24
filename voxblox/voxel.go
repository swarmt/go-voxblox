package voxblox

import "sync"

type TsdfVoxel struct {
	Index IndexType
	sync.RWMutex
	distance float64
	weight   float64
	color    Color
}

func (v *TsdfVoxel) getWeight() float64 {
	v.RLock()
	defer v.RUnlock()
	return v.weight
}

func (v *TsdfVoxel) setWeight(weight float64) {
	v.Lock()
	defer v.Unlock()
	v.weight = weight
}

func (v *TsdfVoxel) getDistance() float64 {
	v.RLock()
	defer v.RUnlock()
	return v.distance
}

func (v *TsdfVoxel) setDistance(distance float64) {
	v.Lock()
	defer v.Unlock()
	v.distance = distance
}

func (v *TsdfVoxel) getColor() Color {
	v.RLock()
	defer v.RUnlock()
	return v.color
}

func (v *TsdfVoxel) setColor(color Color) {
	v.Lock()
	defer v.Unlock()
	v.color = color
}

func NewVoxel(index IndexType) *TsdfVoxel {
	return &TsdfVoxel{
		Index: index,
		color: Color{127, 127, 127},
	}
}
