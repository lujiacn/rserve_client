/*Rserve go client */
package rservcli

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

/*
 return mkint32($cmd) . mkint32($n + 4) . mkint32(0) . mkint32(0) . chr(4) . mkint24($n) . $string;
*/
type PhDr struct {
	Cmd   int32 /* command */
	Len   int32 /* length of the packet minus header (ergo -16) */
	MsgId int32 /* message id (since 1.8) [WAS:data offset behind header (ergo usually 0)] */
	Res   int32 /* high 32-bit of the packet length (since 0103
	and supported on 64-bit platforms only)
	aka "lenhi", but the name was not changed to
	maintain compatibility */
}

//MkpStr
func MakeBytes(cmd int32, data string) []byte {
	data = strings.Replace(data, "\r", "\n", -1)
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, cmd) //eval
	binary.Write(buff, binary.LittleEndian, int32(len(data)+4))
	binary.Write(buff, binary.LittleEndian, int32(0))
	binary.Write(buff, binary.LittleEndian, int32(0))
	binary.Write(buff, binary.LittleEndian, int8(4))
	binary.Write(buff, binary.LittleEndian, int32(len(data)))
	buff.Truncate(len(buff.Bytes()) - 1) // truncate 36bit int to 24 bit int -- 3 bytes
	// data string

	binary.Write(buff, binary.LittleEndian, []byte(data))
	binary.Write(buff, binary.LittleEndian, []byte("\r\n\r\n"))
	return buff.Bytes()
}

func RInit() (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:6311")
	if err != nil {
		fmt.Printf("Cannot connect to %v", err)
	}
	return conn, err
}

func REval(data string) ([]byte, error) {
	conn, err := RInit()
	if err != nil {
		fmt.Printf("%v", err)
	}
	send := MakeBytes(2, data)
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, send)
	binary.Write(buff, binary.LittleEndian, int32(len(send)))
	binary.Write(buff, binary.LittleEndian, int32(0))
	n, err := conn.Write(buff.Bytes())
	if err != nil {
		fmt.Printf("%v", err)

	}
	out := make([]byte, n)
	_, err = conn.Read(out)
	conn.Close()
	return out, err

}
