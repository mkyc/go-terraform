package terra

import (
	"errors"
	"regexp"
	"strconv"
)

// ResourceCount represents counts of resources affected by terraform apply/plan/destroy command.
type ResourceCount struct {
	Add     int
	Change  int
	Destroy int
}

// Regular expressions for terraform commands stdout pattern matching.
const (
	applyRegexp             = `Apply complete! Resources: (\d+) added, (\d+) changed, (\d+) destroyed\.`
	destroyRegexp           = `Destroy complete! Resources: (\d+) destroyed\.`
	planWithChangesRegexp   = `(\033\[1m)?Plan:(\033\[0m)? (\d+) to add, (\d+) to change, (\d+) to destroy\.`
	planWithNoChangesRegexp = `No changes\. Infrastructure is up-to-date\.`
)

// GetResourceCount parses stdout/stderr of apply/plan/destroy commands and returns number of affected resources.
func Count(cmdout string) (*ResourceCount, error) {
	cnt := ResourceCount{}

	terraformCommandPatterns := []struct {
		regexpStr       string
		addPosition     int
		changePosition  int
		destroyPosition int
	}{
		{applyRegexp, 1, 2, 3},
		{destroyRegexp, -1, -1, 1},
		{planWithChangesRegexp, 3, 4, 5},
		{planWithNoChangesRegexp, -1, -1, -1},
	}

	for _, tc := range terraformCommandPatterns {
		pattern, err := regexp.Compile(tc.regexpStr)
		if err != nil {
			return nil, err
		}

		matches := pattern.FindStringSubmatch(cmdout)
		if matches != nil {
			if tc.addPosition != -1 {
				cnt.Add, err = strconv.Atoi(matches[tc.addPosition])
				if err != nil {
					return nil, err
				}
			}

			if tc.changePosition != -1 {
				cnt.Change, err = strconv.Atoi(matches[tc.changePosition])
				if err != nil {
					return nil, err
				}
			}

			if tc.destroyPosition != -1 {
				cnt.Destroy, err = strconv.Atoi(matches[tc.destroyPosition])
				if err != nil {
					return nil, err
				}
			}

			return &cnt, nil
		}
	}

	return nil, errors.New("can't parse Terraform output")
}
