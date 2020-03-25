package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goburrow/serial"
)

var (
	address     string
	showVersion bool
	debugMode   bool
	help        bool
)

func init() {
	flag.StringVar(&address, "address", "", "Serialport address.")
	flag.BoolVar(&showVersion, "v", false, "Show version.")
	flag.BoolVar(&debugMode, "dev", false, "Show debug msg.")
	flag.BoolVar(&help, "h", false, "Help")
}

func main() {
	// flag parse
	flag.Parse()
	flagPrase()

	// address := "/dev/ttyUSB0"
	if address == "" {
		reader := bufio.NewScanner(os.Stdin)
		for {
			fmt.Println("Enter serial port address: ")
			reader.Scan()
			address = reader.Text()
			if address != "" {
				break
			}
		}
	}
	log.Println("Now start communication...", address)

	port, err := serial.Open(
		&serial.Config{
			Address:  address,
			BaudRate: 115200,
			// DataBits: 8,
			StopBits: 1,
			Timeout:  3 * time.Second,
		})
	if err != nil {
		log.Println("Can not open", address, ":", err)
		os.Exit(0)
	}
	defer port.Close()

	for {
		fmt.Println(`Enter the action or enter h for help...`)
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		action := reader.Text()
		switch action {
		case "w":
			fmt.Println("Enter the write string(hex), like AABBCCDD")
			writeScan := bufio.NewScanner(os.Stdin)
			writeScan.Scan()
			writeString := []byte(writeScan.Text())
			if debugMode {
				log.Println(writeString)
			}
			wBufHex := make([]byte, hex.DecodedLen(len(writeString)))
			hex.Decode(wBufHex, writeString)
			_, err = port.Write(wBufHex) //寫資料出去
			if err != nil {
				log.Println("Write fail:", err)
				log.Println("Exit...")
				os.Exit(0)
			}

		case "wr":
			fmt.Println("Enter the write string(hex), like AABBCCDD")
			writeScan := bufio.NewScanner(os.Stdin)
			writeScan.Scan()
			writeString := []byte(writeScan.Text())
			if debugMode {
				log.Println(writeString)
			}
			wBufHex := make([]byte, hex.DecodedLen(len(writeString)))
			hex.Decode(wBufHex, writeString)
			_, err = port.Write(wBufHex) //寫資料出去
			if err != nil {
				log.Println("Write fail:", err)
				log.Println("Exit...")
				os.Exit(0)
			}

			data := make([]byte, 128)
			n, err := port.Read(data) //讀資料回來
			if err != nil {
				log.Println("Read fail:", err)
				log.Println("Exit...")
				os.Exit(0)
			}
			fmt.Println("Origin data: ", data[:n])
			fmt.Println("ASCII data: ", string(data[:n]))
		case "h":
			fmt.Println(`
w to write.
r to read.
wr to write and read.
exit to exit.
			`)
		case "exit":
			os.Exit(0)
		}
	}
}

func flagPrase() {
	if showVersion {
		log.Println(`
Version 1.0.1
Author: Alfred Wu.
WebSite: www.jicommand.com
		`)
		os.Exit(0)
	}

	if help {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}
}
