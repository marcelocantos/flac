package fcache

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/pierrec/lz4"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
)

func modTime(fs afero.Fs, path string) (time.Time, error) {
	stat, err := fs.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return stat.ModTime(), nil
}

func logCacheTiming() func(mode, path string) {
	// start := time.Now()
	return func(mode, path string) {
		// end := time.Now()
		// log.Printf("Cache %s (%.2f ms): %s",
		// 	mode,
		// 	float64(end.Sub(start))/float64(time.Millisecond),
		// 	path)
	}
}

func Load(
	fs afero.Fs,
	path string,
	createAndSave func(src io.Reader, cache io.Writer) error,
	load func(cache io.Reader) error,
) (err error) {
	mod, err := modTime(fs, path)
	if err != nil {
		return err
	}

	cachePath := path + ".cache"
	cacheMod, _ := modTime(fs, cachePath)

	if cacheMod.After(mod) {
		defer logCacheTiming()("HIT", path)
		cache, err := fs.Open(cachePath)
		if err == nil {
			defer cache.Close()
			if err = load(cache); err == nil {
				return nil
			}
		}
		log.Print("Cache load failed, reverting to cache miss")
	}

	defer logCacheTiming()("MISS", path)
	src, err := fs.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	cache, err := fs.Create(cachePath)
	if err != nil {
		cache.Close()
		fs.Remove(cachePath)
		return err
	}
	defer func() {
		if err != nil {
			fs.Remove(cachePath)
		}
	}()
	defer cache.Close()

	return createAndSave(src, cache)
}

func Proto(
	fs afero.Fs,
	path string,
	target proto.Message,
	create func(src io.Reader) error,
) error {
	err := Load(
		fs, path,
		func(src io.Reader, cache io.Writer) error { // createAndSave
			if err := create(src); err != nil {
				return err
			}
			data, err := proto.Marshal(target)
			if err != nil {
				return err
			}
			w := lz4.NewWriter(cache)
			defer w.Close()
			_, err = w.Write(data)
			if err != nil {
				return err
			}
			return nil
		},
		func(cache io.Reader) error { // load
			data, err := ioutil.ReadAll(lz4.NewReader(cache))
			if err != nil {
				return err
			}
			if len(data) == 0 {
				return errors.New("empty cache")
			}
			err = proto.Unmarshal(data, target)
			return err
		},
	)
	if err != nil {
		return err
	}
	return err
}
