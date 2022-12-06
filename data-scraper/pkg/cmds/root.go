package cmds

import "github.com/spf13/cobra"

func NewRootCommand() (*cobra.Command, error) {
	root := cobra.Command{
		Use:   "data-scraper",
		Short: "data-scraper: dedicated tool of kaniuse for scraping data from kubernetes API",
	}
	apiLifecycleCmd, err := NewApiLifecycleCommand()
	if err != nil {
		return nil, err
	}
	root.AddCommand(apiLifecycleCmd)
	kindsCmd, err := NewKindsCmd()
	if err != nil {
		return nil, err
	}
	root.AddCommand(kindsCmd)
	return &root, nil
}
