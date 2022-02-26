package client

import (
	"log"
	"time"

	"github.com/Antriko/go-network/shared"
	"github.com/gookit/color"
)

// Global scope
var (
	DataReadChan  chan interface{}
	DataWriteChan chan shared.Serializable
)

var red = color.New(color.FgBlack, color.BgRed).Render
var yellow = color.New(color.FgBlack, color.BgYellow).Render
var green = color.New(color.FgBlack, color.BgGreen).Render
var magenta = color.New(color.FgBlack, color.BgMagenta).Render

func serverConn() {
	log.SetFlags(log.Lshortfile)
	Handler := shared.NewConnectionHandler()
	Handler.Dial("127.0.0.1", func(conn *shared.DualConnection) {
		log.Println(yellow(" Handling connections "))
		// Move channels into global scope
		DataWriteChan = conn.DataWriteChan
		DataReadChan = conn.DataReadChan

		// Initial connection packet
		DataWriteChan <- &shared.C2SUserLoginPacket{
			Username:           player.username,
			UserModelSelection: player.UserModelSelection,
		}
		log.Println(DataWriteChan)

		for {
			select {
			case r := <-DataReadChan:
				log.Printf("%T", r)
				switch typed := r.(type) {
				case *shared.A2APingPacket:
					log.Println("pinged", typed.Reliability)
					DataWriteChan <- &shared.A2APongPacket{Reliability: typed.Reliability}
				case *shared.A2APongPacket:
					log.Println("ponged", typed.Reliability)

				case *shared.S2CUserLoginResponsePacket: // Send constant coord updates to server

					go func() {
						timer := time.NewTimer(time.Second / 20)
						for {
							select {
							case <-timer.C:
								DataWriteChan <- &shared.C2SUpdatePlayerPosPacket{
									X:      player.pos.X,
									Y:      player.pos.Y,
									Z:      player.pos.Z,
									Facing: player.rotation.facing,
								}
								timer.Reset(time.Second)
							}
						}
					}()

					// Get selection of models of players that are already logged in.
					for _, value := range typed.UsersModelSelection {
						connectedPlayers[value.Username] = shared.OtherPlayer{
							Username:           value.Username,
							UserModelSelection: value.UserModelSelection,
						}
					}

					// Done with setup
					DataWriteChan <- &shared.C2SUserFinishedSetupPacket{}
				case *shared.S2CUserInformationPacket:
					log.Println(typed.User.UserModelSelection)
					connectedPlayers[typed.User.Username] = shared.OtherPlayer{
						Username:           typed.User.Username,
						UserModelSelection: typed.User.UserModelSelection,
					}
				case *shared.S2CUpdatePlayersPosPacket:
					players = typed.Coords
				case *shared.S2CChatMessagePacket:
					chatHistory = append(chatHistory, ChatMessage{
						typed.Username,
						typed.Type,
						typed.Time,
						typed.Message,
					})
				}

			}
		}
	})
}
