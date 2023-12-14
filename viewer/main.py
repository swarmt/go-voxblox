import time
import grpc
import rerun as rr
from mesh_service_pb2 import GetMeshRequest
from mesh_service_pb2_grpc import MeshServiceStub

if __name__ == "__main__":
    client = grpc.insecure_channel("localhost:50051")
    stub = MeshServiceStub(client)
    request = GetMeshRequest()

    rr.init("go-voxblox", spawn=True)
    rr.log("world", rr.ViewCoordinates.RIGHT_HAND_Z_UP)

    frame = 0
    try:
        while True:
            result = stub.GetMeshBlocks(request)
            rr.set_time_sequence("frame_idx", frame)

            for mesh_block in result:
                rr.log(
                    f"world/{mesh_block.index}", rr.Asset3D(contents=mesh_block.bytes)
                )

            frame += 1
            time.sleep(0.1)

    except KeyboardInterrupt:
        print("Interrupted by user, shutting down.")
