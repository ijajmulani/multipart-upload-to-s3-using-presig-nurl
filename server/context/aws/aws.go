package aws

type StartUploadArgs struct {
	ContentType string
	BucketType  string
}

type StartUploadResp struct {
	UploadID string
	MediaID  string
	Error    string
}

type GetUploadURLArgs struct {
	MediaID    string
	BucketType string
	PartNumber int64
	UploadID   string
}

type GetUploadURLResp struct {
	URL   string
	Error string
}

type CompleteUploadURLArgs struct {
	MediaID    string
	BucketType string
	Parts      []CompletedParts
	UploadID   string
}

type CompleteUploadURLResp struct {
	Status bool
	Error  string
}

type CompletedParts struct {
	ETag       string
	PartNumber int64
}
