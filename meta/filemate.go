package meta

import (
	"github.com/KenianShi/filestore-server/db"
	"sort"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(fMeta FileMeta) { //todo 此处改成 filemeta 的所属方法比较好
	fileMetas[fMeta.FileSha1] = fMeta
}

func UpdateFileMetaDB(fMeta FileMeta) bool{
	return db.OnFileUploadFinished(fMeta.FileSha1,fMeta.FileName,fMeta.Location,fMeta.FileSize)
}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (*FileMeta,error){
	var fmeta FileMeta
	tfile,err := db.GetFileMeta(fileSha1)
	if err != nil {
		return &fmeta,err
	}
	fmeta.FileSha1 = tfile.FileHash
	fmeta.FileSize = tfile.FileSize.Int64
	fmeta.FileName = tfile.FileName.String
	fmeta.Location = tfile.FileAddr.String
	return &fmeta,nil
}


func GetLatestFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count] //todo 此处会存在数据越界的问题
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
