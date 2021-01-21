package pkg

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/secmask/go-redisproto"
)

const keySeparator = ":"

// Server accepting connections, running RESP with custom commands
type Server struct {
	ps   PeerStore
	node *BufferedNode

	ln net.Listener

	ctx context.Context
}

// NewServer creates a new server. This binds the given port, but does not
// yet accept incomming connections
func NewServer(ctx context.Context, port uint16, ps PeerStore, node *BufferedNode) (*Server, error) {
	s := &Server{
		ps:   ps,
		ctx:  ctx,
		node: node,
	}

	// create a default listenerconfig so we can pass the context
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tcp listener")
	}

	// TODO: tls stuff if needed

	s.ln = listener

	return s, nil
}

// Run the server until the context is done, or an error is encountered while
// accepting the connection. A new goroutine is spawned per new connection.
func (s *Server) Run() error {
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:

		}

		con, err := s.ln.Accept()
		if err != nil {
			// context is done
			if s.ctx.Err() != nil {
				return nil
			}
			// TODO: should we log and continue here?
			return errors.Wrap(err, "could not accept connection")
		}

		go func() {
			defer con.Close()
			s.handleCon(con)
		}()
	}
}

// Close the server and its connections
func (s *Server) Close() error {
	return errors.Wrap(s.ln.Close(), "failed to close listener")
}

func (s *Server) handleCon(conn net.Conn) {
	parser := redisproto.NewParser(conn)
	// writer := redisproto.NewWriter(bufio.NewWriter(conn))
	writer := redisproto.NewWriter(conn)

	var c connection = newUnauthenticatedConn(s)

	for {
		// Don't use the `Commands()` channel here, as that exits on any error,
		// including protocol errors
		command, err := parser.ReadCommand()
		if err != nil {
			if errors.Is(err, &redisproto.ProtocolError{}) {
				writer.WriteError(err.Error())
				continue
			}
			if errors.Is(err, io.EOF) {
				log.Debug().Msg("client closed connection")
				return
			}
			log.Error().Err(err).Msg("failed to read command")
			return
		}

		cmd := strings.ToUpper(string(command.Get(0)))
		switch cmd {
		case "PING":
			log.Debug().Msg("client PING command")
			err = writer.WriteSimpleString("PONG")
		case "HELLO":
			log.Debug().Msg("client HELLO command")
			if command.ArgCount() != 1 {
				err = writer.WriteError(errInvalidArgCount.Error())
				break
			}
			err = writer.WriteObjectsSlice(s.helloInfo())
		case "AUTH":
			log.Debug().Msg("client AUTH command")
			if command.ArgCount() != 3 {
				err = writer.WriteError(errInvalidArgCount.Error())
				break
			}

			// first arg is the dtid
			var dtid uint64 // declare here so we don't shadow err below
			dtid, err = strconv.ParseUint(string(command.Get(1)), 10, 64)
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			// second arg is the signature.
			rawSig := command.Get(2)

			if err = c.Auth(dtid, rawSig); err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			// upgrade connection
			c = newAuthenticatedConn(dtid, s)
			err = writer.WriteSimpleString("Authenticated")

		case "LPUSH":
			log.Debug().Msg("client LPUSH command")
			if command.ArgCount() != 3 {
				err = writer.WriteError(errInvalidArgCount.Error())
				break
			}

			var dtid uint64
			var subject string
			dtid, subject, err = parseKey(string(command.Get(1)))
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			if err = c.LPush(dtid, subject, command.Get(2)); err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			err = writer.WriteSimpleString("OK")
		case "LPOP":
			log.Debug().Msg("client LPOP command")
			if command.ArgCount() != 2 {
				err = writer.WriteError(errInvalidArgCount.Error())
				break
			}

			var dtid uint64
			var subject string
			dtid, subject, err = parseKey(string(command.Get(1)))
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			var msg Message
			msg, err = c.LPop(dtid, subject)
			if err != nil {
				if errors.Is(err, errNoMessage) {
					err = writer.WriteBulksSlice(nil)
					break
				}
				err = writer.WriteError(err.Error())
				break
			}

			err = writer.WriteObjectsSlice([]interface{}{createKey(msg.Sender, msg.Topic), msg.Payload})
		case "LLEN":
			log.Debug().Msg("client LLEN command")
			if command.ArgCount() != 2 {
				err = writer.WriteError(errInvalidArgCount.Error())
				break
			}

			var dtid uint64
			var subject string
			dtid, subject, err = parseKey(string(command.Get(1)))
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			var count uint64
			count, err = c.LLen(dtid, subject)
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			err = writer.WriteInt(int64(count))
		case "LRANGE":
			log.Debug().Msg("client LRANGE command")
			if command.ArgCount() != 4 {
				err = writer.WriteError(errInvalidArgCount.Error())
				break
			}

			var dtid uint64
			var subject string
			dtid, subject, err = parseKey(string(command.Get(1)))
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			var start, end int
			start, err = strconv.Atoi(string(command.Get(2)))
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			end, err = strconv.Atoi(string(command.Get(3)))
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			var messages []Message
			messages, err = c.LRange(dtid, subject, start, end)
			if err != nil {
				err = writer.WriteError(err.Error())
				break
			}

			output := make([]interface{}, len(messages))
			for i := range messages {
				output[2*i] = createKey(messages[i].Sender, messages[i].Topic)
				output[2*i+1] = messages[i].Payload
			}
			err = writer.WriteObjectsSlice(output)
		default:
			log.Debug().Str("CMD", cmd).Msg("client sent unknown command")
			err = writer.WriteError(errInvalidCommand.Error())
		}

		if err != nil {
			log.Error().Err(err).Msg("could not write to connection")
			return
		}
	}
}

var (
	errInvalidCommand      = errors.New("unknown command")
	errInvalidArgCount     = errors.New("invalid amount of argument for command")
	errAuthorizationFailed = errors.New("authorization failed")
	errMalformedKey        = errors.New("malformed key")
)

const (
	serverVersion = "0.1.0"
	protoVersion  = 1
)

func (s *Server) helloInfo() []interface{} {
	// old style hello, trimmed down
	return []interface{}{
		"server",
		"tfagent",
		"version",
		serverVersion,
		"proto",
		protoVersion,
		"id",
		s.peerID(),
	}
}

func (s *Server) peerID() string {
	return s.node.PeerID()
}

// parseKey parses a key as expected by the protocol, and returns the dtid and
// subject from the key, or an error if parsing fails
func parseKey(key string) (uint64, string, error) {
	// key is formatted as <receiver_dtid>:<subject>
	// subject can be empty
	keyParts := strings.Split(key, keySeparator)
	if len(keyParts) != 2 {
		return 0, "", errMalformedKey
	}

	dtid, err := strconv.ParseUint(keyParts[0], 10, 64)
	if err != nil {
		return 0, "", errors.Wrap(err, "could not parse dtid")
	}

	return dtid, keyParts[1], nil
}

func createKey(dtid uint64, subject string) string {
	return fmt.Sprintf("%d%s%s", dtid, keySeparator, subject)
}
