package player

import (
	"encoding/binary"
	"fmt"
	"io"
)

func loadSound(reader io.ReadCloser) ([][]byte, error) {
	var buffer = make([][]byte, 0)

	var opuslen int16
	var err error

	for {
		// Read opus frame length from dca reader.
		err = binary.Read(reader, binary.LittleEndian, &opuslen)

		// If this is the end of the reader, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := reader.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca reader :", err)
			return nil, err
		}

		// Read encoded pcm from dca reader.
		InBuf := make([]byte, opuslen)
		err = binary.Read(reader, binary.LittleEndian, &InBuf)

		// Should not be any end of reader errors
		if err != nil {
			fmt.Println("Error reading from dca reader :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}
