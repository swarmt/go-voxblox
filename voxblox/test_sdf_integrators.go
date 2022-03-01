package voxblox

var world *SimulationWorld

func init() {
	voxelSize := 0.1

	minBound := Point{-5.0, -5.0, -1.0}
	maxBound := Point{5.0, 5.0, 6.0}

	world = NewSimulationWorld(voxelSize, minBound, maxBound)

	//cylinderCenter := Point{0.0, 0.0, 2.0}
	//cylinderRadius := 2.0
	//cylinderHeight := 4.0
	world.AddObject()
}
