package main

import (
	"encoding/json"
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
		panic(err)
	}

	// リージョン一覧の取得
	regions, err := loop.FetchRegion(profiles[0], "ap-northeast-1")
	if err != nil {
		panic(err)
	}

	// ループリストを作成
	var infos data.Infomations
	for _, profile := range profiles {
		for _, region := range regions {
			infos = append(infos, data.Infomation{Profile: profile, Region: region})
		}
	}

	dataChan := make(chan data.Data, 50)
	errChan := make(chan error, 1)
	// APIコールを並行処理で実施
	go func() {
		defer close(dataChan)
		defer close(errChan)
		for _, info := range infos {
			DescribeEC2(info.Profile, info.Region, dataChan, errChan)
		}
	}()

	// 並行処理でデータを受け取ったものから出力処理
	for {
		select {
		case data := <-dataChan:
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("marshal error! profile: %s, region: %s -> %s", data.Infomation.Profile, data.Infomation.Region, err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", jsonData)
		case err := <-errChan:
			// errChanをクローズするので一度はここを通るため、nil判定をする
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				return
			}
		}
	}
}
