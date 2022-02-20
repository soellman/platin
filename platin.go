package platin

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	portControl  = 50006
	portMedia    = 7777
	dialTimeout  = 5 * time.Second
	readDeadline = time.Second
)

var (
	// control port 50006
	cmdSources   = []byte{0x01}
	cmdSetSource = []byte{0x03}
	cmdPower     = []byte{0x0d}
	cmdSetPower  = []byte{0x0f}

	// media port 7777
	cmdRegister  = []byte{0x02, 0x03}
	cmdVolume    = []byte{0x01, 0x40}
	cmdSetVolume = []byte{0x02, 0x40}

	// ErrSourceNotFound is returned when attempting to set a non-existent source
	ErrSourceNotFound = errors.New("source not found")
)

type sources struct {
	Sources []Source `json:"src"`
}

// Source represents a source on a Platin hub
type Source struct {
	Index  int     `json:"ix"`
	Name   string  `json:"name"`
	Active boolean `json:"sts"`
}

type boolean bool

func (b *boolean) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "0" {
		*b = false
	} else if str == "1" {
		*b = true
	} else {
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", str)
	}
	return nil
}

// Hub controls a Platin hub
type Hub struct {
	host string
}

// NewHub creates a new platin hub
func NewHub(host string) *Hub {
	return &Hub{host: host}
}

// Power returns the current power state
func (h *Hub) Power() (bool, error) {
	msg := newControlMessage(cmdPower)
	resp := controlResponse{}
	err := h.connectAndSendReceive(portControl, msg, &resp)
	if err != nil {
		return false, err
	}
	return resp.payload[0] == 0x01, nil
}

// SetPower sets the power to the speakers
func (h *Hub) SetPower(on bool) error {
	current, err := h.Power()
	if err != nil {
		return err
	}

	if current == on {
		return nil
	}

	onByte := 0x00
	if on {
		onByte = 0x01
	}
	msg := newControlMessage(append(cmdSetPower, byte(onByte)))
	return h.connectAndSend(portControl, msg)
}

// TogglePower will toggle the power to the speakers
func (h *Hub) TogglePower() error {
	msg := newControlMessage(append(cmdSetPower, 0x00))
	return h.connectAndSend(portControl, msg)
}

// ActiveSource returns the source that is currently selected
func (h *Hub) ActiveSource() (Source, error) {
	source := Source{}
	src, err := h.Sources()
	if err != nil {
		return source, err
	}

	for _, s := range src {
		if s.Active == true {
			return s, nil
		}
	}

	return source, ErrSourceNotFound
}

// Sources returns a list of sources on the hub
func (h *Hub) Sources() ([]Source, error) {
	s := sources{}
	msg := newControlMessage(cmdSources)
	resp := controlResponse{}
	err := h.connectAndSendReceive(portControl, msg, &resp)
	if err != nil {
		return s.Sources, err
	}

	err = json.Unmarshal(resp.payload, &s)
	return s.Sources, err
}

// SetSource sets the active source to the provided name
// Returns an error if source does not exist
func (h *Hub) SetSource(name string) error {
	sources, err := h.Sources()
	if err != nil {
		return err
	}

	var index int
	found := false
	for _, source := range sources {
		if name == source.Name {
			index = source.Index
			found = true
			break
		}
	}
	if !found {
		return ErrSourceNotFound
	}

	cmd := append(cmdSetSource, byte(index))
	msg := newControlMessage(cmd)
	err = h.connectAndSend(portControl, msg)
	return err
}

// Volume returns the current volume level
func (h *Hub) Volume() (int, error) {
	vol := 0
	err := h.connectAndRun(portMedia, func(conn net.Conn) error {
		if err := mediaRegister(conn); err != nil {
			return err
		}

		msg := newMediaMessage(cmdVolume, "")
		resp := mediaResponse{}
		err := sendAndReceive(conn, msg, &resp)
		if err != nil {
			return err
		}

		vol, err = strconv.Atoi(resp.Payload())
		return err
	})

	return vol, err
}

// SetVolume sets the volume level
// Valid levels are 0 - 100
func (h *Hub) SetVolume(vol int) error {
	return h.connectAndRun(portMedia, func(conn net.Conn) error {
		if err := mediaRegister(conn); err != nil {
			return err
		}

		msg := newMediaMessage(cmdSetVolume, strconv.Itoa(vol))
		resp := mediaResponse{}
		return sendAndReceive(conn, msg, &resp)
	})
}

func (h *Hub) connectAndSend(port int, msg message) error {
	host := fmt.Sprintf("%s:%d", h.host, port)
	conn, err := net.DialTimeout("tcp", host, dialTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(readDeadline))

	_, err = conn.Write(msg.Bytes())
	return err
}

func (h *Hub) connectAndSendReceive(port int, msg message, resp response) error {
	host := fmt.Sprintf("%s:%d", h.host, port)
	conn, err := net.DialTimeout("tcp", host, dialTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(readDeadline))

	_, err = conn.Write(msg.Bytes())
	return resp.Read(conn)
}

func (h *Hub) connectAndRun(port int, f func(conn net.Conn) error) error {
	host := fmt.Sprintf("%s:%d", h.host, port)
	conn, err := net.DialTimeout("tcp", host, dialTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return f(conn)
}

func sendAndReceive(conn net.Conn, msg message, resp response) error {
	conn.SetReadDeadline(time.Now().Add(readDeadline))
	_, err := conn.Write(msg.Bytes())
	if err != nil {
		return err
	}
	return resp.Read(conn)
}

// Media endpoint wants new connections to "register"
func mediaRegister(conn net.Conn) error {
	addr := conn.LocalAddr().(*net.TCPAddr)
	msg := newMediaMessage(cmdRegister, fmt.Sprint(addr.IP))
	resp := mediaResponse{}
	return sendAndReceive(conn, msg, &resp)
}

type message interface {
	Bytes() []byte
}

type controlMessage struct {
	m []byte
}

func newControlMessage(m []byte) *controlMessage {
	return &controlMessage{m: m}
}

func (m *controlMessage) Bytes() []byte {
	// format: 00 [msg length (1B)] 03 02 + command payload
	// length is len(payload) + 2
	b := append([]byte{
		0x00,
		byte(len(m.m) + 2),
		0x03,
		0x02,
	}, m.m...)
	return b
}

type response interface {
	Read(net.Conn) error
}

type controlResponse struct {
	length  []byte // length, two bytes
	rtype   []byte // response type, three bytes (0x00 0x02 cmd)
	payload []byte // body
}

func (r controlResponse) String() string {
	// format: [length (2B)] [type (3B)] [payload]
	return fmt.Sprintf("[% x] [% x] [% x]", r.length, r.rtype, r.payload)
}

func (r *controlResponse) Read(conn net.Conn) error {
	msg := make([]byte, 5)
	_, err := conn.Read(msg)
	if err != nil {
		return err
	}

	r.length = msg[0:2]
	r.rtype = msg[2:]

	length := binary.BigEndian.Uint16(r.length) - 3
	if length > 0 {
		r.payload = make([]byte, length)
		_, err = conn.Read(r.payload)
		if err != nil {
			return err
		}
	}

	return nil
}

type mediaMessage struct {
	cmd     []byte
	payload string
}

func newMediaMessage(cmd []byte, payload string) *mediaMessage {
	return &mediaMessage{cmd: cmd, payload: payload}
}

func (m *mediaMessage) Bytes() []byte {
	// format: 00 00 [msg type (2b)] 00 00 00 00 [length (1b)] 00 + optional payload
	// Guessing that byte 7-8 together represent length. no matter, our messages are short
	b := []byte{0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00}
	b[2] = m.cmd[0]
	b[3] = m.cmd[1]
	b[8] = byte(len(m.payload))
	if len(m.payload) > 0 {
		b = append(b, []byte(m.payload)...)
	}
	return b
}

type mediaResponse struct {
	rtype   []byte // response type, three bytes
	unknown []byte // unknown content, three bytes
	length  []byte // payload length, two bytes
	payload []byte // body, in ASCII
}

func (r *mediaResponse) Read(conn net.Conn) error {
	msg := make([]byte, 10)
	_, err := conn.Read(msg)
	if err != nil {
		return err
	}
	r.rtype = msg[2:5]
	r.unknown = msg[5:8]
	r.length = msg[8:]

	length := binary.BigEndian.Uint16(r.length)
	if length > 0 {
		r.payload = make([]byte, length)
		if _, err = conn.Read(r.payload); err != nil {
			return err
		}
	}

	return nil
}

func (r mediaResponse) Payload() string {
	return string(r.payload)
}

func (r mediaResponse) String() string {
	// format: 00 00 [type (3B)] [unknown (3B)] [payload length (2B)] [payload (in ASCII)]
	str := fmt.Sprintf("00 00 [% x] [% x] [% x]", r.rtype, r.unknown, r.length)
	if len(r.payload) > 0 {
		str = fmt.Sprintf("%s %q", str, string(r.payload))
	}
	return str
}
