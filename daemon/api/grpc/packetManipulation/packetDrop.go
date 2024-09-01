package packetmanipulation

import (
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
)

func DropPackets(
	conn net.Conn,
	target net.Conn,
	wg *sync.WaitGroup,
	dropRate float64,
	corruptRate float64,
) {
	defer conn.Close()
	defer target.Close()
	defer wg.Done()

	go func() {
		_, err := io.Copy(target, conn)
		if err != nil && err != io.EOF {
			log.Println("Error copying data from client to server:", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := target.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from target:", err)
			}
			return
		}

		if rand.Float64() < dropRate {
			if rand.Float64() < 0.5 {
				log.Printf("Packet dropped: %s\n", string(buf[:n]))
				return
			}

			isCritical := rand.Float64() < corruptRate

			if isCritical {
				corruptIndex := rand.Intn(
					100,
				) // Assuming headers are within the first 100 bytes
				if corruptIndex < n {
					buf[corruptIndex] = '^' // Arbitrary critical corruption
					log.Printf(
						"Critical corruption at byte %d: %s\n",
						corruptIndex,
						string(buf[:n]),
					)
				}
			} else {
				corruptIndex := rand.Intn(n)
				buf[corruptIndex] = '#' // Arbitrary non-critical corruption
				log.Printf("Non-Critical corruption at byte %d: %s\n", corruptIndex, string(buf[:n]))
			}
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Println("Error writing to client:", err)
			return
		}

		log.Printf(
			"Server response: %s\n----------------------------------------",
			string(buf[:n]),
		)
	}
}
