package main

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"reflect"
	"strconv"
)

func GetCheckpointId(vbId uint16, groupName string) string {
	key := Prefix + groupName + ":checkpoint:" + strconv.Itoa(int(vbId))
	crc := crc32.Checksum([]byte(fmt.Sprintf("%x", vbId)), crc32.IEEETable)
	return fmt.Sprintf("%v#%08x", key, crc)
}

func IsMetadata(data interface{}) bool {
	t := reflect.TypeOf(data)

	_, exist := t.FieldByName("Key")

	if exist {
		key := reflect.ValueOf(data).FieldByName("Key").Bytes()
		return bytes.HasPrefix(key, []byte(Prefix))
	}

	return false
}

func ChunkSlice(slice []uint16, chunks int) [][]uint16 {
	maxChunkSize := ((len(slice) - 1) / chunks) + 1
	numFullChunks := chunks - (maxChunkSize*chunks - len(slice))

	result := make([][]uint16, chunks)

	startIndex := 0

	for i := 0; i < chunks; i++ {
		endIndex := startIndex + maxChunkSize

		if i >= numFullChunks {
			endIndex--
		}

		result[i] = slice[startIndex:endIndex]

		startIndex = endIndex
	}

	return result
}