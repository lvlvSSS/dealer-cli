package extract

import (
	"bufio"
	file_util "dealer-cli/utils/files"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestReadByteAndOffset(t *testing.T) {
	source := `100000008A68FV] 2022-06-13 00:05:41,401 [My-Thread-001] INFO  sdfsdfsf13918723hksjds,mx.mcvxcvj sd fsldf019329012-slsdkf sd sdf :null
100000008A68FV] 2022-06-13 00:05:41,401 [My-Thread-001] INFO  sdfsdfsf13918723hksjds,mx.mcvxcvj sd fsldf019329012-slsdkf sd sdf :null
[478d41d2298f4de8b004afff6db771fd] [] [1001G3100000008A68FV] 2022-06-13 00:07:15,259 [My-Thread-001] INFO  sdfsfsfsfsdfgfghfghfg]sdf\] \sx;x;cvx/.cvxc;vcx1 123;'s;'fsd   sdf-s=df1	sdfsfDSFSDF:null
`
	target := `[478d41d2298f4de8b004afff6db771fd] [] [1001G3100000008A68FV] 2022-06-13 00:07:15,259 [My-Thread-001] INFO  sdfsfsfsfsdfgfghfghfg]sdf\] \sx;x;cvx/.cvxc;vcx1 123;'s;'fsd   sdf-s=df1	sdfsfDSFSDF:null
`

	headlineFormat := `^\[[\S\s]+\]\s\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2},\d{3}\s\[\S+?\]\s[A-Z]+\s[\s\S]+`
	compile := regexp.MustCompile(headlineFormat)
	number, err := readlineNumberByRegexTemplate(bufio.NewReaderSize(strings.NewReader(source), 1024*10), compile)
	if err != nil {
		t.Logf(err.Error() + "\n")
		t.Fail()
		return
	}
	result := source[number:]
	t.Logf("result : %s\n", result)
	if result != target {
		t.Fail()
	}
}

func TestBufferChan(t *testing.T) {
	messages := make(chan int, 10)
	for i := 0; i < 10; i++ {
		messages <- i
	}
re:
	for {
		select {
		case re := <-messages:
			fmt.Println("received: " + strconv.Itoa(re))
		default:
			fmt.Println("default end")
			break re
		}
	}

}

func TestFilePathAbs(t *testing.T) {
	root, _ := os.Getwd()
	rootAbs, _ := filepath.Abs(root)
	if root != rootAbs {
		t.Fail()
	}
	t.Logf("root : %s \n", rootAbs)

	abs, _ := filepath.Abs(".\\extractor_test.go")
	t.Logf("root : %s \n", abs)
}

func TestGetAllSubFiles(t *testing.T) {
	current, _ := os.Getwd()
	t.Logf("current path : %s \n", current)
	join := filepath.Join(current, "\\", "..")
	files, _ := file_util.GetAllSubFiles(join)
	fmt.Printf("All sub files : %v", files)
	if len(files) != 7 {
		t.Fail()
	}
	t.Logf("join path : %s \n", join)
	filepath.WalkDir(join, func(path string, d fs.DirEntry, err error) error {
		t.Logf("file : %s \n", path)
		return nil
	})
}
