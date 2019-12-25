package loop

import (
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func FetchProfile(filePath string) ([]string, error) {
	data, err := ioutil.ReadFile(filePath)
	res := []string{}
	if err != nil {
		return []string{}, err
	}
	file := string(data)

	temp := strings.Split(file, "\n")

	for _, item := range temp {
		if strings.Contains(item, "[profile ") {
			item = strings.Replace(item, "[profile ", "", -1)
			item = strings.Replace(item, "]", "", -1)
			res = append(res, item)
		}
	}
	return res, nil
}

func FetchRegion(profileName string, regionName string) ([]string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profileName}))

	svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion(regionName),
	)
	awsRegions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	regions := make([]string, 0, len(awsRegions.Regions))
	for _, region := range awsRegions.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}
