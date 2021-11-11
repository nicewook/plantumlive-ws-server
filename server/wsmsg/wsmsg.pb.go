// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: server/wsmsg/wsmsg.proto

package wsmsg

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

type Type int32

const (
	// command type
	Type_Unknown   Type = 0
	Type_Connected Type = 1
	Type_Join      Type = 2
	Type_Leave     Type = 3
	// message type
	Type_Msg Type = 100
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0:   "Unknown",
		1:   "Connected",
		2:   "Join",
		3:   "Leave",
		100: "Msg",
	}
	Type_value = map[string]int32{
		"Unknown":   0,
		"Connected": 1,
		"Join":      2,
		"Leave":     3,
		"Msg":       100,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_server_wsmsg_wsmsg_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_server_wsmsg_wsmsg_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_server_wsmsg_wsmsg_proto_rawDescGZIP(), []int{0}
}

type WebsocketMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      Type   `protobuf:"varint,1,opt,name=type,proto3,enum=wsmsg.Type" json:"type,omitempty"`
	SessionId string `protobuf:"bytes,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Username  string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Message   string `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *WebsocketMessage) Reset() {
	*x = WebsocketMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_wsmsg_wsmsg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebsocketMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebsocketMessage) ProtoMessage() {}

func (x *WebsocketMessage) ProtoReflect() protoreflect.Message {
	mi := &file_server_wsmsg_wsmsg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebsocketMessage.ProtoReflect.Descriptor instead.
func (*WebsocketMessage) Descriptor() ([]byte, []int) {
	return file_server_wsmsg_wsmsg_proto_rawDescGZIP(), []int{0}
}

func (x *WebsocketMessage) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_Unknown
}

func (x *WebsocketMessage) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *WebsocketMessage) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *WebsocketMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_server_wsmsg_wsmsg_proto protoreflect.FileDescriptor

var file_server_wsmsg_wsmsg_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x77, 0x73, 0x6d, 0x73, 0x67, 0x2f, 0x77,
	0x73, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x77, 0x73, 0x6d, 0x73,
	0x67, 0x22, 0x88, 0x01, 0x0a, 0x10, 0x57, 0x65, 0x62, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x77, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x40, 0x0a, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10,
	0x00, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x10, 0x01,
	0x12, 0x08, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x65,
	0x61, 0x76, 0x65, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x10, 0x64, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x2f, 0x77, 0x73, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_server_wsmsg_wsmsg_proto_rawDescOnce sync.Once
	file_server_wsmsg_wsmsg_proto_rawDescData = file_server_wsmsg_wsmsg_proto_rawDesc
)

func file_server_wsmsg_wsmsg_proto_rawDescGZIP() []byte {
	file_server_wsmsg_wsmsg_proto_rawDescOnce.Do(func() {
		file_server_wsmsg_wsmsg_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_wsmsg_wsmsg_proto_rawDescData)
	})
	return file_server_wsmsg_wsmsg_proto_rawDescData
}

var file_server_wsmsg_wsmsg_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_server_wsmsg_wsmsg_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_server_wsmsg_wsmsg_proto_goTypes = []interface{}{
	(Type)(0),                // 0: wsmsg.Type
	(*WebsocketMessage)(nil), // 1: wsmsg.WebsocketMessage
}
var file_server_wsmsg_wsmsg_proto_depIdxs = []int32{
	0, // 0: wsmsg.WebsocketMessage.type:type_name -> wsmsg.Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_server_wsmsg_wsmsg_proto_init() }
func file_server_wsmsg_wsmsg_proto_init() {
	if File_server_wsmsg_wsmsg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_wsmsg_wsmsg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebsocketMessage); i {
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
			RawDescriptor: file_server_wsmsg_wsmsg_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_server_wsmsg_wsmsg_proto_goTypes,
		DependencyIndexes: file_server_wsmsg_wsmsg_proto_depIdxs,
		EnumInfos:         file_server_wsmsg_wsmsg_proto_enumTypes,
		MessageInfos:      file_server_wsmsg_wsmsg_proto_msgTypes,
	}.Build()
	File_server_wsmsg_wsmsg_proto = out.File
	file_server_wsmsg_wsmsg_proto_rawDesc = nil
	file_server_wsmsg_wsmsg_proto_goTypes = nil
	file_server_wsmsg_wsmsg_proto_depIdxs = nil
}
