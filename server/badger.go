package server

import (
	"os"
	"path"

	"github.com/dgraph-io/badger/v3"
	libbadger "github.com/otamoe/go-library/badger"
	libviper "github.com/otamoe/go-library/viper"
	"github.com/spf13/viper"
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := path.Join(userHomeDir, ".vptun/server/badger")

	defaultOptions := libbadger.DefaultOptions()
	libviper.SetDefault("badger.indexDir", path.Join(baseDir, "index"), "Badger index dir")
	libviper.SetDefault("badger.valueDir", path.Join(baseDir, "value"), "Badger value dir")
	libviper.SetDefault("badger.memTableSize", defaultOptions.MemTableSize, "Badger mem table size")
	libviper.SetDefault("badger.blockCacheSize", defaultOptions.BlockCacheSize, "Badger block cache size")
	libviper.SetDefault("badger.indexCacheSize", defaultOptions.IndexCacheSize, "Badger index cache size")
}

func BadgerIndexDirOption() (out libbadger.OutOption) {
	out.Option = func(b badger.Options) (badger.Options, error) {
		b.Dir = viper.GetString("badger.indexDir")
		return b, nil
	}
	return
}

func BadgerValueDirOption() (out libbadger.OutOption) {
	out.Option = func(b badger.Options) (badger.Options, error) {
		b.ValueDir = viper.GetString("badger.valueDir")
		return b, nil
	}
	return
}
func BadgerMemTableSizeOption() (out libbadger.OutOption) {
	out.Option = func(b badger.Options) (badger.Options, error) {
		if val := viper.GetInt64("badger.memTableSize"); val > 1024*1024*4 {
			b.MemTableSize = val
		}
		return b, nil
	}
	return
}

func BadgerBlockCacheSizeOption() (out libbadger.OutOption) {
	out.Option = func(b badger.Options) (badger.Options, error) {
		if val := viper.GetInt64("badger.blockCacheSize"); val > 1024*1024*16 {
			b.BlockCacheSize = val
		}
		return b, nil
	}
	return
}

func BadgerIndexCacheSizeOption() (out libbadger.OutOption) {
	out.Option = func(b badger.Options) (badger.Options, error) {
		if val := viper.GetInt64("badger.indexCacheSize"); val > 1024*1024*16 {
			b.IndexCacheSize = val
		}
		return b, nil
	}
	return
}
