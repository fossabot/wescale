//
//Copyright 2019 The Vitess Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// This file contains the service definition for making management API
// calls to mysqlctld.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.3
// source: mysqlctl.proto

package mysqlctl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	topodata "github.com/wesql/wescale/go/vt/proto/topodata"
	vttime "github.com/wesql/wescale/go/vt/proto/vttime"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Status is an enum representing the possible status of a backup.
type BackupInfo_Status int32

const (
	BackupInfo_UNKNOWN    BackupInfo_Status = 0
	BackupInfo_INCOMPLETE BackupInfo_Status = 1
	BackupInfo_COMPLETE   BackupInfo_Status = 2
	// A backup status of INVALID should be set if the backup is complete
	// but unusable in some way (partial upload, corrupt file, etc).
	BackupInfo_INVALID BackupInfo_Status = 3
	// A backup status of VALID should be set if the backup is both
	// complete and usuable.
	BackupInfo_VALID BackupInfo_Status = 4
)

// Enum value maps for BackupInfo_Status.
var (
	BackupInfo_Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "INCOMPLETE",
		2: "COMPLETE",
		3: "INVALID",
		4: "VALID",
	}
	BackupInfo_Status_value = map[string]int32{
		"UNKNOWN":    0,
		"INCOMPLETE": 1,
		"COMPLETE":   2,
		"INVALID":    3,
		"VALID":      4,
	}
)

func (x BackupInfo_Status) Enum() *BackupInfo_Status {
	p := new(BackupInfo_Status)
	*p = x
	return p
}

func (x BackupInfo_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BackupInfo_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlctl_proto_enumTypes[0].Descriptor()
}

func (BackupInfo_Status) Type() protoreflect.EnumType {
	return &file_mysqlctl_proto_enumTypes[0]
}

func (x BackupInfo_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BackupInfo_Status.Descriptor instead.
func (BackupInfo_Status) EnumDescriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{10, 0}
}

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MysqldArgs []string `protobuf:"bytes,1,rep,name=mysqld_args,json=mysqldArgs,proto3" json:"mysqld_args,omitempty"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{0}
}

func (x *StartRequest) GetMysqldArgs() []string {
	if x != nil {
		return x.MysqldArgs
	}
	return nil
}

type StartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartResponse) Reset() {
	*x = StartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartResponse) ProtoMessage() {}

func (x *StartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartResponse.ProtoReflect.Descriptor instead.
func (*StartResponse) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{1}
}

type ShutdownRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WaitForMysqld bool `protobuf:"varint,1,opt,name=wait_for_mysqld,json=waitForMysqld,proto3" json:"wait_for_mysqld,omitempty"`
}

func (x *ShutdownRequest) Reset() {
	*x = ShutdownRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShutdownRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownRequest) ProtoMessage() {}

func (x *ShutdownRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownRequest.ProtoReflect.Descriptor instead.
func (*ShutdownRequest) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{2}
}

func (x *ShutdownRequest) GetWaitForMysqld() bool {
	if x != nil {
		return x.WaitForMysqld
	}
	return false
}

type ShutdownResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ShutdownResponse) Reset() {
	*x = ShutdownResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShutdownResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownResponse) ProtoMessage() {}

func (x *ShutdownResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownResponse.ProtoReflect.Descriptor instead.
func (*ShutdownResponse) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{3}
}

type RunMysqlUpgradeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RunMysqlUpgradeRequest) Reset() {
	*x = RunMysqlUpgradeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunMysqlUpgradeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunMysqlUpgradeRequest) ProtoMessage() {}

func (x *RunMysqlUpgradeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunMysqlUpgradeRequest.ProtoReflect.Descriptor instead.
func (*RunMysqlUpgradeRequest) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{4}
}

type RunMysqlUpgradeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RunMysqlUpgradeResponse) Reset() {
	*x = RunMysqlUpgradeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunMysqlUpgradeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunMysqlUpgradeResponse) ProtoMessage() {}

func (x *RunMysqlUpgradeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunMysqlUpgradeResponse.ProtoReflect.Descriptor instead.
func (*RunMysqlUpgradeResponse) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{5}
}

type ReinitConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReinitConfigRequest) Reset() {
	*x = ReinitConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReinitConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReinitConfigRequest) ProtoMessage() {}

func (x *ReinitConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReinitConfigRequest.ProtoReflect.Descriptor instead.
func (*ReinitConfigRequest) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{6}
}

type ReinitConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReinitConfigResponse) Reset() {
	*x = ReinitConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReinitConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReinitConfigResponse) ProtoMessage() {}

func (x *ReinitConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReinitConfigResponse.ProtoReflect.Descriptor instead.
func (*ReinitConfigResponse) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{7}
}

type RefreshConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RefreshConfigRequest) Reset() {
	*x = RefreshConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshConfigRequest) ProtoMessage() {}

func (x *RefreshConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshConfigRequest.ProtoReflect.Descriptor instead.
func (*RefreshConfigRequest) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{8}
}

type RefreshConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RefreshConfigResponse) Reset() {
	*x = RefreshConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshConfigResponse) ProtoMessage() {}

func (x *RefreshConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshConfigResponse.ProtoReflect.Descriptor instead.
func (*RefreshConfigResponse) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{9}
}

// BackupInfo is the read-only attributes of a mysqlctl/backupstorage.BackupHandle.
type BackupInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Directory   string                `protobuf:"bytes,2,opt,name=directory,proto3" json:"directory,omitempty"`
	Keyspace    string                `protobuf:"bytes,3,opt,name=keyspace,proto3" json:"keyspace,omitempty"`
	Shard       string                `protobuf:"bytes,4,opt,name=shard,proto3" json:"shard,omitempty"`
	TabletAlias *topodata.TabletAlias `protobuf:"bytes,5,opt,name=tablet_alias,json=tabletAlias,proto3" json:"tablet_alias,omitempty"`
	Time        *vttime.Time          `protobuf:"bytes,6,opt,name=time,proto3" json:"time,omitempty"`
	// Engine is the name of the backupengine implementation used to create
	// this backup.
	Engine string            `protobuf:"bytes,7,opt,name=engine,proto3" json:"engine,omitempty"`
	Status BackupInfo_Status `protobuf:"varint,8,opt,name=status,proto3,enum=mysqlctl.BackupInfo_Status" json:"status,omitempty"`
}

func (x *BackupInfo) Reset() {
	*x = BackupInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlctl_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackupInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackupInfo) ProtoMessage() {}

func (x *BackupInfo) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlctl_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackupInfo.ProtoReflect.Descriptor instead.
func (*BackupInfo) Descriptor() ([]byte, []int) {
	return file_mysqlctl_proto_rawDescGZIP(), []int{10}
}

func (x *BackupInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BackupInfo) GetDirectory() string {
	if x != nil {
		return x.Directory
	}
	return ""
}

func (x *BackupInfo) GetKeyspace() string {
	if x != nil {
		return x.Keyspace
	}
	return ""
}

func (x *BackupInfo) GetShard() string {
	if x != nil {
		return x.Shard
	}
	return ""
}

func (x *BackupInfo) GetTabletAlias() *topodata.TabletAlias {
	if x != nil {
		return x.TabletAlias
	}
	return nil
}

func (x *BackupInfo) GetTime() *vttime.Time {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *BackupInfo) GetEngine() string {
	if x != nil {
		return x.Engine
	}
	return ""
}

func (x *BackupInfo) GetStatus() BackupInfo_Status {
	if x != nil {
		return x.Status
	}
	return BackupInfo_UNKNOWN
}

var File_mysqlctl_proto protoreflect.FileDescriptor

var file_mysqlctl_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x1a, 0x0e, 0x74, 0x6f, 0x70, 0x6f,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x76, 0x74, 0x74, 0x69,
	0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x79, 0x73, 0x71,
	0x6c, 0x64, 0x5f, 0x61, 0x72, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x6d,
	0x79, 0x73, 0x71, 0x6c, 0x64, 0x41, 0x72, 0x67, 0x73, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x0a, 0x0f, 0x53, 0x68,
	0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a,
	0x0f, 0x77, 0x61, 0x69, 0x74, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x77, 0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x4d,
	0x79, 0x73, 0x71, 0x6c, 0x64, 0x22, 0x12, 0x0a, 0x10, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x0a, 0x16, 0x52, 0x75, 0x6e,
	0x4d, 0x79, 0x73, 0x71, 0x6c, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x19, 0x0a, 0x17, 0x52, 0x75, 0x6e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x55,
	0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15,
	0x0a, 0x13, 0x52, 0x65, 0x69, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x16, 0x0a, 0x14, 0x52, 0x65, 0x69, 0x6e, 0x69, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x0a,
	0x14, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x17, 0x0a, 0x15, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xe6,
	0x02, 0x0a, 0x0a, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x12, 0x38, 0x0a, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x61, 0x6c, 0x69, 0x61,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x6f, 0x70, 0x6f, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x52, 0x0b,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x20, 0x0a, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x76, 0x74, 0x74, 0x69,
	0x6d, 0x65, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x12, 0x33, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c,
	0x2e, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x4b, 0x0a, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x49, 0x4e, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x10,
	0x01, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x02, 0x12,
	0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x04, 0x32, 0x8a, 0x03, 0x0a, 0x08, 0x4d, 0x79, 0x73, 0x71,
	0x6c, 0x43, 0x74, 0x6c, 0x12, 0x3a, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x16, 0x2e,
	0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c,
	0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x43, 0x0a, 0x08, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x12, 0x19, 0x2e, 0x6d,
	0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63,
	0x74, 0x6c, 0x2e, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x58, 0x0a, 0x0f, 0x52, 0x75, 0x6e, 0x4d, 0x79, 0x73, 0x71,
	0x6c, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x12, 0x20, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c,
	0x63, 0x74, 0x6c, 0x2e, 0x52, 0x75, 0x6e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x55, 0x70, 0x67, 0x72,
	0x61, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6d, 0x79, 0x73,
	0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x75, 0x6e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x55, 0x70,
	0x67, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4f, 0x0a, 0x0c, 0x52, 0x65, 0x69, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x1d, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x69, 0x6e, 0x69,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x69, 0x6e, 0x69, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x52, 0x0a, 0x0d, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x1e, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x27, 0x5a, 0x25, 0x76, 0x69, 0x74, 0x65, 0x73, 0x73, 0x2e, 0x69,
	0x6f, 0x2f, 0x76, 0x69, 0x74, 0x65, 0x73, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x76, 0x74, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x63, 0x74, 0x6c, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mysqlctl_proto_rawDescOnce sync.Once
	file_mysqlctl_proto_rawDescData = file_mysqlctl_proto_rawDesc
)

func file_mysqlctl_proto_rawDescGZIP() []byte {
	file_mysqlctl_proto_rawDescOnce.Do(func() {
		file_mysqlctl_proto_rawDescData = protoimpl.X.CompressGZIP(file_mysqlctl_proto_rawDescData)
	})
	return file_mysqlctl_proto_rawDescData
}

var file_mysqlctl_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mysqlctl_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_mysqlctl_proto_goTypes = []interface{}{
	(BackupInfo_Status)(0),          // 0: mysqlctl.BackupInfo.Status
	(*StartRequest)(nil),            // 1: mysqlctl.StartRequest
	(*StartResponse)(nil),           // 2: mysqlctl.StartResponse
	(*ShutdownRequest)(nil),         // 3: mysqlctl.ShutdownRequest
	(*ShutdownResponse)(nil),        // 4: mysqlctl.ShutdownResponse
	(*RunMysqlUpgradeRequest)(nil),  // 5: mysqlctl.RunMysqlUpgradeRequest
	(*RunMysqlUpgradeResponse)(nil), // 6: mysqlctl.RunMysqlUpgradeResponse
	(*ReinitConfigRequest)(nil),     // 7: mysqlctl.ReinitConfigRequest
	(*ReinitConfigResponse)(nil),    // 8: mysqlctl.ReinitConfigResponse
	(*RefreshConfigRequest)(nil),    // 9: mysqlctl.RefreshConfigRequest
	(*RefreshConfigResponse)(nil),   // 10: mysqlctl.RefreshConfigResponse
	(*BackupInfo)(nil),              // 11: mysqlctl.BackupInfo
	(*topodata.TabletAlias)(nil),    // 12: topodata.TabletAlias
	(*vttime.Time)(nil),             // 13: vttime.Time
}
var file_mysqlctl_proto_depIdxs = []int32{
	12, // 0: mysqlctl.BackupInfo.tablet_alias:type_name -> topodata.TabletAlias
	13, // 1: mysqlctl.BackupInfo.time:type_name -> vttime.Time
	0,  // 2: mysqlctl.BackupInfo.status:type_name -> mysqlctl.BackupInfo.Status
	1,  // 3: mysqlctl.MysqlCtl.Start:input_type -> mysqlctl.StartRequest
	3,  // 4: mysqlctl.MysqlCtl.Shutdown:input_type -> mysqlctl.ShutdownRequest
	5,  // 5: mysqlctl.MysqlCtl.RunMysqlUpgrade:input_type -> mysqlctl.RunMysqlUpgradeRequest
	7,  // 6: mysqlctl.MysqlCtl.ReinitConfig:input_type -> mysqlctl.ReinitConfigRequest
	9,  // 7: mysqlctl.MysqlCtl.RefreshConfig:input_type -> mysqlctl.RefreshConfigRequest
	2,  // 8: mysqlctl.MysqlCtl.Start:output_type -> mysqlctl.StartResponse
	4,  // 9: mysqlctl.MysqlCtl.Shutdown:output_type -> mysqlctl.ShutdownResponse
	6,  // 10: mysqlctl.MysqlCtl.RunMysqlUpgrade:output_type -> mysqlctl.RunMysqlUpgradeResponse
	8,  // 11: mysqlctl.MysqlCtl.ReinitConfig:output_type -> mysqlctl.ReinitConfigResponse
	10, // 12: mysqlctl.MysqlCtl.RefreshConfig:output_type -> mysqlctl.RefreshConfigResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_mysqlctl_proto_init() }
func file_mysqlctl_proto_init() {
	if File_mysqlctl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mysqlctl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartRequest); i {
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
		file_mysqlctl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartResponse); i {
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
		file_mysqlctl_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShutdownRequest); i {
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
		file_mysqlctl_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShutdownResponse); i {
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
		file_mysqlctl_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunMysqlUpgradeRequest); i {
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
		file_mysqlctl_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunMysqlUpgradeResponse); i {
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
		file_mysqlctl_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReinitConfigRequest); i {
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
		file_mysqlctl_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReinitConfigResponse); i {
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
		file_mysqlctl_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshConfigRequest); i {
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
		file_mysqlctl_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshConfigResponse); i {
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
		file_mysqlctl_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackupInfo); i {
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
			RawDescriptor: file_mysqlctl_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mysqlctl_proto_goTypes,
		DependencyIndexes: file_mysqlctl_proto_depIdxs,
		EnumInfos:         file_mysqlctl_proto_enumTypes,
		MessageInfos:      file_mysqlctl_proto_msgTypes,
	}.Build()
	File_mysqlctl_proto = out.File
	file_mysqlctl_proto_rawDesc = nil
	file_mysqlctl_proto_goTypes = nil
	file_mysqlctl_proto_depIdxs = nil
}
