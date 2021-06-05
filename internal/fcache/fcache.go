package fcache

import (
	"io"
	"time"

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

func Load(
	fs afero.Fs,
	path string,
	createAndSave func(src io.Reader, cache io.Writer) error,
	load func(cache io.Reader) error,
) error {
	mod, err := modTime(fs, path)
	if err != nil {
		return err
	}

	cachePath := path + ".cache"
	cacheMod, _ := modTime(fs, cachePath)

	if cacheMod.After(mod) {
		// log.Printf("Cache HIT: %s", path)
		cache, err := fs.Open(cachePath)
		if err == nil {
			defer cache.Close()
			if err = load(cache); err == nil {
				return nil
			}
		}
		// log.Print("Failed to load from cache, reverting to cache miss")
	}

	// log.Printf("Cache MISS: %s", path)
	src, err := fs.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	cache, err := fs.Create(cachePath)
	if err != nil {
		return err
	}
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
			_, err = cache.Write(data)
			if err != nil {
				return err
			}
			return nil
		},
		func(cache io.Reader) error { // load
			data, err := afero.ReadAll(cache)
			if err != nil {
				return err
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
