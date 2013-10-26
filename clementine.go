// Example
//
//     package main
//
//     import "github.com/brunoqc/go-clementine"
//
//     func main() {
//         clementine := clementine.Clementine{
//             Host:     "127.0.0.1",
//             Port:     5500,
//             AuthCode: 28615,
//         }
//         errPause := clementine.SimpleStop()
//         if errPause != nil {
//             panic(errPause)
//         }
//     }
package clementine

import (
	"bufio"
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/brunoqc/go-clementine/pb_remote"
	"io"
	"net"
)

const (
	// TODO: rename (it's the size in bytes of the Message's size)
	sizeSize = 4
)

var (
	msgPlay = &pb_remote.Message{
		Type: pb_remote.MsgType_PLAY.Enum(),
	}
	msgPause = &pb_remote.Message{
		Type: pb_remote.MsgType_PAUSE.Enum(),
	}
	msgStop = &pb_remote.Message{
		Type: pb_remote.MsgType_STOP.Enum(),
	}
	msgDisconnect = &pb_remote.Message{
		Type: pb_remote.MsgType_DISCONNECT.Enum(),
	}
)

func recvMessage(conn net.Conn) (*pb_remote.Message, error) {
	buf := make([]byte, sizeSize)

	if _, err := io.ReadAtLeast(conn, buf, sizeSize); err != nil {
		return nil, err
	} else {
		lenbuf := bytes.NewBuffer(buf)
		var length int32

		if err = binary.Read(lenbuf, binary.BigEndian, &length); err != nil {
			return nil, err
		} else {
			buf = make([]byte, length)

			if _, err = io.ReadAtLeast(conn, buf, int(length)); err != nil {
				return nil, err
			} else {

				message := &pb_remote.Message{}
				if err = proto.Unmarshal(buf, message); err != nil {
					return nil, err
				} else {
					return message, nil
				}
			}
		}
	}
}

func sendMessage(conn net.Conn, msg *pb_remote.Message) error {
	if data, errMarshal := proto.Marshal(msg); errMarshal != nil {
		return errMarshal
	} else {
		size := new(bytes.Buffer)
		if err := binary.Write(size, binary.BigEndian, uint32(len(data))); err != nil {
			return err
		} else {
			writer := bufio.NewWriter(conn)
			writer.Write(size.Bytes())
			writer.Write(data)
			writer.Flush()
			return nil
		}
	}
}

type Clementine struct {
	Host     string
	Port     int
	AuthCode int
	conn     net.Conn
}

func (c *Clementine) connect() error {
	var errDial error
	c.conn, errDial = net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if errDial != nil {
		return errDial
	} else {
		msg := &pb_remote.Message{
			Type: pb_remote.MsgType_CONNECT.Enum(),
			RequestConnect: &pb_remote.RequestConnect{
				AuthCode:          proto.Int(c.AuthCode),
				SendPlaylistSongs: proto.Bool(false),
				Downloader:        proto.Bool(false),
			},
		}

		errSendMessage := sendMessage(c.conn, msg)
		if errSendMessage != nil {
			return errSendMessage
		} else {
			tmp, err := recvMessage(c.conn)
			switch {
			case err != nil:
				return err
			case tmp.GetType() == pb_remote.MsgType_DISCONNECT:
				return errors.New(tmp.GetResponseDisconnect().GetReasonDisconnect().String())
			default:
				return nil
			}
		}
	}
}

func (c *Clementine) sendSimplePlayPause(msg *pb_remote.Message) error {
	if errConnect := c.connect(); errConnect != nil {
		return errConnect
	} else if err := sendMessage(c.conn, msg); err != nil {
		return err
	} else {
		return sendMessage(c.conn, msgDisconnect)
	}
}

// SimplePlay connect to Clementine, send the 'Play' command and disconnect.
func (c *Clementine) SimplePlay() error {
	return c.sendSimplePlayPause(msgPlay)
}

// SimplePause connect to Clementine, send the 'Pause' command and disconnect.
func (c *Clementine) SimplePause() error {
	return c.sendSimplePlayPause(msgPause)
}

// SimpleStop connect to Clementine, send the 'Stop' command and disconnect.
func (c *Clementine) SimpleStop() error {
	return c.sendSimplePlayPause(msgStop)
}
