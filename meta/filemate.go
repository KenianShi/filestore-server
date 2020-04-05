package meta

import (
	"github.com/KenianShi/filestore-server/db"
	"sort"
	"fmt"
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

func GetLatestFileMetasDB(count int) ([]FileMeta,error){
	tfiles,err := db.GetFileMetaList(count)
	if err != nil {
		fmt.Println("get LatestFileMeta From DB failed")
		return nil,err
	}
	var fileMetas []FileMeta
	for _,tfile := range tfiles{
		fileMeta := FileMeta{}
		fileMeta.FileSha1 = tfile.FileHash
		fileMeta.Location = tfile.FileAddr.String
		fileMeta.FileSize = tfile.FileSize.Int64
		fileMeta.FileName = tfile.FileName.String

		fileMetas = append(fileMetas,fileMeta)
	}
	return fileMetas,nil

	//for i:= 0;i<len(tfiles);i++{
	//	fileMeta := FileMeta{}
	//	fileMeta.FileSha1 = tfiles[i].FileHash
	//	fileMeta.FileName = tfiles[i].FileName.String
	//	fileMeta.FileSize = tfiles[i].FileSize.Int64
	//	fileMeta.Location = tfiles[i].FileAddr.String
	//	fileMetas = append(fileMetas,fileMeta)
	//}


}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
