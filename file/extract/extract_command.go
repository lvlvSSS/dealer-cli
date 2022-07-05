package extract

import (
	"dealer-cli/docs"
	LOG "dealer-cli/log"
	file_util "dealer-cli/utils/files"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

var extractFlags = []cli.Flag{
	headlineFlag,
	headlineYamlFlag,

	targetFlag,
	targetYamlFlag,

	fileFormatFlag,
	fileFormatYamlFlag,

	fileLocationFlag,
	fileLocationYamlFlag,

	xmlFlag,
	xmlYamlFlag,

	goroutinesFlag,
	goroutinesYamlFlag,

	fileSourceDirFlag,
	fileSourceDirYamlFlag,
}

var DefaultLogName = "./Logs/dealer-cli.log"
var DefaultLevel zapcore.Level = zapcore.InfoLevel

var ExtractCommand = &cli.Command{
	Name:  "extract",
	Usage: "Extract some specified message from a file and put it in another file",
	Flags: extractFlags,
	Action: func(c *cli.Context) error {
		if err := LOG.InitLog(c); err != nil {
			log.Warn(fmt.Sprintf("dealer_cli extract - log init errors[%s]", err))
			return err
		}
		if c.Bool("check") || c.Bool("dealer.check") {
			checkFlag(c)
			return nil
		}
		var fileSourceDir = ""
		if fileSourceDir = c.String("file-source-dir"); len(strings.TrimSpace(fileSourceDir)) != 0 {
			log.Debug(fmt.Sprintf("dealer_cli extract - file-source-dir[%s]", fileSourceDir))
		} else if fileSourceDir = c.String("dealer.file.extract.file-source-dir"); len(strings.TrimSpace(fileSourceDir)) != 0 {
			log.Debug(fmt.Sprintf("dealer_cli extract - file.extract.file-source-dir[%s]", fileSourceDir))
		} else {
			return errors.New("dealer_cli extract - file-source-dir is empty")
		}
		files, err := file_util.GetAllSubFiles(fileSourceDir)
		if err != nil {
			return err
		}
		wg := &sync.WaitGroup{}
		for _, file := range files {
			var fileExtractor = &GeneralFileExtractor{
				Messenger:      &LogXmlMessage{},
				FileSourcePath: file,
			}
			wg.Add(1)
			go func(group *sync.WaitGroup) {
				defer group.Done()
				if err := fileExtractor.Init(c); err != nil {
					log.Info(fmt.Sprintf("dealer-cli file[%s] extractor - Init errors : %s", fileExtractor.FileSourcePath, err))
					return
				}
				if err := fileExtractor.Extract(); err != nil {
					log.Info(fmt.Sprintf("dealer-cli file[%s] extractor - Extract errors : %s", fileExtractor.FileSourcePath, err))
					return
				}
				if err := fileExtractor.Done(); err != nil {
					log.Info(fmt.Sprintf("dealer-cli file[%s] extractor - Extract Done : %s", fileExtractor.FileSourcePath, err))
					return
				}
			}(wg)
		}
		wg.Wait()
		return nil
	},
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Command wrong, please use 'file extract help' to check ... \n")
		return err
	},
	Before: altsrc.InitInputSourceWithContext(extractFlags, altsrc.NewYamlSourceFromFlagFunc(docs.APP_LOAD_YAML)),
}

func checkFlag(c *cli.Context) {
	log.Warn("dealer-cli file extract - check flags begins ...")
	mode := c.String("mode")
	if len(strings.TrimSpace(mode)) == 0 {
		mode = "nil"
	}
	log.Info("mode : " + mode)

	modeYaml := c.String("dealer.mode")
	if len(strings.TrimSpace(modeYaml)) == 0 {
		modeYaml = "nil"
	}
	log.Info("mode Yaml : " + modeYaml)

	logPath := c.String("log")
	if len(strings.TrimSpace(logPath)) == 0 {
		logPath = "nil"
	}
	log.Info("log : " + logPath)

	logPathYaml := c.String("dealer.log")
	if len(strings.TrimSpace(logPathYaml)) == 0 {
		logPathYaml = "nil"
	}
	log.Info("log Yaml : " + logPathYaml)

	root, _ := os.Getwd()
	log.Info("root : " + root)
	headline := c.String("headline")
	if len(strings.TrimSpace(headline)) == 0 {
		headline = "nil"
	}
	log.Info("headline : " + headline)

	headlineYaml := c.String("dealer.file.extract.headline")
	if len(strings.TrimSpace(headlineYaml)) == 0 {
		headlineYaml = "nil"
	}
	log.Info("headline yaml : " + headlineYaml)

	target := c.String("target")
	if len(strings.TrimSpace(target)) == 0 {
		target = "nil"
	}
	log.Info("target : " + target)

	targetYaml := c.String("dealer.file.extract.target")
	if len(strings.TrimSpace(targetYaml)) == 0 {
		targetYaml = "nil"
	}
	log.Info("target yaml : " + targetYaml)

	fileFormat := c.String("file-format")
	if len(strings.TrimSpace(fileFormat)) == 0 {
		fileFormat = "nil"
	}
	log.Info("file-format : " + fileFormat)

	fileFormatYaml := c.String("dealer.file.extract.file-format")
	if len(strings.TrimSpace(fileFormatYaml)) == 0 {
		fileFormatYaml = "nil"
	}
	log.Info("file-format yaml : " + fileFormatYaml)

	lo := c.String("location")
	if len(strings.TrimSpace(lo)) == 0 {
		lo = "nil"
	}
	log.Info("location : " + lo)

	loYaml := c.String("dealer.file.extract.location")
	if len(strings.TrimSpace(loYaml)) == 0 {
		lo = "nil"
	}
	log.Info("location yaml : " + loYaml)

	xml := c.Bool("xml")
	log.Info(fmt.Sprintf("xml : %v", xml))

	xmlYaml := c.Bool("dealer.file.extract.xml")
	log.Info(fmt.Sprintf("xml yaml : %v", xmlYaml))

	goroutines := c.Int("goroutines")
	log.Info(fmt.Sprintf("goroutines : %d", goroutines))

	goroutinesYaml := c.Int("dealer.file.extract.goroutines")
	log.Info(fmt.Sprintf("goroutines yaml : %d", goroutinesYaml))

	fileSourceDir := c.String("file-source-dir")
	if len(strings.TrimSpace(fileSourceDir)) == 0 {
		fileSourceDir = "nil"
	}
	log.Info("file-source-dir : " + fileSourceDir)

	fileSourceDirYaml := c.String("dealer.file.extract.file-source-dir")
	if len(strings.TrimSpace(fileSourceDirYaml)) == 0 {
		fileFormatYaml = "nil"
	}
	log.Info("file-source-dir yaml : " + fileSourceDirYaml)
}
