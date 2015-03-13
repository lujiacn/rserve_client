package rservcli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestREval(t *testing.T) {
	t_byte, _ := ioutil.ReadFile("test/data.txt")
	t_file, _ := ioutil.TempFile(os.TempDir(), "test_")
	//Make String
	data := string(t_byte)
	data = strings.Replace(data, "\r", "\n", -1)
	data = strings.Replace(data, "&&file_name", t_file.Name(), -1)
	fmt.Println(t_file.Name())
	REval(data)

}
