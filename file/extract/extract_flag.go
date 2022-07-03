package extract

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"path/filepath"
)

var location, _ = os.Getwd()
var headlineFlag = &cli.StringFlag{
	Name:     "headline",
	Usage:    "Specify the head of each message",
	FilePath: filepath.Join(location, "./file.extract.headline.dealer"),
}
var headlineYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "file.extract.headline",
	Usage: "Specify the head of each message, same as the flag 'headline'",
})

var targetFlag = &cli.StringFlag{
	Name:     "target",
	Usage:    "Specifies the target content to extract",
	FilePath: filepath.Join(location, "./file.extract.target.dealer"),
}
var targetYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "file.extract.target",
	Usage: "Specifies the target content to extract, same as the flag 'target'",
})

var fileFormatFlag = &cli.StringFlag{
	Name:     "file-format",
	Usage:    "Specifies the format of the file generated after extraction",
	FilePath: filepath.Join(location, "./file.extract.file-format.dealer"),
}
var fileFormatYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "file.extract.file-format",
	Usage: "Specifies the format of the file generated after extraction, same as the flag 'file-format' ",
})

var fileLocationFlag = &cli.StringFlag{
	Name:     "location",
	Usage:    "Specifies the location of the file generated after extraction",
	FilePath: filepath.Join(location, "./file.extract.location.dealer"),
}
var fileLocationYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "file.extract.location",
	Usage: "Specifies the location of the file generated after extraction, same as the flag 'location' ",
})

var xmlFlag = &cli.BoolFlag{
	Name:  "xml",
	Usage: "Specifies that the string of source is xml",
	Value: false,
}

var xmlYamlFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:  "file.extract.xml",
	Usage: "Specifies that the string of source is xml, same as the flag 'xml' ",
	Value: false,
})

var goroutinesFlag = &cli.IntFlag{
	Name:  "goroutines",
	Usage: "the max of goroutines numbers to analysis one file",
	Value: 1,
}

var goroutinesYamlFlag = altsrc.NewIntFlag(&cli.IntFlag{
	Name:  "file.extract.goroutines",
	Usage: "the max of goroutines numbers to analysis one file, same as the flag 'goroutines'",
})

var fileSourceDirFlag = &cli.StringFlag{
	Name:     "file-source-dir",
	Usage:    "Specifies the directory of the files source",
	FilePath: filepath.Join(location, "./file.source.dir.benchmark"),
}

var fileSourceDirYamlFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "file.extract.file-source-dir",
	Usage: "Specifies the directory of the files source, same as the flag 'file-source-dir' ",
})
