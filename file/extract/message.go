package extract

import (
	"bytes"
	"dealer-cli/utils/converter"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/beevik/etree"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io"
	"regexp"
	"strings"
)

// Messenger - could handle the data string to target file locations
// Messenger should be thread-safe.
type Messenger interface {
	Init(c *cli.Context) error
	WriteTo(writer io.Writer, source []byte) error
	Dest(destFormat string, target []byte) (string, error)
	Format(source []byte) ([]byte, error)
	Clone() Messenger
}

var fileXmlFormatPattern = `#\{(\S+?)\}`

// LogXmlMessage - extract the target string from source, and then write to the location with specified file name.
type LogXmlMessage struct {
	Target string // the extracting target of every message, using the regex
}

// Init - suggest to use Init to create LogXmlMessage object.
func (message *LogXmlMessage) Init(c *cli.Context) error {
	// xml flag - is not set, then can't use the LogXmlMessage to handle.
	if !c.Bool("xml") && !c.Bool("dealer.file.extract.xml") {
		return fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.Init] target file is not xml ")
	}

	if message.Target = c.String("target"); len(strings.TrimSpace(message.Target)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.LogXmlMessage.Init] - target[%s]", message.Target))
	} else if message.Target = c.String("dealer.file.extract.target"); len(strings.TrimSpace(message.Target)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.LogXmlMessage.Init] - file.extract.target[%s]", message.Target))
	} else {
		return fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.Init] init failed : targetFormat[%s]", message.Target)
	}
	return nil
}

func (message *LogXmlMessage) WriteTo(writer io.Writer, source []byte) error {
	bytesWrite, err := writer.Write(source)
	if err != nil {
		return err
	}
	if bytesWrite != len(source) {
		return errors.New(fmt.Sprintf("[dealer_cli.file.extract.LogXmlMessage.WriteTo] source[%s] not fully written", converter.BytesToString(source)))
	}
	return nil
}

// Dest
// destFormat - the format of file name.
// target - the source xml string.
func (message *LogXmlMessage) Dest(destFormat string, target []byte) (string, error) {
	compile, err := regexp.Compile(fileXmlFormatPattern)
	if err != nil {
		return "",
			fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.Format] file format failed : file format[%s], error[%#v]",
				fileXmlFormatPattern, err)
	}
	originFileFormats := compile.FindAllStringSubmatch(destFormat, -1)
	fileFormatPaths := make([]string, 0, len(originFileFormats))
	for _, f := range originFileFormats {
		if len(f) != 2 {
			return "",
				fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.Format] file format failed : length of each originFileFormats[%#v] not equals 2", originFileFormats)
		}
		fileFormatPaths = append(fileFormatPaths, f[1])
	}
	fileFormats := make([]string, 0, len(fileFormatPaths))
	doc := etree.NewDocument()
	err = doc.ReadFromBytes(target)
	if err != nil {
		return "", fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.generateFileName] load xml failed : xml[%s], error[%#v]", converter.BytesToString(target), err)
	}

	for _, formatPath := range fileFormatPaths {
		elements := doc.FindElements(formatPath)
		fileFormats = append(fileFormats,
			linq.From(elements).Select(func(element interface{}) interface{} {
				return element.(*etree.Element).Text()
			}).Distinct().Aggregate(func(element1 interface{}, element2 interface{}) interface{} {
				return element1.(string) + "," + element2.(string)
			}).(string))
	}
	for i, path := range originFileFormats {
		destFormat = strings.Replace(destFormat, path[0], fileFormats[i], 1)
	}
	return destFormat, nil
}

func (message *LogXmlMessage) Format(source []byte) ([]byte, error) {
	if source = bytes.TrimSpace(source); len(source) == 0 {
		return nil, errors.New("[dealer_cli.file.extract.LogXmlMessage.Init] no Source ... ")
	}
	if message.Target = strings.TrimSpace(message.Target); len(message.Target) == 0 {
		return nil, fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.Format] Format failed : no target")
	}
	// get the target regex.
	compile, err := regexp.Compile(message.Target)
	if err != nil {
		return nil, fmt.Errorf("[dealer_cli.file.extract.LogXmlMessage.Format] Format failed : Compile[%s], error[%#v]", message.Target, err)
	}
	targetBytes := compile.Find(source)
	if len(bytes.TrimSpace(targetBytes)) == 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.LogXmlMessage.Format] : LogXmlMessage.source[%s], the target is not found", converter.BytesToString(source)))
		return nil, nil
	}
	return targetBytes, nil
}

func (message *LogXmlMessage) Clone() Messenger {
	newMessage := *message
	return &newMessage
}
