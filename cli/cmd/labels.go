package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/setavenger/blindbitd/pb"

	"github.com/setavenger/blindbitd/cli/lib"
)

// labelsCmd represents the chain command
var (
	labelsCmd = &cobra.Command{
		Use:   "labels",
		Short: "Operations related to labels",
		Long:  ``,
		// no Run so it goes directly to help
	}

	newLabelComment string

	labelsNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new label",
		Long:  `This command creates a new label and returns the new address`,
		Run: func(cmd *cobra.Command, args []string) {
			if newLabelComment == "" {
				log.Fatalln("comment for new label should not be empty")
			}

			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			newLabel, err := client.CreateNewLabel(ctx, &pb.NewLabelRequest{Comment: newLabelComment})
			if err != nil {
				log.Fatalf("could not create new label: %v\n", err)
			}
			fmt.Printf("New label created: %s\n", newLabel.Address)
		},
	}
)

func init() {
	RootCmd.AddCommand(labelsCmd)
	labelsCmd.AddCommand(labelsNewCmd)
	labelsCmd.AddCommand(labelsListCmd)

	labelsNewCmd.PersistentFlags().StringVar(&newLabelComment, "comment", "", "Set a comment for the new label. The comment allows you to identify the label address later on.")
	err := cobra.MarkFlagRequired(labelsNewCmd.PersistentFlags(), "comment")
	if err != nil {
		log.Fatalln(err)
	}
}
