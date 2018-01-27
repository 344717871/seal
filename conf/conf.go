package conf

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/yaml"
)

var GlobalConfInfo ConfInfo

type RtmpConfInfo struct {
	Listen    string `yaml:"listen"`
	TimeOut   uint32 `yaml:"timeout"`
	ChunkSize uint32 `yaml:"chunk_size"`
}

type ConfInfo struct {
	Rtmp RtmpConfInfo `yaml:"rtmp"`
}

func (t *ConfInfo) Loads(c string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

	}()

	var f *os.File
	if f, err = os.Open(c); err != nil {
		log.Println("Open config failed, err is", err)
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("config file loads failed, ", err.Error())
		return err
	}

	err = yaml.Unmarshal(data, t)
	if err != nil {
		log.Println("error:", err.Error())
		return err
	}

	return nil
}

func (c *ConfInfo) Default() {
	c.Rtmp.Listen = "1935"
	c.Rtmp.TimeOut = 30
	c.Rtmp.ChunkSize = 6000
}