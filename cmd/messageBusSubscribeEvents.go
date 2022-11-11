/*
Copyright © 2022 IceWhaleTech

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

// messageBusSubscribeEventsCmd represents the messageBusSubscribeEvents command
var messageBusSubscribeEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "subscribe to events in message bus",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		sourceID, err := cmd.Flags().GetString(FlagMessageBusSourceID)
		if err != nil {
			log.Fatalln(err.Error())
		}

		eventNames, err := cmd.Flags().GetString(FlagMessageBusEventNames)
		if err != nil {
			log.Fatalln(err.Error())
		}

		var wsURL string

		if eventNames == "" {
			wsURL = fmt.Sprintf("ws://%s/%s/event/%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, sourceID)
		} else {
			wsURL = fmt.Sprintf("ws://%s/%s/event/%s?names=%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, sourceID, eventNames)
		}

		bufferSize, err := cmd.Flags().GetUint(FlagMessageBusMessageBufferSize)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ws, err := websocket.Dial(wsURL, "", "http://localhost")
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer ws.Close()

		log.Println("subscribed to", wsURL)

		for {
			msg := make([]byte, bufferSize)
			n, err := ws.Read(msg)
			if err != nil {
				log.Fatalln(err.Error())
			}
			log.Println(string(msg[:n]))
		}
	},
}

func init() {
	messageBusSubscribeCmd.AddCommand(messageBusSubscribeEventsCmd)

	messageBusSubscribeEventsCmd.Flags().UintP(FlagMessageBusMessageBufferSize, "m", 1024, "message buffer size")
	messageBusSubscribeEventsCmd.Flags().StringP(FlagMessageBusSourceID, "s", "", "source id")
	messageBusSubscribeEventsCmd.Flags().StringP(FlagMessageBusEventNames, "n", "", "event names (seperated by comma)")

	if err := messageBusSubscribeEventsCmd.MarkFlagRequired(FlagMessageBusSourceID); err != nil {
		log.Fatalln(err.Error())
	}
}