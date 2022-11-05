package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"

	log "github.com/sirupsen/logrus"
)

var (
	host       string
	port       int
	streamname string
	bps        int
	channel    int
	delay      int
	loglevel   int

	bpsOpts = []int{0, 110, 150, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 31250,
		38400, 57600, 115200, 128000, 230400, 250000, 256000, 460800, 921600,
		1000000, 1500000, 2000000, 3000000}
)

type (
	// connection represents the configurable fields of a config.toml
	connection struct {
		Host       string
		Port       int
		Streamname string
	}

	// config maps toml headers
	config struct {
		Connection map[string]connection
	}
)

func setLogLevel(loglevel int) {
	switch loglevel {
	case 0:
		log.SetLevel(log.WarnLevel)
	case 1:
		log.SetLevel(log.InfoLevel)
	case 2:
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	flag.StringVar(&host, "host", "", "vban host")
	flag.StringVar(&host, "h", "", "vban host (shorthand)")
	flag.IntVar(&port, "port", 6980, "vban server port")
	flag.IntVar(&port, "p", 6980, "vban server port (shorthand)")
	flag.StringVar(&streamname, "streamname", "Command1", "stream name for text requests")
	flag.StringVar(&streamname, "s", "Command1", "stream name for text requests (shorthand)")
	flag.IntVar(&bps, "bps", 0, "vban bps")
	flag.IntVar(&bps, "b", 0, "vban bps (shorthand)")
	flag.IntVar(&channel, "channel", 0, "vban channel")
	flag.IntVar(&channel, "c", 0, "vban channel (shorthand)")
	flag.IntVar(&delay, "delay", 20, "delay between requests")
	flag.IntVar(&delay, "d", 20, "delay between requests (shorthand)")
	flag.IntVar(&loglevel, "loglevel", 0, "log level")
	flag.IntVar(&loglevel, "l", 0, "log level (shorthand)")
	flag.Parse()

	setLogLevel(loglevel)

	c, err := vbanConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	header := newRequestHeader(streamname, indexOf(bpsOpts, bps), channel)
	for _, arg := range flag.Args() {
		err := send(c, header, arg)
		if err != nil {
			log.Error(err)
		}
	}
}

// vbanConnect establishes a VBAN connection to remote host
func vbanConnect() (*net.UDPConn, error) {
	if host == "" {
		conn, err := connFromToml()
		if err != nil {
			return nil, err
		}
		host = conn.Host
		port = conn.Port
		streamname = conn.Streamname
		if host == "" {
			err := errors.New("must provide a host with --host flag or config.toml")
			return nil, err
		}
	}
	CONNECT := fmt.Sprintf("%s:%d", host, port)

	s, _ := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to ", c.RemoteAddr().String())

	return c, nil
}

// connFromToml parses connection info from config.toml
func connFromToml() (*connection, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	f := filepath.Join(homeDir, ".vbantxt_cli", "config.toml")
	if _, err := os.Stat(f); err != nil {
		err := fmt.Errorf("unable to locate %s", f)
		return nil, err
	}

	var c config
	_, err = toml.DecodeFile(f, &c.Connection)
	if err != nil {
		return nil, err
	}
	conn := c.Connection["connection"]
	return &conn, nil
}

// send sends a VBAN text request over UDP to remote host
func send(c *net.UDPConn, h *requestHeader, msg string) error {
	log.Debug("Sending '", msg, "' to: ", c.RemoteAddr().String())
	data := []byte(msg)
	_, err := c.Write(append(h.header(), data...))
	if err != nil {
		return err
	}
	var a uint32
	_ = binary.Read(bytes.NewReader(h.framecounter[:]), binary.LittleEndian, &a)
	binary.LittleEndian.PutUint32(h.framecounter[:], a+1)

	time.Sleep(time.Duration(delay) * time.Millisecond)

	return nil
}
