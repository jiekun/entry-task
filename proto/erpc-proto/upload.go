// @Author: 2014BDuck
// @Date: 2021/8/5

package erpc_proto

type UploadService interface {
	UploadFile(*UploadRequest) (*UploadReply, error)
}

type UploadRequest struct {
	FileType uint32
	FileName string
	Content  []byte
}

type UploadReply struct {
	FileUrl  string
	FileName string
}
