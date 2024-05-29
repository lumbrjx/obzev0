// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.21.12
// source: tcpService.proto

package tcpService

import (
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

type TcpConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReqDelay int32  `protobuf:"varint,1,opt,name=reqDelay,proto3" json:"reqDelay,omitempty"`
	ResDelay int32  `protobuf:"varint,2,opt,name=resDelay,proto3" json:"resDelay,omitempty"`
	Server   string `protobuf:"bytes,3,opt,name=server,proto3" json:"server,omitempty"`
	Client   string `protobuf:"bytes,4,opt,name=client,proto3" json:"client,omitempty"`
}

func (x *TcpConfig) Reset() {
	*x = TcpConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tcpService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TcpConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpConfig) ProtoMessage() {}

func (x *TcpConfig) ProtoReflect() protoreflect.Message {
	mi := &file_tcpService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpConfig.ProtoReflect.Descriptor instead.
func (*TcpConfig) Descriptor() ([]byte, []int) {
	return file_tcpService_proto_rawDescGZIP(), []int{0}
}

func (x *TcpConfig) GetReqDelay() int32 {
	if x != nil {
		return x.ReqDelay
	}
	return 0
}

func (x *TcpConfig) GetResDelay() int32 {
	if x != nil {
		return x.ResDelay
	}
	return 0
}

func (x *TcpConfig) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *TcpConfig) GetClient() string {
	if x != nil {
		return x.Client
	}
	return ""
}

type RequestForTcp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *TcpConfig `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *RequestForTcp) Reset() {
	*x = RequestForTcp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tcpService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestForTcp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestForTcp) ProtoMessage() {}

func (x *RequestForTcp) ProtoReflect() protoreflect.Message {
	mi := &file_tcpService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestForTcp.ProtoReflect.Descriptor instead.
func (*RequestForTcp) Descriptor() ([]byte, []int) {
	return file_tcpService_proto_rawDescGZIP(), []int{1}
}

func (x *RequestForTcp) GetConfig() *TcpConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

type ResponseFromTcp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric string `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *ResponseFromTcp) Reset() {
	*x = ResponseFromTcp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tcpService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseFromTcp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseFromTcp) ProtoMessage() {}

func (x *ResponseFromTcp) ProtoReflect() protoreflect.Message {
	mi := &file_tcpService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseFromTcp.ProtoReflect.Descriptor instead.
func (*ResponseFromTcp) Descriptor() ([]byte, []int) {
	return file_tcpService_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseFromTcp) GetMetric() string {
	if x != nil {
		return x.Metric
	}
	return ""
}

var File_tcpService_proto protoreflect.FileDescriptor

var file_tcpService_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x63, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x74, 0x63, 0x70, 0x53, 0x65, 0x72, 0x76, 0x22, 0x73, 0x0a, 0x09, 0x54,
	0x63, 0x70, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x44,
	0x65, 0x6c, 0x61, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x72, 0x65, 0x71, 0x44,
	0x65, 0x6c, 0x61, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x44, 0x65, 0x6c, 0x61, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x72, 0x65, 0x73, 0x44, 0x65, 0x6c, 0x61, 0x79,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x22, 0x3b, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x54, 0x63,
	0x70, 0x12, 0x2a, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x74, 0x63, 0x70, 0x53, 0x65, 0x72, 0x76, 0x2e, 0x54, 0x63, 0x70, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x29, 0x0a,
	0x0f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x54, 0x63, 0x70,
	0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x32, 0x52, 0x0a, 0x0a, 0x54, 0x63, 0x70, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x63, 0x70, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x74, 0x63, 0x70, 0x53, 0x65,
	0x72, 0x76, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x54, 0x63, 0x70,
	0x1a, 0x18, 0x2e, 0x74, 0x63, 0x70, 0x53, 0x65, 0x72, 0x76, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x54, 0x63, 0x70, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b,
	0x2f, 0x74, 0x63, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_tcpService_proto_rawDescOnce sync.Once
	file_tcpService_proto_rawDescData = file_tcpService_proto_rawDesc
)

func file_tcpService_proto_rawDescGZIP() []byte {
	file_tcpService_proto_rawDescOnce.Do(func() {
		file_tcpService_proto_rawDescData = protoimpl.X.CompressGZIP(file_tcpService_proto_rawDescData)
	})
	return file_tcpService_proto_rawDescData
}

var file_tcpService_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tcpService_proto_goTypes = []interface{}{
	(*TcpConfig)(nil),       // 0: tcpServ.TcpConfig
	(*RequestForTcp)(nil),   // 1: tcpServ.RequestForTcp
	(*ResponseFromTcp)(nil), // 2: tcpServ.ResponseFromTcp
}
var file_tcpService_proto_depIdxs = []int32{
	0, // 0: tcpServ.RequestForTcp.config:type_name -> tcpServ.TcpConfig
	1, // 1: tcpServ.TcpService.StartTcpServer:input_type -> tcpServ.RequestForTcp
	2, // 2: tcpServ.TcpService.StartTcpServer:output_type -> tcpServ.ResponseFromTcp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tcpService_proto_init() }
func file_tcpService_proto_init() {
	if File_tcpService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tcpService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TcpConfig); i {
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
		file_tcpService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestForTcp); i {
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
		file_tcpService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseFromTcp); i {
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
			RawDescriptor: file_tcpService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tcpService_proto_goTypes,
		DependencyIndexes: file_tcpService_proto_depIdxs,
		MessageInfos:      file_tcpService_proto_msgTypes,
	}.Build()
	File_tcpService_proto = out.File
	file_tcpService_proto_rawDesc = nil
	file_tcpService_proto_goTypes = nil
	file_tcpService_proto_depIdxs = nil
}
