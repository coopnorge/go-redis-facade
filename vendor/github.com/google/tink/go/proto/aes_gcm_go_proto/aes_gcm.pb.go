// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.11.4
// source: third_party/tink/proto/aes_gcm.proto

package aes_gcm_go_proto

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

// only allowing IV size in bytes = 12 and tag size in bytes = 16
// Thus, accept no params.
type AesGcmKeyFormat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeySize uint32 `protobuf:"varint,2,opt,name=key_size,json=keySize,proto3" json:"key_size,omitempty"`
	Version uint32 `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *AesGcmKeyFormat) Reset() {
	*x = AesGcmKeyFormat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_third_party_tink_proto_aes_gcm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AesGcmKeyFormat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AesGcmKeyFormat) ProtoMessage() {}

func (x *AesGcmKeyFormat) ProtoReflect() protoreflect.Message {
	mi := &file_third_party_tink_proto_aes_gcm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AesGcmKeyFormat.ProtoReflect.Descriptor instead.
func (*AesGcmKeyFormat) Descriptor() ([]byte, []int) {
	return file_third_party_tink_proto_aes_gcm_proto_rawDescGZIP(), []int{0}
}

func (x *AesGcmKeyFormat) GetKeySize() uint32 {
	if x != nil {
		return x.KeySize
	}
	return 0
}

func (x *AesGcmKeyFormat) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

// key_type: type.googleapis.com/google.crypto.tink.AesGcmKey
type AesGcmKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version  uint32 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	KeyValue []byte `protobuf:"bytes,3,opt,name=key_value,json=keyValue,proto3" json:"key_value,omitempty"`
}

func (x *AesGcmKey) Reset() {
	*x = AesGcmKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_third_party_tink_proto_aes_gcm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AesGcmKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AesGcmKey) ProtoMessage() {}

func (x *AesGcmKey) ProtoReflect() protoreflect.Message {
	mi := &file_third_party_tink_proto_aes_gcm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AesGcmKey.ProtoReflect.Descriptor instead.
func (*AesGcmKey) Descriptor() ([]byte, []int) {
	return file_third_party_tink_proto_aes_gcm_proto_rawDescGZIP(), []int{1}
}

func (x *AesGcmKey) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *AesGcmKey) GetKeyValue() []byte {
	if x != nil {
		return x.KeyValue
	}
	return nil
}

var File_third_party_tink_proto_aes_gcm_proto protoreflect.FileDescriptor

var file_third_party_tink_proto_aes_gcm_proto_rawDesc = []byte{
	0x0a, 0x24, 0x74, 0x68, 0x69, 0x72, 0x64, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x79, 0x2f, 0x74, 0x69,
	0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x65, 0x73, 0x5f, 0x67, 0x63, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x22, 0x46, 0x0a, 0x0f, 0x41, 0x65,
	0x73, 0x47, 0x63, 0x6d, 0x4b, 0x65, 0x79, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x19, 0x0a,
	0x08, 0x6b, 0x65, 0x79, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x6b, 0x65, 0x79, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0x42, 0x0a, 0x09, 0x41, 0x65, 0x73, 0x47, 0x63, 0x6d, 0x4b, 0x65, 0x79, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6b, 0x65, 0x79,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6b, 0x65,
	0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x4f, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x74, 0x69, 0x6e, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x74, 0x69, 0x6e, 0x6b,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x65, 0x73, 0x5f, 0x67, 0x63, 0x6d, 0x5f, 0x67,
	0x6f, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_third_party_tink_proto_aes_gcm_proto_rawDescOnce sync.Once
	file_third_party_tink_proto_aes_gcm_proto_rawDescData = file_third_party_tink_proto_aes_gcm_proto_rawDesc
)

func file_third_party_tink_proto_aes_gcm_proto_rawDescGZIP() []byte {
	file_third_party_tink_proto_aes_gcm_proto_rawDescOnce.Do(func() {
		file_third_party_tink_proto_aes_gcm_proto_rawDescData = protoimpl.X.CompressGZIP(file_third_party_tink_proto_aes_gcm_proto_rawDescData)
	})
	return file_third_party_tink_proto_aes_gcm_proto_rawDescData
}

var file_third_party_tink_proto_aes_gcm_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_third_party_tink_proto_aes_gcm_proto_goTypes = []interface{}{
	(*AesGcmKeyFormat)(nil), // 0: google.crypto.tink.AesGcmKeyFormat
	(*AesGcmKey)(nil),       // 1: google.crypto.tink.AesGcmKey
}
var file_third_party_tink_proto_aes_gcm_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_third_party_tink_proto_aes_gcm_proto_init() }
func file_third_party_tink_proto_aes_gcm_proto_init() {
	if File_third_party_tink_proto_aes_gcm_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_third_party_tink_proto_aes_gcm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AesGcmKeyFormat); i {
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
		file_third_party_tink_proto_aes_gcm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AesGcmKey); i {
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
			RawDescriptor: file_third_party_tink_proto_aes_gcm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_third_party_tink_proto_aes_gcm_proto_goTypes,
		DependencyIndexes: file_third_party_tink_proto_aes_gcm_proto_depIdxs,
		MessageInfos:      file_third_party_tink_proto_aes_gcm_proto_msgTypes,
	}.Build()
	File_third_party_tink_proto_aes_gcm_proto = out.File
	file_third_party_tink_proto_aes_gcm_proto_rawDesc = nil
	file_third_party_tink_proto_aes_gcm_proto_goTypes = nil
	file_third_party_tink_proto_aes_gcm_proto_depIdxs = nil
}
