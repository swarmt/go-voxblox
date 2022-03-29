go-voxblox
---
[![Go](https://github.com/swarmt/go-voxblox/actions/workflows/go.yml/badge.svg)](https://github.com/swarmt/go-voxblox/actions/workflows/go.yml)
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

## Performance

[Cow and Lady Dataset](https://projects.asl.ethz.ch/datasets/doku.php?id=iros2017/)

Note: This needs to be decompressed to run real time with ```rosbag decompress```

![Cow and Lady Dataset](.readme/cow-and-lady.png)

## TODO

* Merged integrator weights and speed
* Better / more unit tests
* Cache distant blocks with protobuf
* ROS integration
* gRPC mesh server
* Logging
* System tests
* Stress test / map size
* ICP?
* Linear indexing on voxels?
* CUDA?

## References

* [CHISEL](http://www.roboticsproceedings.org/rss11/p40.pdf)

