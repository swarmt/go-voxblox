package voxblox

type TsdfMap struct {
	voxelSize     float64
	voxelsPerSide int
	layer         *TsdfLayer
}

func NewTsdfMap(voxelSize float64, voxelsPerSide int) *TsdfMap {
	m := new(TsdfMap)
	m.voxelSize = voxelSize
	m.voxelsPerSide = voxelsPerSide
	m.layer = NewTsdfLayer(voxelSize, voxelsPerSide)
	return m
}

func (t *TsdfMap) GetTsdfLayerPtr() *TsdfLayer {
	return t.layer
}
