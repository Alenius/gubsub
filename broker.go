package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type channelWorker struct {
	msgChannel chan gsMsg
	config     gsConfig
	quitCh     chan bool
	closed     bool
}

func (w *channelWorker) Start(conn net.Conn) {
	w.msgChannel = make(chan gsMsg, 10) // some buffer size to avoid blocking
	go func() {
		for {
			select {
			case msg := <-w.msgChannel:
				log.Println("msg", msg)
				writeGsMsg(msg, conn)
			case <-w.quitCh:
				log.Println("Quitting worker", w.config.Id)
				quitMsg := gsMsg{}.CreateCloseMsg()
				writeGsMsg(quitMsg, conn)
				w.closed = true
				return
			}

		}

	}()
}

type threadSafeSlice struct {
	sync.Mutex
	workers []*channelWorker
}

func (slice *threadSafeSlice) Push(w *channelWorker) {
	slice.Lock()
	defer slice.Unlock()

	slice.workers = append(slice.workers, w)
}

func (slice *threadSafeSlice) SendMessage(msg gsMsg) {
	slice.Lock()
	defer slice.Unlock()

	for _, worker := range slice.workers {
		if msg.Topic == worker.config.Topic {
			worker.msgChannel <- msg
		}
	}
}

func (slice *threadSafeSlice) CloseAllWorkers() {
	for _, worker := range slice.workers {
		worker.quitCh <- true
	}

	// wait for all close msgs to be written
	for {
		allClosed := true
		for _, worker := range slice.workers {
			allClosed = worker.closed
		}

		if allClosed {
			return
		}
	}

}

func checkServerError(conn net.Conn, err error) {
	if err != nil {
		log.Println("Critical error", err.Error())
		closeConnection(conn, 1)
	}
}

func readConfig(conn net.Conn) (gsConfig, error) {
	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		checkTcpMsgError(err)

		if len(msg) > 0 {
			gs_config := gsConfig{}
			err = json.Unmarshal(msg, &gs_config)
			checkError(err)
			return gs_config, nil
		}
	}
}

func writeToLedger(conn net.Conn, msg string) error {
	path := "ledger.txt"

	fileExists, _ := FileExists(path)

	if !fileExists {
		os.Create(path)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	checkError(err)
	defer f.Close()

	bytes, err := f.Write([]byte(msg))
	checkServerError(conn, err)
	_, err = f.Write([]byte{'\n'})
	checkServerError(conn, err)
	f.Sync()

	log.Println("Wrote to file: ", bytes)
	return err
}

func startBroker() {
	log.Println("starting broker")

	consumerListener, err := net.Listen("tcp", "127.0.0.1:8080")
	checkError(err)
	producerListener, err := net.Listen("tcp", "127.0.0.1:8081")
	checkError(err)

	workerSlice := threadSafeSlice{}
	doneCh := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			if conn, err := consumerListener.Accept(); err == nil {
				log.Println("Consumer connected")
				config, err := readConfig(conn)
				checkError(err)

				worker := channelWorker{config: config, quitCh: make(chan bool), closed: false}
				worker.Start(conn)
				workerSlice.Push(&worker)
			}
		}
	}()

	go func() {
		for {
			if conn, err := producerListener.Accept(); err == nil {
				go handleProducerConnection(conn, &workerSlice)
			}
		}
	}()

	go func() {
		<-sigs
		workerSlice.CloseAllWorkers()
		doneCh <- true
	}()

	<-doneCh
}

func handleProducerConnection(conn net.Conn, workerSlice *threadSafeSlice) {
	log.Println("Producer connected")
	for {
		gsMsg, err := readGsMsg(conn)
		checkError(err)

		writeToLedger(conn, gsMsg.Stringify())
		workerSlice.SendMessage(gsMsg)
	}
}

func handleConsumerConnection(conn net.Conn, worker *channelWorker) {
	defer closeConnection(conn, 0)

	config, err := readConfig(conn)
	checkError(err)

	worker.config = config

	for {
		// msg := <-channel
		// log.Println("channel id", config.Id)
		// log.Println("msg", msg)
		// writeGsMsg(msg, conn)
	}
}

func closeConnection(conn net.Conn, exitCode int) {
	log.Println("Closing connection")
	close_msg := gsMsg{}.CreateCloseMsg()
	writeGsMsg(close_msg, conn)
	conn.Close()
	syscall.Exit(exitCode)
}
