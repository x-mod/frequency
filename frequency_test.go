package frequency

import (
	"log"
	"testing"
	"time"
)

func TestFrequency_Reserve(t *testing.T) {
	freq := New(
		Second(2),
		Limit(time.Second*10, 5, 5),
	)
	for i := 0; i < 11; i++ {
		if d, ok := freq.ReserveN(1); ok {
			<-time.After(d)
			log.Println("reserve ...", i)
		}
	}

	for i := 0; i < 11; i++ {
		freq.WaitN(1)
		log.Println("wait ...", i)
	}

}
