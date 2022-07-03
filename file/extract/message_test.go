package extract

import (
	"dealer-cli/utils/converter"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestFileXmlFormatPattern(t *testing.T) {
	str := "#{abc.ded}-#{cde.abc}"
	result := []string{"abc.ded", "cde.abc"}
	compile, _ := regexp.Compile(fileXmlFormatPattern)
	allString := compile.FindAllStringSubmatch(str, -1)
	fmt.Println(len(allString))
	for index, str := range allString {
		if len(str) != 2 || str[1] != result[index] {
			t.Fail()
		}
	}
}

func TestLogXmlMessage_Clone(t *testing.T) {
	message := &LogXmlMessage{
		Target: "2",
	}
	newMessage := message.Clone()
	t.Logf("old: %p, new: %p", message, newMessage)
	if message == newMessage {
		t.Fail()
	}
}
func TestMessageFormat(t *testing.T) {
	source := `[958f8afff1e4403d861d3963bb50f368] [8a5e2d44-da45-432b-a3aa-2c5143894976] [1001G3100000008A64LB] 2022-04-18 11:41:45,339 
INFO the xml content is ：【<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ROOT>
    <Location>DEMO</Location>
	<OPTYPE>501</OPTYPE>
    <User>
		<User_Name>Lily</User_Name>
		<User_Age>22</User_Age>
	</User>
    <User>
		<User_Name>Lucy</User_Name>
		<User_Age>23</User_Age>
	</User>
</ROOT>
】`
	wrongSource := `[958f8afff1e4403d861d3963bb50f368] [8a5e2d44-da45-432b-a3aa-2c5143894976] [1001G3100000008A64LB] 2022-04-18 11:41:45,339 
INFO the xml content is ：【<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ROOT>
    <Location>DEMO</Location>	
    <User>
		<User_Name>Lily</User_Name>
		<User_Age>22</User_Age>
	</User>
    <User>
		<User_Name>Lucy</User_Name>
		<User_Age>23</User_Age>
	</User>
</ROOT>
】`
	target_str := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ROOT>
    <Location>DEMO</Location>
	<OPTYPE>501</OPTYPE>
    <User>
		<User_Name>Lily</User_Name>
		<User_Age>22</User_Age>
	</User>
    <User>
		<User_Name>Lucy</User_Name>
		<User_Age>23</User_Age>
	</User>
</ROOT>`

	message := &LogXmlMessage{
		Target: `<\?xml version="1.0" encoding="UTF-8" standalone="yes"\?>\s*<ROOT>[\s\S]+<OPTYPE>501</OPTYPE>[\s\S]+</ROOT>`,
	}
	fileFormat := `#{./ROOT/User/User_Name} - #{./ROOT/User/User_Age}.log`
	targetFileName := `Lily,Lucy - 22,23.log`
	formatStr, _ := message.Format(converter.StringToBytes(wrongSource))
	if formatStr != nil {
		t.Fail()
		return
	}
	result, _ := message.Format(converter.StringToBytes(source))
	if converter.BytesToString(result) != target_str {
		t.Fail()
		return
	}

	fileXmlName, err := message.Dest(fileFormat, result)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
		return
	}
	if strings.Compare(fileXmlName, targetFileName) != 0 {
		t.Logf("origin[%s] , target[%s] not equals.", fileXmlName, targetFileName)
		t.Fail()
		return
	}

}
