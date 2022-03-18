go-voxblox
---

A Go implementation of [Voxblox](https://github.com/ethz-asl/voxblox)

Voxblox system diagram

![System Diagram](.readme/system-diagram.png)

go-voxblox system diagram

```mermaid
graph LR
  tsdfIntegrator[TSDF Integrator]
  tsdfMap[TSDF Map]
  meshIntegrator[Mesh Integrator]
  meshLayer[Mesh Layer]

  Sensor -- Pointcloud2 --> tsdfIntegrator
  poseEstimate -- 6DOF Pose --> tsdfIntegrator
  tsdfIntegrator --> tsdfMap --> tsdfIntegrator
  tsdfMap --> meshIntegrator
  meshIntegrator --> meshLayer --> meshIntegrator
  meshLayer --> gRPC 
```

![Unit Test Cylinder](.readme/cylinder-mesh.png)

## TODO

* Mesh generation
* Remove distant blocks
* Integrators
    * Merged
    * Fast
* ROS integration
* gRPC mesh server
* Logging
* System tests
* ICP?
* Linear indexing on voxels?
* CUDA?

## References

* [CHISEL](http://www.roboticsproceedings.org/rss11/p40.pdf)

