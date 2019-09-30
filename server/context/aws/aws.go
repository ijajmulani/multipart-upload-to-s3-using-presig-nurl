package aws

type StartMultipartUploadArgs struct {
	ContentType string
	BucketType  string
}

type StartMultipartUploadResp struct {
	UploadID string
	MediaID  string
	Error    string
}

type GetUploadURLArgs struct {
	BucketType  string
	ContentType string
}

type GetUploadURLResp struct {
	URL   string
	Error string
}

type GetMultipartUploadURLArgs struct {
	MediaID    string
	BucketType string
	PartNumber int64
	UploadID   string
}

type GetMultipartUploadURLResp struct {
	URL   string
	Error string
}

type CompleteMultipartUploadArgs struct {
	MediaID    string
	BucketType string
	Parts      []CompletedParts
	UploadID   string
}

type CompleteMultipartUploadResp struct {
	Status bool
	Error  string
}

type CompletedParts struct {
	ETag       string
	PartNumber int64
}
