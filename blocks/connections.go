package blocks

import (
	"encoding/json"
	"log"
	"time"
)

func Connection(b *Block) {
	var last BMsg
	var rate float64 // rate in messages per second of this block
	var N float64    // number of messages passed through this block
	var t time.Time
	for {
		select {
		case msg := <-b.InChan:
			last = msg
			broadcast(b.OutChans, msg)
			// rate calc
			if t.IsZero() {
				// this is the connection's first message
				t = time.Now()
				break
			}
			N++
			dt := time.Since(t).Seconds()
			rate = ((N-1.0)/N)*rate + (1.0/N)*dt
			t = time.Now()
		case query := <-b.Routes["last_message"]:
			mj, err := json.Marshal(last)
			if err != nil {
				log.Println(err.Error())
			}
			query.ResponseChan <- mj
		case query := <-b.Routes["rate"]:
			out := map[string]float64{"rate": rate}
			mj, err := json.Marshal(out)
			if err != nil {
				log.Println(err.Error())
			}
			query.ResponseChan <- mj
		case msg := <-b.AddChan:
			updateOutChans(msg, b)
		case <-b.QuitChan:
			quit(b)
			return
		}
	}

}