package dirlisting

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/00000vish/filebeam/internals/config"
	"github.com/00000vish/filebeam/internals/logger"
	"github.com/00000vish/filebeam/internals/dirserver/dirfileinfo"
)

type DirListing struct {
	Path           string
	Files          []*dirfileinfo.DirFileInfo
	SubDirectories []*DirListing
}

func Create(config *config.Config) *DirListing {
	listing := DirListing{
		Path: config.Directory,
	}

	listing.Scan()

	return &listing
}

func (d *DirListing) IsEmpty() bool {
	if(len(d.Files) != 0){
		return false;
	}

	if(len(d.SubDirectories) != 0){
		return false;
	}

	return true;
}

func DeSerializeFromJson(data []byte) (*DirListing, error) {
	res := DirListing{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		logger.Log(err.Error())
		return nil, err
	}
	return &res, nil
}

func (d *DirListing) SerializeToJson() []byte {
	json, err := json.Marshal(d)
	if err != nil {
		logger.Log(err.Error())
	}
	return json
}

func (d *DirListing) Scan() {
	err := filepath.Walk(d.Path, d.walkDirectory)
	if err != nil {
		log.Println(err)
	}
}

func (d *DirListing) ToString() {
	fmt.Println(fmt.Sprintf("Directory : %s", d.Path))
	for _, file := range d.Files {
		fmt.Println(fmt.Sprintf("\t\t%s", file.FileName))
	}
	for _, subDir := range d.SubDirectories {
		subDir.ToString()
	}
}

func (d *DirListing) IsEquals(other *DirListing) bool {
	for _, localFile := range d.Files{
		found := false
		for _, remoteFile := range other.Files{
			if localFile.IsEquals(remoteFile){
				found = true
				break;
			} 
		}
		if !found {
			return false
		}
	}

	for _, localSubDir := range d.SubDirectories{
		found := false
		for _, remoteSubDir := range other.SubDirectories{
			if localSubDir.IsEquals(remoteSubDir){
				found = true;
				break;
			}
			if !found {
				return false
			} 
		}
	}

	return true
}

func (d *DirListing) GetDifference(other *DirListing) *DirListing {
	diffListing := DirListing{}

	if(d.IsEquals(other)){
		return &diffListing;
	}

	for _, localFile := range d.Files{
		found := false
		for _, remoteFile := range other.Files{
			if localFile.IsEquals(remoteFile){
				found = true
				break;
			} 
		}
		if !found {
			diffListing.Files = append(diffListing.Files, localFile)
		}
	}

	//TODO: implement subdirectories
	
	return &diffListing
}

func (d *DirListing) walkDirectory(path string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return err
	}

	listing := d

	fileName := info.Name()

	fileDir := strings.Replace(path, fileName, "", 1)
	if len(fileDir) == 0 {
		return nil
	}

	fileDir = fileDir[:len(fileDir)-1]

	if fileDir != listing.Path {
		for _, subListing := range listing.SubDirectories {
			if subListing.Path == fileDir {
				listing = subListing
				break
			}
		}

		if listing == d {
			listing = &DirListing{
				Path: fileDir,
			}
			d.SubDirectories = append(d.SubDirectories, listing)
		}
	}

	fileInfo := dirfileinfo.DirFileInfo{}
	fileInfo.FileName = fileName
	fileInfo.Path = path

	listing.Files = append(listing.Files, &fileInfo)

	return nil
}
