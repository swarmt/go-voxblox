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

[Cow and Lady Dataset](https://projects.asl.ethz.ch/datasets/doku.php?id=iros2017/)

Note: This needs to be decompressed to run real time with ```rosbag decompress```
![Cow and Lady Dataset](.readme/cow-and-lady.png)

## TODO

* Merge duplicate vertices
* Better unit tests
* Remove distant blocks
* ROS integration
* gRPC mesh server
* Logging
* System tests
* ICP?
* Linear indexing on voxels?
* CUDA?

## References

* [CHISEL](http://www.roboticsproceedings.org/rss11/p40.pdf)

