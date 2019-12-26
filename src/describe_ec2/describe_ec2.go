package main

import (
	"fmt"

	"data"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// DescribeEC2 retiurn error
func DescribeEC2(profileName string, regionName string, dataChan chan data.Data, errChan chan error) {
	var ret []*ec2.Reservation
	var data data.Data

	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profileName}))
	svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion(regionName),
	)

	configInput := new(ec2.DescribeInstancesInput)

	for {
		req, err := svc.DescribeInstances(configInput)
		if err != nil {
			errChan <- fmt.Errorf("describe error! profile: %s, region: %s -> %s", profileName, regionName, err)
			return
		}

		ret = append(ret, req.Reservations...)

		if req.NextToken == nil {
			break
		}
		configInput = configInput.SetNextToken(*req.NextToken)
	}

	data.Result = ret
	data.Infomation.Profile = profileName
	data.Infomation.Region = regionName

	dataChan <- data
	return
}
