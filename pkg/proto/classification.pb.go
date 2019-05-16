// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/classification.proto

package elections_mediawatch_io

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Request message
type UserFeatures struct {
	Followers            int64    `protobuf:"varint,1,opt,name=followers,proto3" json:"followers,omitempty"`
	Friends              int64    `protobuf:"varint,2,opt,name=friends,proto3" json:"friends,omitempty"`
	Statuses             int64    `protobuf:"varint,3,opt,name=statuses,proto3" json:"statuses,omitempty"`
	Favorites            int64    `protobuf:"varint,4,opt,name=favorites,proto3" json:"favorites,omitempty"`
	Lists                int64    `protobuf:"varint,5,opt,name=lists,proto3" json:"lists,omitempty"`
	Ffr                  float64  `protobuf:"fixed64,6,opt,name=ffr,proto3" json:"ffr,omitempty"`
	Stfv                 float64  `protobuf:"fixed64,7,opt,name=stfv,proto3" json:"stfv,omitempty"`
	Fstfv                float64  `protobuf:"fixed64,8,opt,name=fstfv,proto3" json:"fstfv,omitempty"`
	Dates                float64  `protobuf:"fixed64,9,opt,name=dates,proto3" json:"dates,omitempty"`
	Actions              float64  `protobuf:"fixed64,10,opt,name=actions,proto3" json:"actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserFeatures) Reset()         { *m = UserFeatures{} }
func (m *UserFeatures) String() string { return proto.CompactTextString(m) }
func (*UserFeatures) ProtoMessage()    {}
func (*UserFeatures) Descriptor() ([]byte, []int) {
	return fileDescriptor_classification_02a8ac5be12a7f15, []int{0}
}
func (m *UserFeatures) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserFeatures.Unmarshal(m, b)
}
func (m *UserFeatures) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserFeatures.Marshal(b, m, deterministic)
}
func (dst *UserFeatures) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserFeatures.Merge(dst, src)
}
func (m *UserFeatures) XXX_Size() int {
	return xxx_messageInfo_UserFeatures.Size(m)
}
func (m *UserFeatures) XXX_DiscardUnknown() {
	xxx_messageInfo_UserFeatures.DiscardUnknown(m)
}

var xxx_messageInfo_UserFeatures proto.InternalMessageInfo

func (m *UserFeatures) GetFollowers() int64 {
	if m != nil {
		return m.Followers
	}
	return 0
}

func (m *UserFeatures) GetFriends() int64 {
	if m != nil {
		return m.Friends
	}
	return 0
}

func (m *UserFeatures) GetStatuses() int64 {
	if m != nil {
		return m.Statuses
	}
	return 0
}

func (m *UserFeatures) GetFavorites() int64 {
	if m != nil {
		return m.Favorites
	}
	return 0
}

func (m *UserFeatures) GetLists() int64 {
	if m != nil {
		return m.Lists
	}
	return 0
}

func (m *UserFeatures) GetFfr() float64 {
	if m != nil {
		return m.Ffr
	}
	return 0
}

func (m *UserFeatures) GetStfv() float64 {
	if m != nil {
		return m.Stfv
	}
	return 0
}

func (m *UserFeatures) GetFstfv() float64 {
	if m != nil {
		return m.Fstfv
	}
	return 0
}

func (m *UserFeatures) GetDates() float64 {
	if m != nil {
		return m.Dates
	}
	return 0
}

func (m *UserFeatures) GetActions() float64 {
	if m != nil {
		return m.Actions
	}
	return 0
}

// Response message
type UserClass struct {
	Label                string   `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	Score                float64  `protobuf:"fixed64,2,opt,name=score,proto3" json:"score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserClass) Reset()         { *m = UserClass{} }
func (m *UserClass) String() string { return proto.CompactTextString(m) }
func (*UserClass) ProtoMessage()    {}
func (*UserClass) Descriptor() ([]byte, []int) {
	return fileDescriptor_classification_02a8ac5be12a7f15, []int{1}
}
func (m *UserClass) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserClass.Unmarshal(m, b)
}
func (m *UserClass) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserClass.Marshal(b, m, deterministic)
}
func (dst *UserClass) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserClass.Merge(dst, src)
}
func (m *UserClass) XXX_Size() int {
	return xxx_messageInfo_UserClass.Size(m)
}
func (m *UserClass) XXX_DiscardUnknown() {
	xxx_messageInfo_UserClass.DiscardUnknown(m)
}

var xxx_messageInfo_UserClass proto.InternalMessageInfo

func (m *UserClass) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *UserClass) GetScore() float64 {
	if m != nil {
		return m.Score
	}
	return 0
}

type Model struct {
	Active               float32  `protobuf:"fixed32,1,opt,name=active,proto3" json:"active,omitempty"`
	Bot                  float32  `protobuf:"fixed32,2,opt,name=bot,proto3" json:"bot,omitempty"`
	Influencer           float32  `protobuf:"fixed32,3,opt,name=influencer,proto3" json:"influencer,omitempty"`
	New                  float32  `protobuf:"fixed32,4,opt,name=new,proto3" json:"new,omitempty"`
	Normal               float32  `protobuf:"fixed32,5,opt,name=normal,proto3" json:"normal,omitempty"`
	Other                float32  `protobuf:"fixed32,6,opt,name=other,proto3" json:"other,omitempty"`
	Retweeter            float32  `protobuf:"fixed32,7,opt,name=retweeter,proto3" json:"retweeter,omitempty"`
	SuperUser            float32  `protobuf:"fixed32,8,opt,name=super_user,json=superUser,proto3" json:"super_user,omitempty"`
	Unknown              float32  `protobuf:"fixed32,9,opt,name=unknown,proto3" json:"unknown,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Model) Reset()         { *m = Model{} }
func (m *Model) String() string { return proto.CompactTextString(m) }
func (*Model) ProtoMessage()    {}
func (*Model) Descriptor() ([]byte, []int) {
	return fileDescriptor_classification_02a8ac5be12a7f15, []int{2}
}
func (m *Model) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Model.Unmarshal(m, b)
}
func (m *Model) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Model.Marshal(b, m, deterministic)
}
func (dst *Model) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Model.Merge(dst, src)
}
func (m *Model) XXX_Size() int {
	return xxx_messageInfo_Model.Size(m)
}
func (m *Model) XXX_DiscardUnknown() {
	xxx_messageInfo_Model.DiscardUnknown(m)
}

var xxx_messageInfo_Model proto.InternalMessageInfo

func (m *Model) GetActive() float32 {
	if m != nil {
		return m.Active
	}
	return 0
}

func (m *Model) GetBot() float32 {
	if m != nil {
		return m.Bot
	}
	return 0
}

func (m *Model) GetInfluencer() float32 {
	if m != nil {
		return m.Influencer
	}
	return 0
}

func (m *Model) GetNew() float32 {
	if m != nil {
		return m.New
	}
	return 0
}

func (m *Model) GetNormal() float32 {
	if m != nil {
		return m.Normal
	}
	return 0
}

func (m *Model) GetOther() float32 {
	if m != nil {
		return m.Other
	}
	return 0
}

func (m *Model) GetRetweeter() float32 {
	if m != nil {
		return m.Retweeter
	}
	return 0
}

func (m *Model) GetSuperUser() float32 {
	if m != nil {
		return m.SuperUser
	}
	return 0
}

func (m *Model) GetUnknown() float32 {
	if m != nil {
		return m.Unknown
	}
	return 0
}

func init() {
	proto.RegisterType((*UserFeatures)(nil), "elections.mediawatch.io.UserFeatures")
	proto.RegisterType((*UserClass)(nil), "elections.mediawatch.io.UserClass")
	proto.RegisterType((*Model)(nil), "elections.mediawatch.io.Model")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClassificationClient is the client API for Classification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClassificationClient interface {
	Classify(ctx context.Context, in *UserFeatures, opts ...grpc.CallOption) (*UserClass, error)
}

type classificationClient struct {
	cc *grpc.ClientConn
}

func NewClassificationClient(cc *grpc.ClientConn) ClassificationClient {
	return &classificationClient{cc}
}

func (c *classificationClient) Classify(ctx context.Context, in *UserFeatures, opts ...grpc.CallOption) (*UserClass, error) {
	out := new(UserClass)
	err := c.cc.Invoke(ctx, "/elections.mediawatch.io.Classification/Classify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClassificationServer is the server API for Classification service.
type ClassificationServer interface {
	Classify(context.Context, *UserFeatures) (*UserClass, error)
}

func RegisterClassificationServer(s *grpc.Server, srv ClassificationServer) {
	s.RegisterService(&_Classification_serviceDesc, srv)
}

func _Classification_Classify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFeatures)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClassificationServer).Classify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elections.mediawatch.io.Classification/Classify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClassificationServer).Classify(ctx, req.(*UserFeatures))
	}
	return interceptor(ctx, in, info, handler)
}

var _Classification_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elections.mediawatch.io.Classification",
	HandlerType: (*ClassificationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Classify",
			Handler:    _Classification_Classify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/classification.proto",
}

func init() {
	proto.RegisterFile("proto/classification.proto", fileDescriptor_classification_02a8ac5be12a7f15)
}

var fileDescriptor_classification_02a8ac5be12a7f15 = []byte{
	// 399 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x8a, 0x14, 0x31,
	0x10, 0x86, 0xed, 0xcc, 0xce, 0xec, 0x74, 0x21, 0x22, 0x41, 0x34, 0x0c, 0x2a, 0x4b, 0x83, 0xb0,
	0xa7, 0x16, 0xf4, 0xe0, 0x03, 0x2c, 0x78, 0xf3, 0x12, 0x10, 0x8f, 0x92, 0xe9, 0xa9, 0x66, 0x83,
	0xd9, 0x64, 0x49, 0xa5, 0xa7, 0xf1, 0xea, 0xd3, 0xfa, 0x18, 0x52, 0x95, 0x69, 0x77, 0x3c, 0xb8,
	0xb7, 0xfa, 0xbf, 0x74, 0x55, 0x27, 0x1f, 0x05, 0xbb, 0xfb, 0x9c, 0x4a, 0x7a, 0x3f, 0x04, 0x47,
	0xe4, 0x47, 0x3f, 0xb8, 0xe2, 0x53, 0xec, 0x05, 0xea, 0x57, 0x18, 0x70, 0xe0, 0x4c, 0xfd, 0x1d,
	0x1e, 0xbc, 0x9b, 0x5d, 0x19, 0x6e, 0x7b, 0x9f, 0xba, 0x5f, 0x0a, 0x9e, 0x7e, 0x25, 0xcc, 0x9f,
	0xd1, 0x95, 0x29, 0x23, 0xe9, 0xd7, 0xd0, 0x8e, 0x29, 0x84, 0x34, 0x63, 0x26, 0xd3, 0x5c, 0x35,
	0xd7, 0x2b, 0xfb, 0x00, 0xb4, 0x81, 0xcb, 0x31, 0x7b, 0x8c, 0x07, 0x32, 0x4a, 0xce, 0x96, 0xa8,
	0x77, 0xb0, 0xa5, 0xe2, 0xca, 0x44, 0x48, 0x66, 0x25, 0x47, 0x7f, 0xb3, 0xcc, 0x74, 0xc7, 0x94,
	0x7d, 0x41, 0x32, 0x17, 0xa7, 0x99, 0x0b, 0xd0, 0x2f, 0x60, 0x1d, 0x3c, 0x15, 0x32, 0x6b, 0x39,
	0xa9, 0x41, 0x3f, 0x87, 0xd5, 0x38, 0x66, 0xb3, 0xb9, 0x6a, 0xae, 0x1b, 0xcb, 0xa5, 0xd6, 0x70,
	0x41, 0x65, 0x3c, 0x9a, 0x4b, 0x41, 0x52, 0x73, 0xef, 0x28, 0x70, 0x2b, 0xb0, 0x06, 0xa6, 0x07,
	0xc7, 0xff, 0x6a, 0x2b, 0x95, 0xc0, 0x77, 0x77, 0xd5, 0x81, 0x01, 0xe1, 0x4b, 0xec, 0x3e, 0x41,
	0xcb, 0x0e, 0x6e, 0xd8, 0x9c, 0x5c, 0xc7, 0xed, 0x31, 0xc8, 0xe3, 0x5b, 0x5b, 0x03, 0x53, 0x1a,
	0x52, 0x46, 0x79, 0x76, 0x63, 0x6b, 0xe8, 0x7e, 0x37, 0xb0, 0xfe, 0x92, 0x0e, 0x18, 0xf4, 0x4b,
	0xd8, 0xf0, 0xb4, 0x23, 0x4a, 0x9b, 0xb2, 0xa7, 0xc4, 0xcf, 0xd8, 0xa7, 0x22, 0x5d, 0xca, 0x72,
	0xa9, 0xdf, 0x02, 0xf8, 0x38, 0x86, 0x09, 0xe3, 0x80, 0x59, 0x54, 0x29, 0x7b, 0x46, 0xb8, 0x23,
	0xe2, 0x2c, 0x9a, 0x94, 0xe5, 0x92, 0x67, 0xc7, 0x94, 0xef, 0x5c, 0x10, 0x43, 0xca, 0x9e, 0x12,
	0xdf, 0x29, 0x95, 0x5b, 0xac, 0x92, 0x94, 0xad, 0x81, 0x65, 0x67, 0x2c, 0x33, 0x62, 0xc1, 0x2c,
	0xae, 0x94, 0x7d, 0x00, 0xfa, 0x0d, 0x00, 0x4d, 0xf7, 0x98, 0xbf, 0x4f, 0x84, 0x59, 0xac, 0x29,
	0xdb, 0x0a, 0x61, 0x03, 0xec, 0x68, 0x8a, 0x3f, 0x62, 0x9a, 0xa3, 0xb8, 0x53, 0x76, 0x89, 0x1f,
	0x3c, 0x3c, 0xbb, 0xf9, 0x67, 0xb3, 0xf4, 0x37, 0xd8, 0x9e, 0xc8, 0x4f, 0xfd, 0xae, 0xff, 0xcf,
	0x82, 0xf5, 0xe7, 0xcb, 0xb5, 0xeb, 0x1e, 0xfd, 0x4c, 0xa6, 0x75, 0x4f, 0xf6, 0x1b, 0xd9, 0xd9,
	0x8f, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x55, 0x53, 0xc6, 0xed, 0xd1, 0x02, 0x00, 0x00,
}
