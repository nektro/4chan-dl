package main

import (
	"bytes"
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
	"go.etcd.io/bbolt"
)

var (
	doneDir = "./data/"
	dbC     *bbolt.DB
)

func main() {

	flagBoards := pflag.StringArrayP("board", "b", []string{}, "/--board/ to download.")

	flagSaveDir := pflag.String("save-dir", "", "Path to a directory to save to.")
	flagConcurr := pflag.Int("concurrency", 10, "Maximum number of simultaneous downloads.")

	pflag.Parse()

	//

	if len(*flagSaveDir) > 0 {
		doneDir = *flagSaveDir
	}
	doneDir, _ = filepath.Abs(doneDir)
	doneDir += "/4chan.org"
	os.MkdirAll(doneDir, os.ModePerm)

	dbBDB, err := bbolt.Open(doneDir+"/comments.boltdb", 0666, nil)
	dbC = dbBDB
	util.DieOnError(err)

	//

	util.RunOnClose(onClose)
	mbpp.Init(*flagConcurr)

	//

	for _, item := range *flagBoards {
		grabBoard(item)
	}

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
	dbC.Close()
}

func grabBoard(board string) {
	mbpp.CreateJob("/"+board+"/", func(bar *mbpp.BarProxy) {
		req, _ := http.NewRequest(http.MethodGet, "https://p.4chan.org/4chan/board/"+board+"/catalog", nil)
		req.Header.Add("user-agent", "nektro/4chan-dl")
		res, _ := http.DefaultClient.Do(req)
		bys, _ := ioutil.ReadAll(res.Body)
		val, _ := fastjson.ParseBytes(bys)

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
	dir := doneDir + "/" + board + "/" + id
	m := false

	req, _ := http.NewRequest(http.MethodGet, "https://p.4chan.org/4chan/board/"+board+"/thread/"+id, nil)
	req.Header.Add("user-agent", "nektro/4chan-dl")
	res, _ := http.DefaultClient.Do(req)
	bys, _ := ioutil.ReadAll(res.Body)
	val, _ := fastjson.ParseBytes(bys)

	ar := val.GetArray("body", "posts")

	dbC.Update(func(tx *bbolt.Tx) error {
		b1, _ := tx.CreateBucketIfNotExists([]byte(board))
		b2, _ := b1.CreateBucketIfNotExists([]byte(id))
		b2.Put([]byte("title"), ar[0].GetStringBytes("sub"))
		for i, item := range ar {
			b2.Put([]byte("c"+strconv.Itoa(i)), bytes.Join(
				[][]byte{
					[]byte(strconv.FormatInt(int64(item.GetFloat64("no")), 10)),
					item.GetStringBytes("now"),
					item.GetStringBytes("name"),
					item.GetStringBytes("com"),
					item.GetStringBytes("tim"),
					item.GetStringBytes("filename"),
					item.GetStringBytes("ext"),
				},
				[]byte("\n"),
			))
		}
		return nil
	})

	for _, item := range ar {
		t := strconv.Itoa(item.GetInt("tim"))
		f := string(item.GetStringBytes("filename"))
		e := string(item.GetStringBytes("ext"))
		u := "https://i.4cdn.org/" + board + "/" + t + e

		if len(e) == 0 {
			continue
		}
		bar.AddToTotal(1)
		if !m {
			os.MkdirAll(dir, os.ModePerm)
			m = true
		}
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
		ar := val.GetArray("body", "boards")
		bar.AddToTotal(int64(len(ar)))
		for _, item := range ar {
			grabBoard(string(item.GetStringBytes("board")))
			bar.Increment(1)
		}
	})
}
