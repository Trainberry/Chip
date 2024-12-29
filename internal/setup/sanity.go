package setup

import(
	"time"
	"chip/internal/configuration"
	"errors"
	"net/netip"
	"github.com/soypat/seqs"
	"github.com/soypat/seqs/stacks"
	"machine"
	"math/rand"
)

func sanity() {
	waitTime := 4
	failCount := 0

	for {
		time.Sleep(time.Second * time.Duration(waitTime))

		if err := checkAlive(); err == nil {
			machine.Watchdog.Update()
			Logger.Info("Sanity check passed")
			waitTime = 4
		} else {
			Logger.Error("Sanity check failed! Error was: " + err.Error())
			if failCount < 5 {
				machine.Watchdog.Update()
			}
			failCount = failCount + 1
			waitTime = 0
		}
	}
}

func checkAlive() error {
	start := time.Now()

	var err error
	var conn *stacks.TCPConn
	var svAddr netip.AddrPort
	
	// Parse address from configuration
	for svAddr, err = netip.ParseAddrPort(configuration.ServerAddr) ; err != nil ; {
		return errors.New("0: " + err.Error())
	}

	// Set-up port and create connection
	rng := rand.New(rand.NewSource(int64(time.Now().Sub(start))))
	clientAddr := netip.AddrPortFrom(PortStack.Addr(), uint16(rng.Intn(65535-16636)+16636))

	for conn, err = stacks.NewTCPConn(PortStack, stacks.TCPConnConfig{}) ; err != nil ; {
		return errors.New("1: " + err.Error())
	}

	// Make sure to timeout the connection if it takes too long.
	conn.SetDeadline(time.Now().Add(time.Second * 1))

	// Open connection
	if err = conn.OpenDialTCP(clientAddr.Port(), configuration.ServerMAC, svAddr, seqs.Value(uint16(rng.Intn(65535-16636)+16636))) ; err != nil {
		return errors.New("2: " + err.Error())
	}

	count := 0
	for conn.State() != seqs.StateEstablished && count < 10 { 
		time.Sleep(100 * time.Millisecond)
		count += 1
	}

	if conn.State() != seqs.StateEstablished {
		return errors.New("nope")
	}

	conn.Close()

	return nil
}