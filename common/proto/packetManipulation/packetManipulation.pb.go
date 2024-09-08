// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: common/proto/packetManipulation/packetManipulation.proto

package packetManipulation

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PctmConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Server         string          `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Client         string          `protobuf:"bytes,2,opt,name=client,proto3" json:"client,omitempty"`
	DurationConfig *DurationConfig `protobuf:"bytes,3,opt,name=duration_config,json=durationConfig,proto3" json:"duration_config,omitempty"`
}

func (x *PctmConfig) Reset() {
	*x = PctmConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PctmConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PctmConfig) ProtoMessage() {}

func (x *PctmConfig) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PctmConfig.ProtoReflect.Descriptor instead.
func (*PctmConfig) Descriptor() ([]byte, []int) {
	return file_common_proto_packetManipulation_packetManipulation_proto_rawDescGZIP(), []int{0}
}

func (x *PctmConfig) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *PctmConfig) GetClient() string {
	if x != nil {
		return x.Client
	}
	return ""
}

func (x *PctmConfig) GetDurationConfig() *DurationConfig {
	if x != nil {
		return x.DurationConfig
	}
	return nil
}

type DurationConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DurationSeconds int32   `protobuf:"varint,1,opt,name=duration_seconds,json=durationSeconds,proto3" json:"duration_seconds,omitempty"`
	DropRate        float32 `protobuf:"fixed32,2,opt,name=drop_rate,json=dropRate,proto3" json:"drop_rate,omitempty"`          // Added drop rate for packet dropping
	CorruptRate     float32 `protobuf:"fixed32,6,opt,name=corrupt_rate,json=corruptRate,proto3" json:"corrupt_rate,omitempty"` // Added rate for packet corruption
}

func (x *DurationConfig) Reset() {
	*x = DurationConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DurationConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DurationConfig) ProtoMessage() {}

func (x *DurationConfig) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DurationConfig.ProtoReflect.Descriptor instead.
func (*DurationConfig) Descriptor() ([]byte, []int) {
	return file_common_proto_packetManipulation_packetManipulation_proto_rawDescGZIP(), []int{1}
}

func (x *DurationConfig) GetDurationSeconds() int32 {
	if x != nil {
		return x.DurationSeconds
	}
	return 0
}

func (x *DurationConfig) GetDropRate() float32 {
	if x != nil {
		return x.DropRate
	}
	return 0
}

func (x *DurationConfig) GetCorruptRate() float32 {
	if x != nil {
		return x.CorruptRate
	}
	return 0
}

type RequestForManipulationProxy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *PctmConfig `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *RequestForManipulationProxy) Reset() {
	*x = RequestForManipulationProxy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestForManipulationProxy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestForManipulationProxy) ProtoMessage() {}

func (x *RequestForManipulationProxy) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestForManipulationProxy.ProtoReflect.Descriptor instead.
func (*RequestForManipulationProxy) Descriptor() ([]byte, []int) {
	return file_common_proto_packetManipulation_packetManipulation_proto_rawDescGZIP(), []int{2}
}

func (x *RequestForManipulationProxy) GetConfig() *PctmConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

type ResponseFromManipulationProxy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ResponseFromManipulationProxy) Reset() {
	*x = ResponseFromManipulationProxy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseFromManipulationProxy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseFromManipulationProxy) ProtoMessage() {}

func (x *ResponseFromManipulationProxy) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseFromManipulationProxy.ProtoReflect.Descriptor instead.
func (*ResponseFromManipulationProxy) Descriptor() ([]byte, []int) {
	return file_common_proto_packetManipulation_packetManipulation_proto_rawDescGZIP(), []int{3}
}

func (x *ResponseFromManipulationProxy) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_common_proto_packetManipulation_packetManipulation_proto protoreflect.FileDescriptor

var file_common_proto_packetManipulation_packetManipulation_proto_rawDesc = []byte{
	0x0a, 0x38, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x01, 0x0a, 0x0a,
	0x70, 0x63, 0x74, 0x6d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x06, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x59, 0x0a, 0x0f,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x61,
	0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0e, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xa6, 0x01, 0x0a, 0x0e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x32, 0x0a, 0x10, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x52, 0x0f, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x2c,
	0x0a, 0x09, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x02, 0x42, 0x0f, 0xfa, 0x42, 0x0c, 0x0a, 0x0a, 0x1d, 0x00, 0x00, 0x80, 0x3f, 0x2d, 0x00, 0x00,
	0x00, 0x00, 0x52, 0x08, 0x64, 0x72, 0x6f, 0x70, 0x52, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x0c,
	0x63, 0x6f, 0x72, 0x72, 0x75, 0x70, 0x74, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x02, 0x42, 0x0f, 0xfa, 0x42, 0x0c, 0x0a, 0x0a, 0x1d, 0x00, 0x00, 0x80, 0x3f, 0x2d, 0x00,
	0x00, 0x00, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x72, 0x72, 0x75, 0x70, 0x74, 0x52, 0x61, 0x74, 0x65,
	0x22, 0x63, 0x0a, 0x1b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x4d, 0x61,
	0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x12,
	0x44, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x2e, 0x70, 0x63, 0x74, 0x6d, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x39, 0x0a, 0x1d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x46, 0x72, 0x6f, 0x6d, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0xa4, 0x01, 0x0a, 0x19, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70,
	0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x86,
	0x01, 0x0a, 0x16, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x12, 0x33, 0x2e, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x4d, 0x61, 0x6e,
	0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x1a, 0x35,
	0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x46, 0x72, 0x6f, 0x6d, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x50, 0x72, 0x6f, 0x78, 0x79, 0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x2f, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x4d, 0x61, 0x6e, 0x69, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_packetManipulation_packetManipulation_proto_rawDescOnce sync.Once
	file_common_proto_packetManipulation_packetManipulation_proto_rawDescData = file_common_proto_packetManipulation_packetManipulation_proto_rawDesc
)

func file_common_proto_packetManipulation_packetManipulation_proto_rawDescGZIP() []byte {
	file_common_proto_packetManipulation_packetManipulation_proto_rawDescOnce.Do(func() {
		file_common_proto_packetManipulation_packetManipulation_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_packetManipulation_packetManipulation_proto_rawDescData)
	})
	return file_common_proto_packetManipulation_packetManipulation_proto_rawDescData
}

var file_common_proto_packetManipulation_packetManipulation_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_common_proto_packetManipulation_packetManipulation_proto_goTypes = []any{
	(*PctmConfig)(nil),                    // 0: packetManipulationServ.pctmConfig
	(*DurationConfig)(nil),                // 1: packetManipulationServ.DurationConfig
	(*RequestForManipulationProxy)(nil),   // 2: packetManipulationServ.RequestForManipulationProxy
	(*ResponseFromManipulationProxy)(nil), // 3: packetManipulationServ.ResponseFromManipulationProxy
}
var file_common_proto_packetManipulation_packetManipulation_proto_depIdxs = []int32{
	1, // 0: packetManipulationServ.pctmConfig.duration_config:type_name -> packetManipulationServ.DurationConfig
	0, // 1: packetManipulationServ.RequestForManipulationProxy.config:type_name -> packetManipulationServ.pctmConfig
	2, // 2: packetManipulationServ.PacketManipulationService.StartManipulationProxy:input_type -> packetManipulationServ.RequestForManipulationProxy
	3, // 3: packetManipulationServ.PacketManipulationService.StartManipulationProxy:output_type -> packetManipulationServ.ResponseFromManipulationProxy
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_proto_packetManipulation_packetManipulation_proto_init() }
func file_common_proto_packetManipulation_packetManipulation_proto_init() {
	if File_common_proto_packetManipulation_packetManipulation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*PctmConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DurationConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RequestForManipulationProxy); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_packetManipulation_packetManipulation_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ResponseFromManipulationProxy); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_proto_packetManipulation_packetManipulation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_common_proto_packetManipulation_packetManipulation_proto_goTypes,
		DependencyIndexes: file_common_proto_packetManipulation_packetManipulation_proto_depIdxs,
		MessageInfos:      file_common_proto_packetManipulation_packetManipulation_proto_msgTypes,
	}.Build()
	File_common_proto_packetManipulation_packetManipulation_proto = out.File
	file_common_proto_packetManipulation_packetManipulation_proto_rawDesc = nil
	file_common_proto_packetManipulation_packetManipulation_proto_goTypes = nil
	file_common_proto_packetManipulation_packetManipulation_proto_depIdxs = nil
}
