package voxblox

import "sync"

type TsdfVoxel struct {
	Index    IndexType
	distance float64
	weight   float64
	color    Color
	mutex    sync.RWMutex
}

func (v *TsdfVoxel) getWeight() float64 {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.weight
}

func (v *TsdfVoxel) setWeight(weight float64) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.weight = weight
}

func (v *TsdfVoxel) getDistance() float64 {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.distance
}

func (v *TsdfVoxel) setDistance(distance float64) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.distance = distance
}

func (v *TsdfVoxel) getColor() Color {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.color
}

func (v *TsdfVoxel) setColor(color Color) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.color = color
}

func NewVoxel(index IndexType) *TsdfVoxel {
	return &TsdfVoxel{
		Index: index,
		mutex: sync.RWMutex{},
	}
}
