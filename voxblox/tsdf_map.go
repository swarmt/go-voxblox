package voxblox

type TsdfMap struct {
	TsdfVoxelSize     float64
	TsdfVoxelsPerSide int32
	TsdfLayer         *Layer
}

func NewTsdfMap(tsdfVoxelSize float64, tsdfVoxelsPerSide int32) *TsdfMap {
	m := new(TsdfMap)
	m.TsdfVoxelSize = tsdfVoxelSize
	m.TsdfVoxelsPerSide = tsdfVoxelsPerSide
	m.TsdfLayer = NewLayer(m.TsdfVoxelSize, m.TsdfVoxelsPerSide)
	return m
}

func (t *TsdfMap) GetTsdfLayerPtr() *Layer {
	return t.TsdfLayer
}
