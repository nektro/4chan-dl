package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
	"github.com/spf13/pflag"
	"github.com/valyala/fastjson"
)

var (
	DoneDir = "./data/"
)

func main() {

	flagBoards := pflag.StringArrayP("board", "b", []string{}, "/--board/ to download.")

	flagSaveDir := pflag.String("save-dir", "", "Path to a directory to save to.")
	flagConcurr := pflag.Int("concurrency", 10, "Maximum number of simultaneous downloads.")

	pflag.Parse()

	//

	if len(*flagSaveDir) > 0 {
		DoneDir = *flagSaveDir
	}
	DoneDir, _ = filepath.Abs(DoneDir)
	DoneDir += "/4chan.org"
	os.MkdirAll(DoneDir, os.ModePerm)

	//

	util.RunOnClose(onClose)
	mbpp.Init(*flagConcurr)

	//

	for _, item := range *flagBoards {
		grabBoard(item)
	}

	//

	if len(*flagBoards) == 0 {
		grabAllBoards()
	}

	//

	mbpp.Wait()
	time.Sleep(time.Second)
	onClose()
}

func onClose() {
	util.Log(mbpp.GetCompletionMessage())
}

func grabBoard(board string) {
	mbpp.CreateJob("/"+board+"/", func(bar *mbpp.BarProxy) {
		req, _ := http.NewRequest(http.MethodGet, "https://p.4chan.org/4chan/board/"+board+"/catalog", nil)
		req.Header.Add("user-agent", "nektro/4chan-dl")
		res, _ := http.DefaultClient.Do(req)
		bys, _ := ioutil.ReadAll(res.Body)
		val, _ := fastjson.ParseBytes(bys)
		//
		ar1 := val.GetArray("body")
		ids := []string{}
		for _, item := range ar1 {
			ar2 := item.GetArray("threads")
			for _, jtem := range ar2 {
				ids = append(ids, strconv.Itoa(jtem.GetInt("no")))
			}
		}
		for _, item := range ids {
			go grabThread(board, item, bar)
			time.Sleep(time.Second / 4)
		}
	})
}

func grabThread(board, id string, bar *mbpp.BarProxy) {
	dir := DoneDir + "/" + board + "/" + id
	m := false
	//
	req, _ := http.NewRequest(http.MethodGet, "https://p.4chan.org/4chan/board/"+board+"/thread/"+id, nil)
	req.Header.Add("user-agent", "nektro/4chan-dl")
	res, _ := http.DefaultClient.Do(req)
	bys, _ := ioutil.ReadAll(res.Body)
	val, _ := fastjson.ParseBytes(bys)
	//
	ar := val.GetArray("body", "posts")
	bar.AddToTotal(int64(len(ar)))
	for _, item := range ar {
		t := strconv.Itoa(item.GetInt("tim"))
		f := string(item.GetStringBytes("filename"))
		e := string(item.GetStringBytes("ext"))
		u := "https://i.4cdn.org/" + board + "/" + t + e
		//
		if len(e) == 0 {
			bar.Increment(1) // this post does not have an attachment
			continue
		}
		if !m {
			os.MkdirAll(dir, os.ModePerm)
			m = true
		}
		//
		go mbpp.CreateDownloadJob(u, dir+"/"+t+"_"+f+e, bar)
	}
}

func grabAllBoards() {
	mbpp.CreateJob("4chan.org", func(bar *mbpp.BarProxy) {
		req, _ := http.NewRequest(http.MethodGet, "https://p.4chan.org/4chan/boards", nil)
		req.Header.Add("user-agent", "nektro/4chan-dl")
		res, _ := http.DefaultClient.Do(req)
		bys, _ := ioutil.ReadAll(res.Body)
		val, _ := fastjson.ParseBytes(bys)
		//
		ar := val.GetArray("body", "boards")
		bar.AddToTotal(int64(len(ar)))
		for _, item := range ar {
			grabBoard(string(item.GetStringBytes("board")))
			bar.Increment(1)
		}
	})
}
