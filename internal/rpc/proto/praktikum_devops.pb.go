// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: internal/rpc/proto/praktikum_devops.proto

package proto

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

// Запрос на обновление метрик
type AddMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []byte `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"` // json-массив метрик в виде слайса из байт
}

func (x *AddMetricsRequest) Reset() {
	*x = AddMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetricsRequest) ProtoMessage() {}

func (x *AddMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetricsRequest.ProtoReflect.Descriptor instead.
func (*AddMetricsRequest) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{0}
}

func (x *AddMetricsRequest) GetMetrics() []byte {
	if x != nil {
		return x.Metrics
	}
	return nil
}

// Результат обновления метрик
type AddMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"` // сообщение об ошибке
}

func (x *AddMetricsResponse) Reset() {
	*x = AddMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetricsResponse) ProtoMessage() {}

func (x *AddMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetricsResponse.ProtoReflect.Descriptor instead.
func (*AddMetricsResponse) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{1}
}

func (x *AddMetricsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Результат проверки доступности БД
type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{2}
}

// Результат проверки доступности БД
type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int64  `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"` // http-статус
	Error  string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`    // сообщение об ошибке
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{3}
}

func (x *PingResponse) GetStatus() int64 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *PingResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Запрос на обновление метрик
type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mtype string `protobuf:"bytes,1,opt,name=mtype,proto3" json:"mtype,omitempty"` // тип метрики
	Mname string `protobuf:"bytes,2,opt,name=mname,proto3" json:"mname,omitempty"` // имя метрики
	Val   string `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`     // имя метрики
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateMetricRequest) GetMtype() string {
	if x != nil {
		return x.Mtype
	}
	return ""
}

func (x *UpdateMetricRequest) GetMname() string {
	if x != nil {
		return x.Mname
	}
	return ""
}

func (x *UpdateMetricRequest) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

// Результат обновления метрик
type UpdateMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"` // сообщение об ошибке
}

func (x *UpdateMetricResponse) Reset() {
	*x = UpdateMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricResponse) ProtoMessage() {}

func (x *UpdateMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricResponse) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Запрос на обновление метрик
type UpdateJSONMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []byte `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"` // json-представление метрики для обновления
}

func (x *UpdateJSONMetricRequest) Reset() {
	*x = UpdateJSONMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateJSONMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateJSONMetricRequest) ProtoMessage() {}

func (x *UpdateJSONMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateJSONMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateJSONMetricRequest) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateJSONMetricRequest) GetMetrics() []byte {
	if x != nil {
		return x.Metrics
	}
	return nil
}

// Результат обновления метрик
type UpdateJSONMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []byte `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"` /// json-представление обновленной метрики
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`     // сообщение об ошибке
}

func (x *UpdateJSONMetricResponse) Reset() {
	*x = UpdateJSONMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateJSONMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateJSONMetricResponse) ProtoMessage() {}

func (x *UpdateJSONMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateJSONMetricResponse.ProtoReflect.Descriptor instead.
func (*UpdateJSONMetricResponse) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateJSONMetricResponse) GetMetrics() []byte {
	if x != nil {
		return x.Metrics
	}
	return nil
}

func (x *UpdateJSONMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mtype string `protobuf:"bytes,1,opt,name=mtype,proto3" json:"mtype,omitempty"` // тип метрики
	Mname string `protobuf:"bytes,2,opt,name=mname,proto3" json:"mname,omitempty"` // имя метрики
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricRequest.ProtoReflect.Descriptor instead.
func (*GetMetricRequest) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{8}
}

func (x *GetMetricRequest) GetMtype() string {
	if x != nil {
		return x.Mtype
	}
	return ""
}

func (x *GetMetricRequest) GetMname() string {
	if x != nil {
		return x.Mname
	}
	return ""
}

type GetMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` // значение метрики
	Error  string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`   // ошибка получения метрики
}

func (x *GetMetricResponse) Reset() {
	*x = GetMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricResponse.ProtoReflect.Descriptor instead.
func (*GetMetricResponse) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{9}
}

func (x *GetMetricResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *GetMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetJSONMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric []byte `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"` // json-массив метрик в виде слайса из байт
}

func (x *GetJSONMetricRequest) Reset() {
	*x = GetJSONMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJSONMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJSONMetricRequest) ProtoMessage() {}

func (x *GetJSONMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJSONMetricRequest.ProtoReflect.Descriptor instead.
func (*GetJSONMetricRequest) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{10}
}

func (x *GetJSONMetricRequest) GetMetric() []byte {
	if x != nil {
		return x.Metric
	}
	return nil
}

type GetJSONMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result []byte `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` // значение метрики
	Error  string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`   // ошибка получения метрики
}

func (x *GetJSONMetricResponse) Reset() {
	*x = GetJSONMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJSONMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJSONMetricResponse) ProtoMessage() {}

func (x *GetJSONMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpc_proto_praktikum_devops_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJSONMetricResponse.ProtoReflect.Descriptor instead.
func (*GetJSONMetricResponse) Descriptor() ([]byte, []int) {
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP(), []int{11}
}

func (x *GetJSONMetricResponse) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *GetJSONMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_internal_rpc_proto_praktikum_devops_proto protoreflect.FileDescriptor

var file_internal_rpc_proto_praktikum_devops_proto_rawDesc = []byte{
	0x0a, 0x29, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x61, 0x6b, 0x74, 0x69, 0x6b, 0x75, 0x6d, 0x5f, 0x64,
	0x65, 0x76, 0x6f, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x2d, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x22, 0x2a, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x0d, 0x0a,
	0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3c, 0x0a, 0x0c,
	0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x53, 0x0a, 0x13, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x22,
	0x2c, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x33, 0x0a,
	0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x22, 0x4a, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4a, 0x53, 0x4f, 0x4e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3e,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x41,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x2e, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x22, 0x45, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xa7, 0x03, 0x0a, 0x07, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x2f, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x53, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f,
	0x4e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x4a, 0x53, 0x4f, 0x4e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x72, 0x65, 0x62, 0x75, 0x73, 0x32, 0x30, 0x31, 0x35, 0x2f, 0x70, 0x72, 0x61, 0x6b, 0x74,
	0x69, 0x6b, 0x75, 0x6d, 0x2d, 0x64, 0x65, 0x76, 0x6f, 0x70, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_rpc_proto_praktikum_devops_proto_rawDescOnce sync.Once
	file_internal_rpc_proto_praktikum_devops_proto_rawDescData = file_internal_rpc_proto_praktikum_devops_proto_rawDesc
)

func file_internal_rpc_proto_praktikum_devops_proto_rawDescGZIP() []byte {
	file_internal_rpc_proto_praktikum_devops_proto_rawDescOnce.Do(func() {
		file_internal_rpc_proto_praktikum_devops_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_rpc_proto_praktikum_devops_proto_rawDescData)
	})
	return file_internal_rpc_proto_praktikum_devops_proto_rawDescData
}

var file_internal_rpc_proto_praktikum_devops_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_internal_rpc_proto_praktikum_devops_proto_goTypes = []interface{}{
	(*AddMetricsRequest)(nil),        // 0: proto.AddMetricsRequest
	(*AddMetricsResponse)(nil),       // 1: proto.AddMetricsResponse
	(*PingRequest)(nil),              // 2: proto.PingRequest
	(*PingResponse)(nil),             // 3: proto.PingResponse
	(*UpdateMetricRequest)(nil),      // 4: proto.UpdateMetricRequest
	(*UpdateMetricResponse)(nil),     // 5: proto.UpdateMetricResponse
	(*UpdateJSONMetricRequest)(nil),  // 6: proto.UpdateJSONMetricRequest
	(*UpdateJSONMetricResponse)(nil), // 7: proto.UpdateJSONMetricResponse
	(*GetMetricRequest)(nil),         // 8: proto.GetMetricRequest
	(*GetMetricResponse)(nil),        // 9: proto.GetMetricResponse
	(*GetJSONMetricRequest)(nil),     // 10: proto.GetJSONMetricRequest
	(*GetJSONMetricResponse)(nil),    // 11: proto.GetJSONMetricResponse
}
var file_internal_rpc_proto_praktikum_devops_proto_depIdxs = []int32{
	2,  // 0: proto.Metrics.Ping:input_type -> proto.PingRequest
	0,  // 1: proto.Metrics.AddMetrics:input_type -> proto.AddMetricsRequest
	4,  // 2: proto.Metrics.UpdateMetric:input_type -> proto.UpdateMetricRequest
	6,  // 3: proto.Metrics.UpdateJSONMetric:input_type -> proto.UpdateJSONMetricRequest
	8,  // 4: proto.Metrics.GetMetric:input_type -> proto.GetMetricRequest
	10, // 5: proto.Metrics.GetJSONMetric:input_type -> proto.GetJSONMetricRequest
	3,  // 6: proto.Metrics.Ping:output_type -> proto.PingResponse
	1,  // 7: proto.Metrics.AddMetrics:output_type -> proto.AddMetricsResponse
	5,  // 8: proto.Metrics.UpdateMetric:output_type -> proto.UpdateMetricResponse
	7,  // 9: proto.Metrics.UpdateJSONMetric:output_type -> proto.UpdateJSONMetricResponse
	9,  // 10: proto.Metrics.GetMetric:output_type -> proto.GetMetricResponse
	11, // 11: proto.Metrics.GetJSONMetric:output_type -> proto.GetJSONMetricResponse
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_internal_rpc_proto_praktikum_devops_proto_init() }
func file_internal_rpc_proto_praktikum_devops_proto_init() {
	if File_internal_rpc_proto_praktikum_devops_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddMetricsRequest); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddMetricsResponse); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricRequest); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricResponse); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateJSONMetricRequest); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateJSONMetricResponse); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetricRequest); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetricResponse); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetJSONMetricRequest); i {
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
		file_internal_rpc_proto_praktikum_devops_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetJSONMetricResponse); i {
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
			RawDescriptor: file_internal_rpc_proto_praktikum_devops_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_rpc_proto_praktikum_devops_proto_goTypes,
		DependencyIndexes: file_internal_rpc_proto_praktikum_devops_proto_depIdxs,
		MessageInfos:      file_internal_rpc_proto_praktikum_devops_proto_msgTypes,
	}.Build()
	File_internal_rpc_proto_praktikum_devops_proto = out.File
	file_internal_rpc_proto_praktikum_devops_proto_rawDesc = nil
	file_internal_rpc_proto_praktikum_devops_proto_goTypes = nil
	file_internal_rpc_proto_praktikum_devops_proto_depIdxs = nil
}