package cmd

import (
	"fmt"
	"io"

	"github.com/okaryo/dlog/internal/render"
	"github.com/okaryo/dlog/internal/service"
	"github.com/spf13/cobra"
)

func NewRootCmd(svc *service.Service, out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dlog [text]",
		Short:         "Record and view daily work logs",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return svc.AddTodayLog(args[0])
			}

			dayLog, err := svc.GetTodayLog()
			if err != nil {
				return err
			}

			output, err := render.Text(dayLog)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(out, output)
			return err
		},
	}

	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.AddCommand(newLogCmd(svc, out), newMarkdownCmd(svc, out))

	return cmd
}

func newLogCmd(svc *service.Service, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "Show today's logs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dayLog, err := svc.GetTodayLog()
			if err != nil {
				return err
			}

			output, err := render.Text(dayLog)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(out, output)
			return err
		},
	}
}

func newMarkdownCmd(svc *service.Service, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "md",
		Short: "Output today's logs as Markdown",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dayLog, err := svc.GetTodayLog()
			if err != nil {
				return err
			}

			output, err := render.Markdown(dayLog)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(out, output)
			return err
		},
	}
}
