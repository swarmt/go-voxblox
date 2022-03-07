package voxblox

import "sync"

type TsdfVoxel struct {
	Index    IndexType
	distance float64
	weight   float64
	mutex    sync.RWMutex // TODO: I'm going to try using a mutex per voxel for its simplicity. May need to change this later.
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

func NewVoxel(index IndexType) *TsdfVoxel {
	return &TsdfVoxel{
		Index: index,
		mutex: sync.RWMutex{},
	}
}
