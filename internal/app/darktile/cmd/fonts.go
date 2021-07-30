package cmd

import (
	"fmt"

	"github.com/liamg/fontinfo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listFontsCmd)
}

var listFontsCmd = &cobra.Command{
	Use:          "list-fonts",
	Short:        "List fonts on your system which are compatible with darktile",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {

		fonts, err := fontinfo.Match(fontinfo.MatchStyle("Regular"))
		if err != nil {
			return err
		}

		for _, font := range fonts {
			fmt.Println(font.Family)
		}
		return nil
	},
}
