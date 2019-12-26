package main

import (
	"fmt"
	"os"

	goHomeDir "github.com/mitchellh/go-homedir"

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

	dataChan := make(chan data.Data)
	errChan := make(chan error)
	// APIコールを並行処理で実施
	for _, info := range infos {
		go DescribeEC2(info.Profile, info.Region, dataChan, errChan)
		select {
		case data := <-dataChan:
			fmt.Println(data)
		case err := <-errChan:
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
