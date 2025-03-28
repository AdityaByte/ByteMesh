package client

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
)

type metaData struct {
	Filename      string
	FileExtension string
	Location      map[string]string
}

// The Download() function takes the name of the file with extension
// as parameter and sends the request to the namenode server so if the data exists
// at the server and the filename matches with any of the file which is at the server then we
// will get the request with request code 200 ok.

const addr = ":9004"

func Download(filename string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error occured while creating connection to the namenodeserver", err)
	}

	conn.Write([]byte("GET\n dfs-flowchart.png\n"))

	reader := bufio.NewReader(conn)
	statusCode, err := reader.ReadString('\n')

	statusCode = strings.TrimSpace(statusCode)
	fmt.Println("Status code is ", statusCode)
	fmt.Printf("Type of status code %T\n", statusCode)
	if err != nil {
		return fmt.Errorf("Error occured at client side:", err)
	}
	if statusCode == "200" {
		decoder := gob.NewDecoder(conn)
		var recievedMetaData *metaData
		err := decoder.Decode(&recievedMetaData)

		if err != nil {
			return fmt.Errorf("error:", err)
		}

		fmt.Println("Recieved meta data:", recievedMetaData)
	}
	return nil
}
