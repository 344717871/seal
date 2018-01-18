package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (rtmp *RtmpSession) Connect() (err error) {

	var chunk *ChunkStruct
	//todo. should expect what kind msg. if recv not a expected one ,maybe shoule
	//response it first
	for {
		err, chunk = rtmp.RecvMsg()
		if err != nil {
			return
		}

		err = rtmp.DecodeMsg(chunk)
		if err != nil {
			return
		}

		if "Amf0CommandConnectPkg" == chunk.decodeResultType {
			break
		}
	}

	connectPkg := chunk.decodeResult.(Amf0CommandConnectPkg)

	log.Println("rtmp connect result: ", connectPkg)

	tcUrlValue := connectPkg.Amf0ObjectsGetProperty("tcUrl")
	if nil == tcUrlValue {
		err = fmt.Errorf("tcUrl is nil.")
		return
	}

	var rtmpUrlData RtmpUrlData
	err = rtmpUrlData.ParseUrl(tcUrlValue.(string))
	if err != nil {
		return
	}

	err = rtmpUrlData.Discover()
	if err != nil {
		return
	}

	return
}

type RtmpUrlData struct {
	schema string
	host   string
	port   uint16
	app    string
	stream string
	token  string
}

//format: rtmp://127.0.0.1:1935/live/test?token=abc123
func (urlData *RtmpUrlData) ParseUrl(url string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	//url is not contain the stream with token.
	var urlTmp string

	urlTmp = strings.Replace(url, "://", " ", 1)
	urlTmp = strings.Replace(urlTmp, ":", " ", 1)
	urlTmp = strings.Replace(urlTmp, "/", " ", 1)

	urlSplit := strings.Split(urlTmp, " ")

	if 3 == len(urlSplit) {
		//no port, use the default 1935
		urlData.schema = urlSplit[0]
		urlData.host = urlSplit[1]
		urlData.port = 1935
		urlData.app = urlSplit[2]
	} else if 4 == len(urlSplit) {
		urlData.schema = urlSplit[0]
		urlData.host = urlSplit[1]
		port, ok := strconv.Atoi(urlSplit[2])
		if nil == ok {
			//the port is not default
			if port > 0 && port < 65536 {
				urlData.port = uint16(port)
			} else {
				err = fmt.Errorf("tcUrl port format is error, port=", port)
				return
			}

		} else {
			err = fmt.Errorf("tcurl format error when convert port format, err=", ok)
			return
		}
		urlData.app = urlSplit[3]
	} else {
		err = fmt.Errorf("tcUrl format is error. tcUrl=", url)
		return
	}

	return
}

func (urlData *RtmpUrlData) Discover() (err error) {
	if 0 == len(urlData.schema) ||
		0 == len(urlData.host) ||
		0 == len(urlData.app) {
		err = fmt.Errorf("discover url data failed. url data=", urlData)
		return
	}
	return
}
