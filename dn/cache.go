package dn

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/BakerHub/trivial/check"
	"github.com/BakerHub/trivial/fs"
)

type CacheMeta struct {
	Saved string `json:"saved,omitempty"`
	ETag  string `json:"etag,omitempty"`
	MD5   string `json:"md5,omitempty"`
}

type Cache struct {
	Files map[string]*CacheMeta
	wg    sync.WaitGroup
	mux   sync.Mutex

	metaFile string
}

func NewCache(metaFile string) *Cache {
	c := &Cache{
		Files:    make(map[string]*CacheMeta),
		metaFile: metaFile,
	}
	c.loadMeta()
	return c
}

func tempFileFor(pathname string) *os.File {
	pattern := path.Base(pathname) + "*" + path.Ext(pathname)
	file, err := ioutil.TempFile("", pattern)
	check.Check(err)
	return file
}

func shouldRemove(file string) {
	if !fs.Exists(file) {
		return
	}

	err := os.Remove(file)
	if err != nil {
		log.Println(err)
	}
}

func shouldClose(file io.Closer) {
	err := file.Close()
	if err != nil {
		log.Println(err)
	}
}

func shouldMove(oldLocation, newLocation string) {
	var err error
	if fs.Exists(newLocation) {
		err = os.Remove(newLocation)

	} else {
		err = os.MkdirAll(path.Dir(newLocation), 0777)
	}
	if err != nil {
		log.Println(err)
	}

	err = os.Rename(oldLocation, newLocation)
	if err != nil {
		log.Println(err)
	}
}

func downloadFile(url string, savePath string) string {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer shouldClose(resp.Body)

	out := tempFileFor(savePath)
	defer shouldRemove(out.Name())

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	shouldClose(out)

	if err != nil {
		log.Println(err)
		return ""
	}

	fs.EnsureDir(path.Dir(savePath))
	shouldMove(out.Name(), savePath)
	return resp.Header.Get("ETag")
}

func (cache *Cache) Check(url string, cachePath string) (bool, string) {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	if local, ok := cache.Files[url]; ok {
		return true, local.Saved
	}

	cache.Files[url] = &CacheMeta{Saved: cachePath}
	exists := fs.Exists(cachePath)

	return exists, cachePath
}

func (cache *Cache) Fetch(url string, cachePath string) string {
	if ok, localPath := cache.Check(url, cachePath); ok {
		return localPath
	}

	cache.StartDownload(url, cachePath)

	return cachePath
}

func (cache *Cache) addMeta(url string, cachePath string, etag string) {
	if len(etag) == 0 {
		return
	}

	meta, exits := cache.Files[url]

	if !exits {
		meta = &CacheMeta{Saved: cachePath}
		cache.Files[url] = meta
	}

	hex, err := fs.NewFileHashMD5().FromFile(cachePath)
	if err != nil {
		return
	}
	meta.MD5 = hex
	meta.ETag = etag
	meta.Saved = cachePath
}

func (cache *Cache) StartDownload(url string, cachePath string) {
	cache.wg.Add(1)
	go func(url string, cachePath string) {
		defer cache.wg.Done()
		etag := downloadFile(url, cachePath)
		cache.addMeta(url, cachePath, etag)
	}(url, cachePath)
}

func (cache *Cache) loadMeta() {
	if !fs.Exists(cache.metaFile) {
		return
	}

	file, _ := ioutil.ReadFile(cache.metaFile)
	_ = json.Unmarshal(file, &cache.Files)
}

func (cache *Cache) saveMeta() {
	file, _ := json.MarshalIndent(cache.Files, "", "  ")
	_ = ioutil.WriteFile(cache.metaFile, file, 0644)
}

func (cache *Cache) Wait() {
	cache.wg.Wait()
	cache.saveMeta()
}
