package voxblox

type TsdfMap struct {
	TsdfVoxelSize     float64
	TsdfVoxelsPerSide int
	TsdfLayer         *TsdfLayer
}

func NewTsdfMap(tsdfVoxelSize float64, tsdfVoxelsPerSide int) *TsdfMap {
	m := new(TsdfMap)
	m.TsdfVoxelSize = tsdfVoxelSize
	m.TsdfVoxelsPerSide = tsdfVoxelsPerSide
	m.TsdfLayer = NewTsdfLayer(m.TsdfVoxelSize, m.TsdfVoxelsPerSide)
	return m
}

func (t *TsdfMap) GetTsdfLayerPtr() *TsdfLayer {
	return t.TsdfLayer
}
