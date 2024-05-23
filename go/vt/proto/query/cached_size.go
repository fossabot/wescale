/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by Sizegen. DO NOT EDIT.

package query

import hack "github.com/wesql/wescale/go/hack"

func (cached *BindVariable) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(96)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Value []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Value)))
	}
	// field Values []*vitess.io/vitess/go/vt/proto/query.Value
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Values)) * int64(8))
		for _, elem := range cached.Values {
			size += elem.CachedSize(true)
		}
	}
	return size
}
func (cached *Field) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(160)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Name string
	size += hack.RuntimeAllocSize(int64(len(cached.Name)))
	// field Table string
	size += hack.RuntimeAllocSize(int64(len(cached.Table)))
	// field OrgTable string
	size += hack.RuntimeAllocSize(int64(len(cached.OrgTable)))
	// field Database string
	size += hack.RuntimeAllocSize(int64(len(cached.Database)))
	// field OrgName string
	size += hack.RuntimeAllocSize(int64(len(cached.OrgName)))
	// field ColumnType string
	size += hack.RuntimeAllocSize(int64(len(cached.ColumnType)))
	return size
}
func (cached *QueryWarning) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(64)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Message string
	size += hack.RuntimeAllocSize(int64(len(cached.Message)))
	return size
}
func (cached *Target) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(96)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Keyspace string
	size += hack.RuntimeAllocSize(int64(len(cached.Keyspace)))
	// field Shard string
	size += hack.RuntimeAllocSize(int64(len(cached.Shard)))
	// field Cell string
	size += hack.RuntimeAllocSize(int64(len(cached.Cell)))
	return size
}
func (cached *Value) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(80)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Value []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Value)))
	}
	return size
}
