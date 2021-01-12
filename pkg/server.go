package pkg

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/secmask/go-redisproto"
)

// Server accepting connections, running RESP with custom commands
type Server struct {
	ps PeerStore

	ln net.Listener

	ctx context.Context
}

// CreateServer creates a new server. This binds the given port, but does not
// yet accept incomming connections
func CreateServer(ctx context.Context, port uint16, ps PeerStore) (*Server, error) {
	s := &Server{ps: ps}

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
	writer := redisproto.NewWriter(bufio.NewWriter(conn))

	for {
		// Don't use the `Commands()` channel here, as that exits on any error,
		// including protocol errors
		command, err := parser.ReadCommand()
		if err != nil {
			if errors.Is(err, &redisproto.ProtocolError{}) {
				writer.WriteError(err.Error())
			} else {
				log.Error().Err(err).Msg("failed to read command")
				return
			}
		}

		var writeErr error
		cmd := strings.ToUpper(string(command.Get(0)))
		switch cmd {
		case "HELLO":
			if command.ArgCount() != 1 {
				writeErr = writer.WriteError(errInvalidArgCount.Error())
				break
			}
			writeErr = writer.WriteObjectsSlice(s.helloInfo())
		case "AUTH":
			if command.ArgCount() != 3 {
				writeErr = writer.WriteError(errInvalidArgCount.Error())
				break
			}

			// first arg is the dtid
			dtid, err := strconv.ParseUint(string(command.Get(1)), 10, 64)
			if err != nil {
				writeErr = writer.WriteError(err.Error())
				break
			}

			// second arg is the signature. Accept either a raw byte array of
			// len64, or a byte array of len 128, which we then assume is the
			// hex encoded signature
			rawSig := command.Get(2)
			var sig [64]byte
			switch len(rawSig) {
				case 64:
					copy(sig[:], rawSig)
				case 128:
					var data []byte
					data, err = hex.DecodeString(string(rawSig))
					if err != nil {
						break
					}
					copy(sig[:], data)
				default:
					err = errInvalidSignatureLength
			}
			if err != nil {
			writeErr = writer.WriteError(err.Error())
			break
			}

			pk, err := s.ps.PublicKey(dtid)
			if err != nil {
					writeErr = writer.WriteError("invalid signature length")
					break
			}

			if signatureValid(pk, sig) {
				writeErr = writer.WriteSimpleString("OK")
				// TODO upgrade to authenticated conn
			} else {
				writeErr = writer.WriteError("authentication failed")
			}
		case "LPUSH":
		case "LPOP":
		case "LLEN":
		case "LRANGE":
		default:
			writeErr = writer.WriteError(errInvalidCommand.Error())
		}

		if writeErr != nil {
			log.Error().Err(writeErr).Msg("could not write to connection")
			return
		}
	}
}


var errInvalidCommand = errors.New("unknown command")
var errInvalidArgCount = errors.New("invalid amount of argument for command")
var errInvalidSignatureLength = errors.New("invalid signature length")
var errAuthorizationFailed = errors.New("authorization failed")

const serverVersion = "0.1.0"
const protoVersion = 1

func (s *Server) helloInfo() []interface{} {
	// old style hello, trimmed down
	return []interface{}{"server", "tfagent", "version", serverVersion, "proto", protoVersion, "id", s.peerID()}
}

func (s *Server) peerID() string {
	return "//TODO"
}
