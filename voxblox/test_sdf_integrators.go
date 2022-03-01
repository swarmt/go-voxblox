package voxblox

var world *SimulationWorld

func init() {
	// Create a test environment.
	// It consists of a 10x10x7 m environment with a cylinder in the middle.
	voxelSize := 0.1
	minBound := Point{-5.0, -5.0, -1.0}
	maxBound := Point{5.0, 5.0, 6.0}
	world = NewSimulationWorld(voxelSize, minBound, maxBound)
	cylinder := Cylinder{
		Center: Point{0.0, 0.0, 2.0},
		Radius: 2.0,
		Height: 4.0,
	}
	world.AddObject(cylinder)
	world.AddGroundLevel(0.0)

}
