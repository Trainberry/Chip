package setup

import (
	"errors"
	"net/netip"
	"time"

	"github.com/soypat/cyw43439"

	"github.com/soypat/seqs/stacks"
)

const mtu = cyw43439.MTU

type WiFiSetupConfig struct {
	Hostname string
	RequestedIP string
	UDPPorts uint16
	TCPPorts uint16
	WiFiSSID string
	WiFiPassword string
}

func setupWiFi(piDevice *cyw43439.Device, cfg WiFiSetupConfig) (*stacks.PortStack, error) {
	cfg.UDPPorts++ // Add extra UDP port for DHCP client.

	var err error
	var reqAddr netip.Addr

	reqAddr, err = netip.ParseAddr(cfg.RequestedIP)
	if err != nil {
		return nil, err
	}

	if err = piDevice.JoinWPA2(cfg.WiFiSSID, cfg.WiFiPassword); err != nil {
		return nil, errors.New("wifi join failed: " + err.Error())
	}

	mac, _ := piDevice.HardwareAddr6()

	stack := stacks.NewPortStack(stacks.PortStackConfig{
		MAC:             mac,
		MaxOpenPortsUDP: int(cfg.UDPPorts),
		MaxOpenPortsTCP: int(cfg.TCPPorts),
		MTU:             mtu,
		Logger: Logger,
	})


	piDevice.RecvEthHandle(stack.RecvEth)

	// Begin asynchronous packet handling.
	go nicLoop(piDevice, stack)

	stack.SetAddr(reqAddr)
	return stack, nil
}

func nicLoop(dev *cyw43439.Device, Stack *stacks.PortStack) {
	// Maximum number of packets to queue before sending them.
	const (
		queueSize                = 3
		maxRetriesBeforeDropping = 3
	)
	var queue [queueSize][mtu]byte
	var lenBuf [queueSize]int
	var retries [queueSize]int
	markSent := func(i int) {
		queue[i] = [mtu]byte{} // Not really necessary.
		lenBuf[i] = 0
		retries[i] = 0
	}
	for {
		stallRx := true
		// Poll for incoming packets.
		for i := 0; i < 1; i++ {
			gotPacket, err := dev.PollOne()
			if err != nil {
				Logger.Error("poll error: " + err.Error())
			}
			if !gotPacket {
				break
			}
			stallRx = false
		}

		// Queue packets to be sent.
		for i := range queue {
			if retries[i] != 0 {
				continue // Packet currently queued for retransmission.
			}
			var err error
			buf := queue[i][:]
			lenBuf[i], err = Stack.HandleEth(buf[:])
			if err != nil {
				Logger.Error("stack error n(should be 0)=", lenBuf[i], "err=", err.Error())
				lenBuf[i] = 0
				continue
			}
			if lenBuf[i] == 0 {
				break
			}
		}
		stallTx := lenBuf == [queueSize]int{}
		if stallTx {
			if stallRx {
				// Avoid busy waiting when both Rx and Tx stall.
				time.Sleep(51 * time.Millisecond)
			}
			continue
		}

		// Send queued packets.
		for i := range queue {
			n := lenBuf[i]
			if n <= 0 {
				continue
			}
			err := dev.SendEth(queue[i][:n])
			if err != nil {
				// Queue packet for retransmission.
				retries[i]++
				if retries[i] > maxRetriesBeforeDropping {
					markSent(i)
					Logger.Error("dropped outgoing packet:", err.Error())
				}
			} else {
				markSent(i)
			}
		}
	}
}
