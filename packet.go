package main

var r *requestHeader

const VBAN_PROTOCOL_TXT = 0x40

// requestHeader represents a single request header
type requestHeader struct {
	name         string
	bpsIndex     int
	channel      int
	framecounter []byte
}

// newRequestHeader returns a pointer to a requestHeader struct as a singleton
func newRequestHeader(streamname string, bpsI, channel int) *requestHeader {
	if r != nil {
		return r
	}
	return &requestHeader{streamname, bpsI, channel, make([]byte, 4)}
}

// sr defines the samplerate for the request
func (r *requestHeader) sr() byte {
	return byte(VBAN_PROTOCOL_TXT + r.bpsIndex)
}

// nbc defines the channel of the request
func (r *requestHeader) nbc() byte {
	return byte(r.channel)
}

// streamname defines the stream name of the text request
func (r *requestHeader) streamname() []byte {
	b := make([]byte, 16)
	copy(b, r.name)
	return b
}

// header returns a fully formed text request packet header
func (t *requestHeader) header() []byte {
	h := []byte("VBAN")
	h = append(h, t.sr())
	h = append(h, byte(0))
	h = append(h, t.nbc())
	h = append(h, byte(0x10))
	h = append(h, t.streamname()...)
	h = append(h, t.framecounter...)
	return h
}
