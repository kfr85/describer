package main

import (
	"log"
	"os"

	goHomeDir "github.com/mitchellh/go-homedir"
	"golang.org/x/sync/errgroup"

	"data"
	"loop"
)

func main() {
	// パスの取得
	homeDirPath, _ := goHomeDir.Dir()
	path := homeDirPath + "/.aws/config"

	// アカウント一覧の取得
	profiles, err := loop.FetchProfile(path)
	if err != nil {
		panic(nil)
	}

	// リージョン一覧の取得
	regions, err := loop.FetchRegion(profiles[0], "ap-northeast-1")
	if err != nil {
		panic(nil)
	}

	// ループリストを作成
	var infos data.Infomations
	for _, profile := range profiles {
		for _, region := range regions {
			infos = append(infos, data.Infomation{Profile: profile, Region: region})
		}
	}

	// APIコールを並行処理で実施
	sem := make(chan struct{}, 50)
	eg := errgroup.Group{}
	for _, info := range infos {
		info := info
		sem <- struct{}{}
		eg.Go(func() error {
			return DescribeEC2(info.Profile, info.Region)
		})
		<-sem
	}

	// 並行処理中にエラーが出た場合は異常終了する
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
