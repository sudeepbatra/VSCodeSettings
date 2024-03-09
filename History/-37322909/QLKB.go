package socketio

import (
	"errors"
	"sync"
	"time"
)

// Socket represents a Socket.IO client.
type Socket struct {
	id            string
	nsp           string
	connected     bool
	io            *Manager
	auth          map[string]string
	acks          map[int]*ackCallback
	ackMu         sync.Mutex
	sendBuffer    []*packet
	receiveBuffer []*packet
	onAnyIncoming []ListenerFunc
	onAnyOutgoing []ListenerFunc
	closeErrChan  chan struct{}
	closeOnce     sync.Once
}

// packet represents a Socket.IO packet.
type packet struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ackCallback represents an acknowledgment callback.
type ackCallback struct {
	id         int
	timer      *time.Timer
	callback   AckFunc
	expireOnce sync.Once
}

// NewSocket creates a new Socket.IO client socket.
func newSocket(io *Manager, nsp string, opts *Options) *Socket {
	socket := &Socket{
		nsp:           nsp,
		io:            io,
		auth:          make(map[string]string),
		acks:          make(map[int]*ackCallback),
		sendBuffer:    make([]*packet, 0),
		receiveBuffer: make([]*packet, 0),
		closeErrChan:  make(chan struct{}),
	}
	if opts != nil {
		socket.auth = opts.Auth
	}
	return socket
}

// Connect establishes a connection to the Socket.IO server.
func (s *Socket) Connect() error {
	return s.open()
}

// Emit sends an event to the server.
func (s *Socket) Emit(event string, args ...interface{}) error {
	return s.emit(event, args...)
}

// On registers an event handler for the given event.
func (s *Socket) On(event string, fn ListenerFunc) {
	s.io.on(event, fn)
}

// OnAnyIncoming registers a listener for any incoming events.
func (s *Socket) OnAnyIncoming(fn ListenerFunc) {
	s.onAnyIncoming = append(s.onAnyIncoming, fn)
}

// OffAnyIncoming removes all listeners for any incoming events.
func (s *Socket) OffAnyIncoming() {
	s.onAnyIncoming = nil
}

// OffAnyIncomingListener removes a specific listener from the list of any incoming event listeners.
func (s *Socket) OffAnyIncomingListener(fn ListenerFunc) {
	for i, listener := range s.onAnyIncoming {
		if listener == fn {
			s.onAnyIncoming = append(s.onAnyIncoming[:i], s.onAnyIncoming[i+1:]...)
			break
		}
	}
}

// OnAnyOutgoing registers a listener for any outgoing events.
func (s *Socket) OnAnyOutgoing(fn ListenerFunc) {
	s.onAnyOutgoing = append(s.onAnyOutgoing, fn)
}

// OffAnyOutgoing removes all listeners for any outgoing events.
func (s *Socket) OffAnyOutgoing() {
	s.onAnyOutgoing = nil
}

// OffAnyOutgoingListener removes a specific listener from the list of any outgoing event listeners.
func (s *Socket) OffAnyOutgoingListener(fn ListenerFunc) {
	for i, listener := range s.onAnyOutgoing {
		if listener == fn {
			s.onAnyOutgoing = append(s.onAnyOutgoing[:i], s.onAnyOutgoing[i+1:]...)
			break
		}
	}
}

// Close disconnects the socket.
func (s *Socket) Close() error {
	s.closeOnce.Do(func() {
		s.closeErrChan <- struct{}{}
		s.close()
	})
	return nil
}

// open establishes a connection to the Socket.IO server.
func (s *Socket) open() error {
	if s.connected || s.io.isReconnecting() {
		return nil
	}

	s.subEvents()
	s.io.open() // ensure open

	if s.io.readyState == ReadyStateOpen {
		s.onOpen()
	}

	return nil
}

// send sends a message to the server.
func (s *Socket) send(args ...interface{}) {
	s.emit(EVENTMessage, args...)
}

// emit sends an event to the server.
func (s *Socket) emit(event string, args ...interface{}) error {
	if isReservedEvent(event) {
		return errors.New("'" + event + "' is a reserved event name")
	}

	ack := extractAck(args)

	packet := &packet{
		Type: EVENT,
		Data: append([]interface{}{event}, args...),
	}

	if ack != nil {
		ackID := s.nextID()
		s.ackMu.Lock()
		s.acks[ackID] = ack
		s.ackMu.Unlock()
		packet.Type = EVENTAck
		packet.Data = append(packet.Data.([]interface{}), ackID)
		if ack.Timeout() > 0 {
			go s.handleAckTimeout(ackID, ack.Timeout())
		}
	}

	if s.connected {
		s.packet(packet)
	} else {
		s.sendBuffer = append(s.sendBuffer, packet)
	}

	return nil
}

// packet sends a packet to the server.
func (s *Socket) packet(p *packet) {
	if p.Type == EVENT {
		s.handleAnyIncoming(p.Data)
	}

	p.NSP = s.nsp
	s.io.packet(p)
}

// subEvents subscribes to socket-related events.
func (s *Socket) subEvents() {
	if s.subs != nil {
		return
	}

	s.subs = []Handle{
		s.io.on(EventOpen, s.onOpen),
		s.io.on(EventPacket, s.onPacket),
		s.io.on(EventError, s.onError),
		s.io.on(EventClose, s.onClose),
	}
}

// isActive checks if the socket is active.
func (s *Socket) isActive() bool {
	return s.subs != nil
}

// nextID returns the next available ack ID.
func (s *Socket) nextID() int {
	return s.io.nextID()
}

// onOpen is called when the transport is open.
func (s *Socket) onOpen(args ...interface{}) {
	s.connected = true
	s.emitBuffered()
	s.io.call(EventConnect)
}

// onClose is called when the transport is closed.
func (s *Socket) onClose(args ...interface{}) {
	reason := ""
	if len(args) > 0 {
		reason = args[0].(string)
	}

	s.connected = false
	s.id = ""
	s.io.call(EventDisconnect, reason)
}

// onPacket is called when a packet is received.
func (s *Socket) onPacket(args ...interface{}) {
	p, ok := args[0].(*packet)
	if !ok {
		return
	}

	if p.Type == CONNECT {
		if p.Data != nil {
			data, ok := p.Data.(map[string]interface{})
			if !ok {
				return
			}

			if sid, exists := data["sid"].(string); exists {
				s.onConnect(sid)
				return
			}
		}
		s.io.call(EventConnectError, errors.New("Invalid server response"))
		return
	}

	if p.Type == EVENT || p.Type == BINARY_EVENT {
		pData, ok := p.Data.([]interface{})
		if !ok || len(pData) == 0 {
			return
		}

		event, ok := pData[0].(string)
		if !ok {
			return
		}

		args := pData[1:]
		s.handleAnyOutgoing(args)
		s.handleAnyIncoming(pData)

		if s.connected {
			if len(args) > 0 {
				s.io.call(event, args...)
			} else {
				s.io.call(event)
			}
		} else {
			s.receiveBuffer = append(s.receiveBuffer, p)
		}
	}

	if p.Type == ACK || p.Type == BINARY_ACK {
		pData, ok := p.Data.([]interface{})
		if !ok || len(pData) == 0 {
			return
		}

		id, ok := pData[0].(float64)
		if !ok {
			return
		}

		ackData := pData[1:]
		ack := s.popAck(int(id))
		if ack != nil {
			if len(ackData) > 0 {
				ack.callback(ackData...)
			} else {
				ack.callback()
			}
		}
	}
}

// onConnect is called when the socket connects.
func (s *Socket) onConnect(id string) {
	s.connected = true
	s.id = id
	s.emitBuffered()
	s.io.call(EventConnect)
}

// emitBuffered emits buffered events and packets.
func (s *Socket) emitBuffered() {
	for _, p := range s.receiveBuffer {
		if p.Type == EVENT {
			event, ok := p.Data.([]interface{})[0].(string)
			if ok {
				args := p.Data.([]interface{})[1:]
				if len(args) > 0 {
					s.io.call(event, args...)
				} else {
					s.io.call(event)
				}
			}
		}
	}
	s.receiveBuffer = nil

	for _, p := range s.sendBuffer {
		s.packet(p)
	}
	s.sendBuffer = nil
}

// handleAnyIncoming calls any incoming event listeners.
func (s *Socket) handleAnyIncoming(args interface{}) {
	if s.onAnyIncoming != nil {
		listeners := s.onAnyIncoming
		argsArray := toArray(args)

		for _, listener := range listeners {
			listener(argsArray...)
		}
	}
}

// handleAnyOutgoing calls any outgoing event listeners.
func (s *Socket) handleAnyOutgoing(args interface{}) {
	if s.onAnyOutgoing != nil {
		listeners := s.onAnyOutgoing
		argsArray := toArray(args)

		for _, listener := range listeners {
			listener(argsArray...)
		}
	}
}

// popAck removes an acknowledgment callback from the list.
func (s *Socket) popAck(id int) *ackCallback {
	s.ackMu.Lock()
	defer s.ackMu.Unlock()

	ack, exists := s.acks[id]
	if exists {
		delete(s.acks, id)
	}
	return ack
}

// handleAckTimeout handles an acknowledgment timeout.
func (s *Socket) handleAckTimeout(id int, timeout time.Duration) {
	select {
	case <-s.closeErrChan:
		// Socket closed, do nothing
	case <-time.After(timeout):
		s.ackMu.Lock()
		defer s.ackMu.Unlock()

		if ack, exists := s.acks[id]; exists {
			ack.callback(errors.New("acknowledgment timeout"))
			delete(s.acks, id)
		}
	}
}

// close disconnects the socket.
func (s *Socket) close() {
	if s.subs != nil {
		for _, sub := range s.subs {
			sub.Destroy()
		}
		s.subs = nil
	}

	s.ackMu.Lock()
	for _, ack := range s.acks {
		if ack.timer != nil {
			ack.timer.Stop()
		}
	}
	s.ackMu.Unlock()

	s.io.disconnect(s)
}

// isReservedEvent checks if an event name is reserved.
func isReservedEvent(event string) bool {
	switch event {
	case EVENTConnect, EVENTDisconnect, EVENTConnectError:
		return true
	default:
		return false
	}
}

// toArray converts an interface to an array.
func toArray(args interface{}) []interface{} {
	if args == nil {
		return nil
	}

	if arr, ok := args.([]interface{}); ok {
		return arr
	}

	return nil
}
