package server

import(
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	// "github.com/gorilla/websocket"
)

var AllRooms RoomMap

// create room and return room id
func CreateRoomHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "* ")
	roomID := AllRooms.CreateRoom()

	type resp struct{
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)

	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct{
	Message map[string]interface{}
	RoomID string
	Client *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster(){

	for{
		msg := <- broadcast

		for _, client := range AllRooms.Map[msg.RoomID]{
			if(client.Conn != msg.Client){
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil{
					log.Fatal(err)
					client.Conn.Close()
				}
			}
		}
	}
}



// join room request handler will join the client in a particular room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request){

	roomID, ok := r.URL.Query()["roomID"]

	if !ok{
		log.Println("roomID misssing in URL Parameters")
		return 
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Fatal("Web socket upgrade error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	for{
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)

		if err != nil{
			log.Fatal("Read Error: ", err)
		}

		msg.Client = ws
		msg.RoomID = roomID[0]
		log.Println(msg.Message)

		broadcast <- msg
	}
}