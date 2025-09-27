package dirfileinfo

import (
	"encoding/json"
	"os"
	"encoding/base64"

	"github.com/00000vish/filebeam/internals/logger"
)


type DirFileInfo struct {
	Path     string
	FileName string
	Hash     string
	Data 	 string
}

func (f *DirFileInfo) IsEquals(other *DirFileInfo) bool {
	if f.Hash != other.Hash {
		return false
	}
	
	if f.FileName != other.FileName {
		return false
	}

	return true
}

func (f *DirFileInfo) Read()  {
	data, err := os.ReadFile(f.Path)
	if(err == nil){
		f.Data =  base64.StdEncoding.EncodeToString(data)
	}
}

func DeSerializeFromJson(data []byte) (*DirFileInfo, error) {
	res := DirFileInfo{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		logger.Log(err.Error())
		return nil, err
	}
	return &res, nil
}

func (d *DirFileInfo) SerializeToJson() []byte {
	json, err := json.Marshal(d)
	if err != nil {
		logger.Log(err.Error())
	}
	return json
}