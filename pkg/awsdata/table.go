package awsdata

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func tableData(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
