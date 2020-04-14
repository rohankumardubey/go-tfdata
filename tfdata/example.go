// Package tfdata provides interfaces to interact with TFRecord files and TFExamples.
//
//
// Copyright (c) 2020, NVIDIA CORPORATION. All rights reserved.
//

package tfdata

import (
	"bytes"
	"image"
	"image/png"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/NVIDIA/go-tfdata/proto"
	"github.com/NVIDIA/go-tfdata/tfdata/internal/cmn"
)

type (
	TFExamplePipe chan *TFExample

	// TFExample is a wrapper over proto.Example struct generated by protoc from TensorFlow
	// tf.Example proto files. It is a golang representation of tf.Example datastructure.
	// It includes functions for adding elements to tf.Example.Features.
	TFExample struct {
		proto.Example
	}
)

// NewTFExample initializes empty TFExample and returns it.
func NewTFExample() *TFExample {
	ex := proto.Example{
		Features: &proto.Features{Feature: make(map[string]*proto.Feature)},
	}

	return &TFExample{ex}
}

func (e *TFExample) GetFeature(name string) *proto.Feature {
	return e.Features.Feature[name]
}

func (e *TFExample) GetInt64List(name string) []int64 {
	return e.GetFeature(name).GetInt64List().Value
}

func (e *TFExample) AddInt64List(name string, ints []int64) {
	e.Features.Feature[name] = &proto.Feature{Kind: &proto.Feature_Int64List{Int64List: &proto.Int64List{Value: ints}}}
}

func (e *TFExample) AddIntList(name string, ints []int) {
	ints64 := make([]int64, 0, len(ints))
	for _, i := range ints {
		ints64 = append(ints64, int64(i))
	}
	e.AddInt64List(name, ints64)
}

func (e *TFExample) AddInt64(name string, ints ...int64) {
	e.AddInt64List(name, ints)
}

func (e *TFExample) GetInt64(name string) int64 {
	return e.GetInt64List(name)[0]
}

func (e *TFExample) AddInt(name string, ints ...int) {
	e.AddIntList(name, ints)
}

func (e *TFExample) GetFloatList(name string) []float32 {
	return e.GetFeature(name).GetFloatList().Value
}

func (e *TFExample) AddFloatList(name string, floats []float32) {
	e.Features.Feature[name] = &proto.Feature{Kind: &proto.Feature_FloatList{FloatList: &proto.FloatList{Value: floats}}}
}

func (e *TFExample) GetFloat(name string) float32 {
	return e.GetFloatList(name)[0]
}

func (e *TFExample) AddFloat(name string, floats ...float32) {
	e.AddFloatList(name, floats)
}

func (e *TFExample) GetBytesList(name string) []byte {
	f := e.GetFeature(name).GetBytesList().Value
	cmn.Assert(len(f) == 1)
	return f[0]
}

func (e *TFExample) AddBytesList(name string, bytes [][]byte) {
	e.Features.Feature[name] = &proto.Feature{Kind: &proto.Feature_BytesList{BytesList: &proto.BytesList{Value: bytes}}}
}

func (e *TFExample) AddBytes(name string, bytes ...[]byte) {
	e.AddBytesList(name, bytes)
}

func (e *TFExample) AddImage(name string, img image.Image) error {
	buff := bytes.NewBuffer(make([]byte, 0, img.Bounds().Dx()*img.Bounds().Dy()*8))
	if err := png.Encode(buff, img); err != nil {
		return err
	}

	e.AddBytes(name, buff.Bytes())
	return nil
}

// GetImage decodes an image from TFExample from JPEG, PNG or GIF
func (e *TFExample) GetImage(name string) (image.Image, error) {
	b := e.GetBytesList(name)
	img, _, err := image.Decode(bytes.NewBuffer(b))
	return img, err
}
