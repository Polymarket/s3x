/*
 * MinIO Cloud Storage, (C) 2016, 2017 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"io"
	"time"

	"github.com/RTradeLtd/s3x/pkg/hash"
	"github.com/RTradeLtd/s3x/pkg/madmin"
)

// BackendType - represents different backend types.
type BackendType int

// Enum for different backend types.
const (
	Unknown BackendType = iota
	// Filesystem backend.
	BackendFS
	// Multi disk BackendErasure (single, distributed) backend.
	BackendErasure
	// Gateway backend.
	BackendGateway
	// Add your own backend.
)

// StorageInfo - represents total capacity of underlying storage.
type StorageInfo struct {
	Used []uint64 // Used total used per disk.

	Total []uint64 // Total disk space per disk.

	Available []uint64 // Total disk space available per disk.

	MountPaths []string // Disk mountpoints

	// Backend type.
	Backend struct {
		// Represents various backend types, currently on FS, Erasure and Gateway
		Type BackendType

		// Following fields are only meaningful if BackendType is Gateway.
		GatewayOnline bool

		// Following fields are only meaningful if BackendType is Erasure.
		OnlineDisks      madmin.BackendDisks // Online disks during server startup.
		OfflineDisks     madmin.BackendDisks // Offline disks during server startup.
		StandardSCData   int                 // Data disks for currently configured Standard storage class.
		StandardSCParity int                 // Parity disks for currently configured Standard storage class.
		RRSCData         int                 // Data disks for currently configured Reduced Redundancy storage class.
		RRSCParity       int                 // Parity disks for currently configured Reduced Redundancy storage class.

		// List of all disk status, this is only meaningful if BackendType is Erasure.
		Sets [][]madmin.DriveInfo
	}
}

// objectHistogramInterval is an interval that will be
// used to report the histogram of objects data sizes
type objectHistogramInterval struct {
	name       string
	start, end int64
}

// ObjectsHistogramIntervals is the list of all intervals
// of object sizes to be included in objects histogram.
var ObjectsHistogramIntervals = []objectHistogramInterval{
	{"LESS_THAN_1024_B", -1, 1024 - 1},
	{"BETWEEN_1024_B_AND_1_MB", 1024, 1024*1024 - 1},
	{"BETWEEN_1_MB_AND_10_MB", 1024 * 1024, 1024*1024*10 - 1},
	{"BETWEEN_10_MB_AND_64_MB", 1024 * 1024 * 10, 1024*1024*64 - 1},
	{"BETWEEN_64_MB_AND_128_MB", 1024 * 1024 * 64, 1024*1024*128 - 1},
	{"BETWEEN_128_MB_AND_512_MB", 1024 * 1024 * 128, 1024*1024*512 - 1},
	{"GREATER_THAN_512_MB", 1024 * 1024 * 512, -1},
}

// DataUsageInfo represents data usage stats of the underlying Object API
type DataUsageInfo struct {
	// The timestamp of when the data usage info is generated
	LastUpdate time.Time `json:"lastUpdate"`

	ObjectsCount uint64 `json:"objectsCount"`
	// Objects total size
	ObjectsTotalSize      uint64            `json:"objectsTotalSize"`
	ObjectsSizesHistogram map[string]uint64 `json:"objectsSizesHistogram"`

	BucketsCount uint64            `json:"bucketsCount"`
	BucketsSizes map[string]uint64 `json:"bucketsSizes"`
}

// BucketInfo - represents bucket metadata.
type BucketInfo struct {
	// Name of the bucket.
	Name string

	// Date and time when the bucket was created.
	Created time.Time
}

// ObjectInfo - represents object metadata.
type ObjectInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Name string

	// Date and time when the object was last modified.
	ModTime time.Time

	// Total object size.
	Size int64

	// IsDir indicates if the object is prefix.
	IsDir bool

	// Hex encoded unique entity tag of the object.
	ETag string

	// A standard MIME type describing the format of the object.
	ContentType string

	// Specifies what content encodings have been applied to the object and thus
	// what decoding mechanisms must be applied to obtain the object referenced
	// by the Content-Type header field.
	ContentEncoding string

	// Date and time at which the object is no longer able to be cached
	Expires time.Time

	// CacheStatus sets status of whether this is a cache hit/miss
	CacheStatus CacheStatusType
	// CacheLookupStatus sets whether a cacheable response is present in the cache
	CacheLookupStatus CacheStatusType

	// Specify object storage class
	StorageClass string

	// User-Defined metadata
	UserDefined map[string]string

	// User-Defined object tags
	UserTags string

	// List of individual parts, maximum size of upto 10,000
	Parts []ObjectPartInfo `json:"-"`

	// Implements writer and reader used by CopyObject API
	Writer       io.WriteCloser `json:"-"`
	Reader       *hash.Reader   `json:"-"`
	PutObjReader *PutObjReader  `json:"-"`

	metadataOnly bool

	// Date and time when the object was last accessed.
	AccTime time.Time

	// backendType indicates which backend filled this structure
	backendType BackendType
}

// ListPartsInfo - represents list of all parts.
type ListPartsInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Object string

	// Upload ID identifying the multipart upload whose parts are being listed.
	UploadID string

	// The class of storage used to store the object.
	StorageClass string

	// Part number after which listing begins.
	PartNumberMarker int

	// When a list is truncated, this element specifies the last part in the list,
	// as well as the value to use for the part-number-marker request parameter
	// in a subsequent request.
	NextPartNumberMarker int

	// Maximum number of parts that were allowed in the response.
	MaxParts int

	// Indicates whether the returned list of parts is truncated.
	IsTruncated bool

	// List of all parts.
	Parts []PartInfo

	// Any metadata set during InitMultipartUpload, including encryption headers.
	UserDefined map[string]string

	EncodingType string // Not supported yet.
}

// Lookup - returns if uploadID is valid
func (lm ListMultipartsInfo) Lookup(uploadID string) bool {
	for _, upload := range lm.Uploads {
		if upload.UploadID == uploadID {
			return true
		}
	}
	return false
}

// ListMultipartsInfo - represnets bucket resources for incomplete multipart uploads.
type ListMultipartsInfo struct {
	// Together with upload-id-marker, this parameter specifies the multipart upload
	// after which listing should begin.
	KeyMarker string

	// Together with key-marker, specifies the multipart upload after which listing
	// should begin. If key-marker is not specified, the upload-id-marker parameter
	// is ignored.
	UploadIDMarker string

	// When a list is truncated, this element specifies the value that should be
	// used for the key-marker request parameter in a subsequent request.
	NextKeyMarker string

	// When a list is truncated, this element specifies the value that should be
	// used for the upload-id-marker request parameter in a subsequent request.
	NextUploadIDMarker string

	// Maximum number of multipart uploads that could have been included in the
	// response.
	MaxUploads int

	// Indicates whether the returned list of multipart uploads is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of multipart uploads exceeds the limit allowed or specified
	// by max uploads.
	IsTruncated bool

	// List of all pending uploads.
	Uploads []MultipartInfo

	// When a prefix is provided in the request, The result contains only keys
	// starting with the specified prefix.
	Prefix string

	// A character used to truncate the object prefixes.
	// NOTE: only supported delimiter is '/'.
	Delimiter string

	// CommonPrefixes contains all (if there are any) keys between Prefix and the
	// next occurrence of the string specified by delimiter.
	CommonPrefixes []string

	EncodingType string // Not supported yet.
}

// ListObjectsInfo - container for list objects.
type ListObjectsInfo struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response is true),
	// you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: AWS S3 returns NextMarker only if you have delimiter request parameter specified,
	//       MinIO always returns NextMarker.
	NextMarker string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

// ListObjectsV2Info - container for list objects version 2.
type ListObjectsV2Info struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: This element is returned only if you have delimiter request parameter
	// specified.
	ContinuationToken     string
	NextContinuationToken string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

// PartInfo - represents individual part metadata.
type PartInfo struct {
	// Part number that identifies the part. This is a positive integer between
	// 1 and 10,000.
	PartNumber int

	// Date and time at which the part was uploaded.
	LastModified time.Time

	// Entity tag returned when the part was initially uploaded.
	ETag string

	// Size in bytes of the part.
	Size int64

	// Decompressed Size.
	ActualSize int64
}

// MultipartInfo - represents metadata in progress multipart upload.
type MultipartInfo struct {
	// Object name for which the multipart upload was initiated.
	Object string

	// Unique identifier for this multipart upload.
	UploadID string

	// Date and time at which the multipart upload was initiated.
	Initiated time.Time

	StorageClass string // Not supported yet.
}

// CompletePart - represents the part that was completed, this is sent by the client
// during CompleteMultipartUpload request.
type CompletePart struct {
	// Part number identifying the part. This is a positive integer between 1 and
	// 10,000
	PartNumber int

	// Entity tag returned when the part was uploaded.
	ETag string
}

// CompletedParts - is a collection satisfying sort.Interface.
type CompletedParts []CompletePart

func (a CompletedParts) Len() int           { return len(a) }
func (a CompletedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CompletedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

// CompleteMultipartUpload - represents list of parts which are completed, this is sent by the
// client during CompleteMultipartUpload request.
type CompleteMultipartUpload struct {
	Parts []CompletePart `xml:"Part"`
}
