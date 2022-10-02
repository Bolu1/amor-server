package server

import(
	"log"
	"sync"
	"time"
	"math/rand"
	"github.com/gorilla/websocket"
)

type Participant struct{
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct{
	Mutex sync.RWMutex
	Map map[string][]Participant
}

// initializa hashmap
func (r *RoomMap) Init(){
	r.Map = make(map[string][]Participant)
}

func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

// createroom generate id and isert into hashm ap
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()


	rand.Seed(time.Now().UnixNano())
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)

	for i := range b{
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomID := string(b)
	r.Map[roomID] = []Participant{}

	return roomID
}


func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn){
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	
	p := Participant{host, conn}

	log.Println("Inserting into room with roomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)

}

func (r *RoomMap) DeleteRoom(roomID string){
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}