package extract

import (
	"dealer-cli/docs"
	LOG "dealer-cli/log"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
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
		LOG.InitLog(c)
		if c.Bool("check") {
			checkFlag(c)
			return nil
		}
		var fileSourceDir = ""
		if fileSourceDir = c.String("file-source-dir"); len(strings.TrimSpace(fileSourceDir)) != 0 {
			log.Debug(fmt.Sprintf("dealer_cli extract file-source-dir[%s]", fileSourceDir))
		} else if fileSourceDir = c.String("file.extract.file-source-dir"); len(strings.TrimSpace(fileSourceDir)) != 0 {
			log.Debug(fmt.Sprintf("dealer_cli extract file.extract.file-source-dir[%s]", fileSourceDir))
		} else {
			return errors.New("file.extract.file-source-dir is empty")
		}
		files, err := getAllSubFiles(fileSourceDir)
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
					log.Info(fmt.Sprintf("dealer-cli file[%s] extractor Init errors : %s", fileExtractor.FileSourcePath, err))
					return
				}
				if err := fileExtractor.Extract(); err != nil {
					log.Info(fmt.Sprintf("dealer-cli file[%s] extractor Extract errors : %s", fileExtractor.FileSourcePath, err))
					return
				}
				if err := fileExtractor.Done(); err != nil {
					log.Info(fmt.Sprintf("dealer-cli file[%s] extractor Extract Done : %s", fileExtractor.FileSourcePath, err))
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

func getAllSubFiles(parent string) ([]string, error) {
	result := make([]string, 0, 16)
	stat, err := os.Stat(parent)
	if err != nil || !stat.IsDir() {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("extract file[%s] fails, because it is not directory ", parent))
	}
	files, err := os.ReadDir(parent)
	for _, file := range files {
		if file.IsDir() {
			subParent := filepath.Join(parent, file.Name())
			subFiles, err := getAllSubFiles(subParent)
			if err != nil {
				continue
			} else {
				result = append(result, subFiles...)
			}
		} else {
			result = append(result, filepath.Join(parent, file.Name()))
		}
	}
	return result, nil
}

func checkFlag(c *cli.Context) {
	log.Warn("dealer-cli file extract - check flags begins ...")
	getwd, _ := os.Getwd()
	log.Info("root : " + getwd)
	headline := c.String("headline")
	if len(strings.TrimSpace(headline)) == 0 {
		headline = "nil"
	}
	log.Info("headline : " + headline)

	headlineYaml := c.String("file.extract.headline")
	if len(strings.TrimSpace(headlineYaml)) == 0 {
		headlineYaml = "nil"
	}
	log.Info("headline yaml : " + headlineYaml)

	target := c.String("target")
	if len(strings.TrimSpace(target)) == 0 {
		target = "nil"
	}
	log.Info("target : " + target)

	targetYaml := c.String("file.extract.target")
	if len(strings.TrimSpace(targetYaml)) == 0 {
		targetYaml = "nil"
	}
	log.Info("target yaml : " + targetYaml)

	fileFormat := c.String("file-format")
	if len(strings.TrimSpace(fileFormat)) == 0 {
		fileFormat = "nil"
	}
	log.Info("file-format : " + fileFormat)

	fileFormatYaml := c.String("file.extract.file-format")
	if len(strings.TrimSpace(fileFormatYaml)) == 0 {
		fileFormatYaml = "nil"
	}
	log.Info("file-format yaml : " + fileFormatYaml)

	lo := c.String("location")
	if len(strings.TrimSpace(lo)) == 0 {
		lo = "nil"
	}
	log.Info("location : " + lo)

	loYaml := c.String("file.extract.location")
	if len(strings.TrimSpace(loYaml)) == 0 {
		lo = "nil"
	}
	log.Info("location yaml : " + loYaml)

	xml := c.Bool("xml")
	log.Info(fmt.Sprintf("xml : %v", xml))

	xmlYaml := c.Bool("file.extract.xml")
	log.Info(fmt.Sprintf("xml yaml : %v", xmlYaml))

	goroutines := c.Int("goroutines")
	log.Info(fmt.Sprintf("goroutines : %d", goroutines))

	goroutinesYaml := c.Int("file.extract.goroutines")
	log.Info(fmt.Sprintf("goroutines yaml : %d", goroutinesYaml))

	fileSourceDir := c.String("file-source-dir")
	if len(strings.TrimSpace(fileSourceDir)) == 0 {
		fileSourceDir = "nil"
	}
	log.Info("file-source-dir : " + fileSourceDir)

	fileSourceDirYaml := c.String("file.extract.file-source-dir")
	if len(strings.TrimSpace(fileSourceDirYaml)) == 0 {
		fileFormatYaml = "nil"
	}
	log.Info("file-source-dir yaml : " + fileSourceDirYaml)
}