package cmd

import (
	"fmt"
	"io"

	"github.com/okaryo/dlog/internal/model"
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
	cmd.AddCommand(newLogCmd(svc, out), newMarkdownCmd(svc, out), newAmendCmd(svc))

	return cmd
}

func newLogCmd(svc *service.Service, out io.Writer) *cobra.Command {
	var date string

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show logs for today or a specified date",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				dayLog *model.DayLog
				err    error
			)

			if date == "" {
				dayLog, err = svc.GetTodayLog()
			} else {
				dayLog, err = svc.GetLogByDate(date)
			}
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

	cmd.Flags().StringVarP(&date, "date", "d", "", "Show logs for the specified date (YYYY-MM-DD)")

	return cmd
}

func newMarkdownCmd(svc *service.Service, out io.Writer) *cobra.Command {
	var date string

	cmd := &cobra.Command{
		Use:   "md",
		Short: "Output logs as Markdown for today or a specified date",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				dayLog *model.DayLog
				err    error
			)

			if date == "" {
				dayLog, err = svc.GetTodayLog()
			} else {
				dayLog, err = svc.GetLogByDate(date)
			}
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

	cmd.Flags().StringVarP(&date, "date", "d", "", "Output logs for the specified date (YYYY-MM-DD)")

	return cmd
}

func newAmendCmd(svc *service.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "amend [text]",
		Short: "Replace today's most recent log entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.AmendTodayLog(args[0])
		},
	}

	return cmd
}
