package dirserver

import (
	"bytes"
	"fmt"

	"github.com/00000vish/filebeam/internals/config"
	"github.com/00000vish/filebeam/internals/dirserver/dirlisting"
	"github.com/00000vish/filebeam/internals/endpoint"
)

type DirServer struct {
	config   *config.Config
	endpoint *endpoint.Endpoint
	listing  *dirlisting.DirListing
}

func (d *DirServer) Run(config *config.Config, endpoint *endpoint.Endpoint) {
	d.config = config
	d.endpoint = endpoint

	d.listing = dirlisting.Create(d.config)

	d.endpoint.Connect(d.connected, d.reciever)
}

func (d *DirServer) connected() {
	d.sendListing()
}

func (d *DirServer) reciever(data []byte) {
	if bytes.Equal(DIR_LISTING, data[:PROTOCOL_SIZE]) {
		listing, _ := dirlisting.DeSerializeFromJson(data[PROTOCOL_SIZE:])
		difference := d.listing.GetDifference(listing);
		d.beginTranmit(difference);
		return
	}
	fmt.Println(string(data[PROTOCOL_SIZE:]))
}

func (d *DirServer) sendListing() {
	data := d.listing.SerializeToJson()
	d.sendData(DIR_LISTING, data)
}

func (d *DirServer) sendData(protocol []byte, data []byte) {
	toSend := append(protocol, data...)
	d.endpoint.Send(toSend)
}

func (d *DirServer) beginTranmit(dirListing *dirlisting.DirListing) {
	for _, file := range(dirListing.Files){
		file.Read()
	}
}