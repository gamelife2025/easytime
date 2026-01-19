package commands

import (
	"fmt"
	"strings"

	"github.com/gamelife2025/easytime/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "转换时间",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := utils.Get(strings.Join(args, " "))
		if err != nil {
			return err
		}
		PARSE_TIME = t
		ShiftRun(cmd)
		fmt.Println(utils.TimeStd(PARSE_TIME))
		return nil
	},
}
