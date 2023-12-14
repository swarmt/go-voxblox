from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class GetMeshRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetMeshResult(_message.Message):
    __slots__ = ("index", "bytes")
    INDEX_FIELD_NUMBER: _ClassVar[int]
    BYTES_FIELD_NUMBER: _ClassVar[int]
    index: str
    bytes: bytes
    def __init__(self, index: _Optional[str] = ..., bytes: _Optional[bytes] = ...) -> None: ...
