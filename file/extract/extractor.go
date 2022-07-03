package extract

import (
	"bufio"
	"bytes"
	"dealer-cli/utils/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type FileExtractor interface {
	Init(c *cli.Context) error
	Extract() error
	Done() error
}

// MIN_FILE_SIZE the minus size of the file that one goroutine could analysis. default to 100MB
const MIN_FILE_SIZE = 100 * 1024 * 1024

// BUFFER_SIZE the buffer size for bufio.NewReaderSize. 100k is enough
var BUFFER_SIZE = 100 * 1024

func SetBufferSize(size int) {
	BUFFER_SIZE = size
}

type GeneralFileExtractor struct {
	Messenger
	Headline *regexp.Regexp // the headline of every message, using regex.

	FileSourcePath string
	reader         *os.File // the reader of FileSourcePath

	FileTargetFormat string // the output file format
	Location         string // the output directory.

	goroutines int // the max of goroutines numbers to analysis one file. every goroutine will analysis 10MB at least.
	wg         *sync.WaitGroup

	messages chan int // used to count the total message that extracted.
}

func New(target string,
	headline string,
	fileSourcePath string,
	fileTargetFormat string,
	location string,
	goroutines int) (*GeneralFileExtractor, error) {
	if len(strings.TrimSpace(target)) == 0 {
		return nil, errors.New("[dealer_cli.file.extract.New] target is empty")
	}

	if len(strings.TrimSpace(headline)) == 0 {
		return nil, errors.New("[dealer_cli.file.extract.New] headline is empty")
	}
	headlineRegex, err := regexp.Compile(headline)
	if err != nil {
		return nil, err
	}

	if len(strings.TrimSpace(fileSourcePath)) == 0 {
		return nil, errors.New("[dealer_cli.file.extract.New] fileSourcePath is empty")
	}
	fileSourcePathAbs, err := filepath.Abs(fileSourcePath)
	if err != nil {
		return nil, err
	}
	fileSourcePath = fileSourcePathAbs
	sourceFile, err := os.OpenFile(fileSourcePath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("[dealer_cli.file.extract.New] open file[%s] failed: [%s]",
				fileSourcePath, err.Error()))
	}

	if len(strings.TrimSpace(fileTargetFormat)) == 0 {
		return nil, errors.New("[dealer_cli.file.extract.New] fileTargetFormat is empty")
	}

	if len(strings.TrimSpace(location)) == 0 {
		return nil, errors.New("[dealer_cli.file.extract.New] location is empty")
	}
	locationAbs, err := filepath.Abs(location)
	if err != nil {
		return nil, err
	}
	location = locationAbs
	locStat, err := os.Stat(location)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.New(
			fmt.Sprintf("[dealer_cli.file.extract.New] init failed : location[%s] is not directory, error[%s]",
				location, err))
	} else if err != nil && os.IsNotExist(err) {
		os.MkdirAll(location, 0)
	} else if !locStat.IsDir() {
		return nil, errors.New(
			fmt.Sprintf("[dealer_cli.file.extract.New] init failed : location[%s] is not directory", location))
	}

	return &GeneralFileExtractor{
		Messenger:        &LogXmlMessage{Target: target},
		Headline:         headlineRegex,
		FileSourcePath:   fileSourcePath,
		reader:           sourceFile,
		FileTargetFormat: fileTargetFormat,
		Location:         location,
		goroutines:       goroutines,
		wg:               &sync.WaitGroup{},
	}, nil
}

func (extractor *GeneralFileExtractor) Init(c *cli.Context) error {
	if len(strings.TrimSpace(extractor.FileSourcePath)) == 0 {
		return fmt.Errorf(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.split] file source not init "))
	}
	fileSourcePathAbs, err := filepath.Abs(extractor.FileSourcePath)
	if err != nil {
		return err
	}
	extractor.FileSourcePath = fileSourcePathAbs
	sourceFile, err := os.OpenFile(extractor.FileSourcePath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] file[%s] not exists ", extractor.FileSourcePath))
		}
		return err
	}
	if fileStat, err := sourceFile.Stat(); err != nil || fileStat.IsDir() {
		if err == nil {
			log.Warn(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] file[%s] is directory", extractor.FileSourcePath))
			return fmt.Errorf("[dealer_cli.file.extract.GeneralFileExtractor.Init] file[%s] is directory", extractor.FileSourcePath)
		}
		return err
	}
	extractor.reader = sourceFile

	// according to the extension of file, set the file's category.
	if extension := filepath.Ext(extractor.FileSourcePath); len(strings.TrimSpace(extension)) != 0 {
		c.Set(extension, "true")
	}
	// set headline
	var headlineStr = ""
	if headlineStr = c.String("headline"); len(strings.TrimSpace(headlineStr)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] - headline[%s]", headlineStr))
	} else if headlineStr = c.String("file.extract.headline"); len(strings.TrimSpace(headlineStr)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] - file.extract.headline[%s]", headlineStr))
	} else {
		return fmt.Errorf("[dealer_cli.file.extract.GeneralFileExtractor.Init] init failed : headline[%s]", headlineStr)
	}
	headlineCompile, err := regexp.Compile(headlineStr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] compile string[%s] failed ", headlineStr))
	}
	extractor.Headline = headlineCompile

	// set goroutine number
	if extractor.goroutines = c.Int("goroutines"); extractor.goroutines < 1 {
		log.Warn("[dealer_cli.file.extract.GeneralFileExtractor.Init] goroutines number forced to 1 ")
		extractor.goroutines = 1
	}

	// set location
	if extractor.Location = c.String("location"); strings.TrimSpace(extractor.Location) != "" {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] - location[%s]", extractor.Location))
	} else if extractor.Location = c.String("file.extract.location"); strings.TrimSpace(extractor.Location) != "" {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] - file.extract.location[%s]", extractor.Location))
	} else {
		return fmt.Errorf("[dealer_cli.file.extract.GeneralFileExtractor.Init] init failed : location[%s]", extractor.Location)
	}
	locationAbs, err := filepath.Abs(extractor.Location)
	if err != nil {
		return err
	}
	extractor.Location = locationAbs
	locStat, err := os.Stat(extractor.Location)
	if err != nil && !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] init failed : location[%s] is not directory", extractor.Location))
	} else if err != nil && os.IsNotExist(err) {
		os.MkdirAll(extractor.Location, 0)
	} else if !locStat.IsDir() {
		return errors.New(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] init failed : location[%s] is not directory", extractor.Location))
	}

	// set file target format
	if extractor.FileTargetFormat = c.String("file-format"); len(strings.TrimSpace(extractor.FileTargetFormat)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] - file-format[%s]", extractor.FileTargetFormat))
	} else if extractor.FileTargetFormat = c.String("file.extract.file-format"); len(strings.TrimSpace(extractor.FileTargetFormat)) != 0 {
		log.Debug(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] - file.extract.file-format[%s]", extractor.FileTargetFormat))
	} else {
		return fmt.Errorf("[dealer_cli.file.extract.GeneralFileExtractor.Init] init failed : file-format[%s]", extractor.FileTargetFormat)
	}

	extractor.wg = &sync.WaitGroup{}

	// init the Messenger.
	if err := extractor.Messenger.Init(c); err != nil {
		log.Error(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Init] error: %s", err))
		return err
	}
	return nil
}

func (extractor *GeneralFileExtractor) Extract() error {
	files, err := extractor.split()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] can't extract file[%s]", extractor.FileSourcePath))
	}
	extractor.messages = make(chan int, len(files))
	for _, file := range files {
		extractor.wg.Add(1)
		oneSourceFile := file
		go func(waitGroup *sync.WaitGroup) {
			var totalMessages = 0
			defer oneSourceFile.close()
			defer waitGroup.Done()
			for {
				oneMessage, readErr := oneSourceFile.readOneMessageFromFile()
				if readErr != nil && readErr != io.EOF {
					log.Error(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] [%s] readOneMessageFromLog failed %s",
						oneSourceFile.String(), readErr.Error()))
					return
				} else if readErr != nil && readErr == io.EOF {
					log.Info(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] [%s] finished . ",
						oneSourceFile.String()))
					break
				}
				targetBytes, formatErr := extractor.Messenger.Format(oneMessage)
				if formatErr != nil {
					log.Error(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] [%s] Format failed %s",
						oneSourceFile.String(), formatErr.Error()))
					return
				}
				if targetBytes == nil {
					continue
				}
				totalMessages++
				dest, destErr := extractor.Messenger.Dest(extractor.FileTargetFormat, targetBytes)
				if destErr != nil {
					log.Error(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] [%s] Dest failed %s",
						oneSourceFile.String(), destErr.Error()))
					return
				}
				targetFile, openErr := os.OpenFile(filepath.Join(extractor.Location, dest), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
				if openErr != nil {
					log.Error(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] [%s] create file[%s] failed %s",
						oneSourceFile.String(), filepath.Join(extractor.Location, dest), openErr.Error()))
					return
				}
				extractor.Messenger.WriteTo(targetFile, targetBytes)
				targetFile.Sync()
				targetFile.Close()
			}
			log.Info(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Extract] file[%s] extract %d messages ...", oneSourceFile, totalMessages))
			extractor.messages <- totalMessages
		}(extractor.wg)
	}
	return nil
}

func (extractor *GeneralFileExtractor) Done() error {
	if extractor.messages == nil {
		log.Error("[dealer_cli.file.extract.GeneralFileExtractor.Done] GeneralFileExtractor.messages not initialized")
		return errors.New("[dealer_cli.file.extract.GeneralFileExtractor.Done] GeneralFileExtractor.messages not initialized")
	}
	extractor.wg.Wait()
	totalMessages := extractor.count()
	close(extractor.messages)
	log.Info(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.Done] file[%s] extracted into [%d] messages by %d goroutines",
		extractor.FileSourcePath, totalMessages, extractor.goroutines))
	return extractor.reader.Close()
}

func (extractor *GeneralFileExtractor) count() int {
	totalMessages := 0
	for {
		select {
		case mess := <-extractor.messages:
			totalMessages += mess
		default:
			return totalMessages
		}
	}
}

func (extractor *GeneralFileExtractor) split() ([]*generalFile, error) {
	if extractor.reader == nil {
		return nil, fmt.Errorf(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.split] reader not init "))
	}
	info, err := extractor.reader.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.split] file[%s] not exists", info.Name()))
		}
		return nil, err
	}
	// the algorithm to define the number of goroutines:
	// 1. if the size is less than MIN_FILE_SIZE bytes 1.5 times, then 1,
	// 2. the size divided by MIN_FILE_SIZE to get the count, if the count is less than extractor.goroutines, then count.
	if extractor.goroutines <= 0 {
		extractor.goroutines = runtime.NumCPU()
	}
	if info.Size() <= (1.5 * MIN_FILE_SIZE) {
		extractor.goroutines = 1
	} else if counts := int(info.Size() / MIN_FILE_SIZE); counts < extractor.goroutines {
		extractor.goroutines = counts
	}

	offsets := make([]int64, extractor.goroutines, extractor.goroutines)
	for offsetIndex, _ := range offsets {
		newOffset, err := extractor.offset(int64(MIN_FILE_SIZE * offsetIndex))
		if err != nil {
			log.Error(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.split] get new offset error[%s] ", err.Error()))
			return nil, err
		}
		offsets[offsetIndex] = newOffset
		if newOffset == info.Size() {
			log.Warn(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.split] offsets[%d] is less than goroutines[%d]", offsetIndex+1, extractor.goroutines))
			break
		}
	}
	if len(offsets) == extractor.goroutines {
		offsets = append(offsets, info.Size())
	}

	totalFiles := len(offsets) - 1
	extractor.goroutines = totalFiles
	newReaders := make([]*generalFile, 0, totalFiles)
	for newReaderIndex := 0; newReaderIndex < totalFiles; newReaderIndex++ {
		newReader, err := os.OpenFile(extractor.FileSourcePath, os.O_RDONLY, os.ModeAppend)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.split] open file[%s] failed in [%v]time", extractor.FileSourcePath, newReaderIndex))
		}

		newReaders = append(newReaders,
			&generalFile{
				start:    offsets[newReaderIndex],
				offset:   offsets[newReaderIndex],
				end:      offsets[newReaderIndex+1],
				reader:   newReader,
				headline: extractor.Headline,
				index:    newReaderIndex,
				total:    totalFiles,
			})
	}
	return newReaders, nil
}
func (extractor *GeneralFileExtractor) offset(originOffset int64) (int64, error) {
	if extractor.reader == nil {
		return -1, fmt.Errorf(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.offset] no source file "))
	}
	info, err := extractor.reader.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.offset] file[%s] not exists", info.Name()))
		}
		return -1, err
	}
	_, err = extractor.reader.Seek(originOffset, 0)
	if err != nil {
		return -1, err
	}
	ioReader := bufio.NewReaderSize(extractor.reader, BUFFER_SIZE)
	readBytes, err := readlineNumberByRegexTemplate(ioReader, extractor.Headline)
	if err != nil {
		return -1, err
	}
	return originOffset + int64(readBytes), nil
}

func readlineNumberByRegexTemplate(reader *bufio.Reader, template *regexp.Regexp) (int, error) {
	var readBytes = 0
	for {
		slice, err := reader.ReadSlice('\n')
		if err != nil && err != io.EOF {
			return -1, errors.Wrap(err, fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.readlineNumberByRegexTemplate] read line failed"))
		}
		if template.Match(slice) {
			break
		}
		readBytes += len(slice)
		if err == io.EOF {
			break
		}
	}
	return readBytes, nil
}

type generalFile struct {
	start    int64
	offset   int64
	end      int64
	reader   *os.File
	headline *regexp.Regexp
	index    int // the index of the goroutines numbers of os.File
	total    int // total number of the goroutines number
}

// readOneMessageFromFile - read a message from a file. that message begins with the headline.
func (file *generalFile) readOneMessageFromFile() ([]byte, error) {
	if file.offset >= file.end {
		return nil, io.EOF
	}
	seek, seekErr := file.reader.Seek(file.offset, 0)
	if seekErr != nil {
		return nil, errors.Wrap(seekErr, fmt.Sprintf("[dealer_cli.file.extract.generalFile.readOneMessageFromFile] file[%s] seek[%d] error. ", file.reader.Name(), file.offset))
	}
	if seek != file.offset {
		log.Warn(fmt.Sprintf("[dealer_cli.file.extract.generalFile.readOneMessageFromFile] file[%s] seek[%d] returns wrong[%d].", file.reader.Name(), file.offset, seek))
		file.offset = seek
	}

	bufReader := bufio.NewReaderSize(file.reader, BUFFER_SIZE)
	slice, err := bufReader.ReadSlice('\n')
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, fmt.Sprintf("[dealer_cli.file.extract.generalFile.readOneMessageFromFile] read line failed"))
	}
	//log.Info(fmt.Sprintf("out - readOneMessageFromFile : %s", converter.BytesToString(slice)))
	if !file.headline.Match(slice) {
		return nil, errors.New(fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.readOneMessageFromFile] the beginning in file[%s] is not fit for the head[%s]", file.reader.Name(), file.headline.String()))
	}
	bytesBuffer := bytes.NewBuffer(make([]byte, 0, BUFFER_SIZE))
	bytesBuffer.Write(slice)
	file.offset += int64(len(slice))
	for {
		if err == io.EOF || file.offset >= file.end {
			break
		}
		slice, err = bufReader.ReadSlice('\n')
		//log.Info(fmt.Sprintf("in - readOneMessageFromFile : %s", converter.BytesToString(slice)))
		if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, fmt.Sprintf("[dealer_cli.file.extract.GeneralFileExtractor.readOneMessageFromFile] try to read one message failed"))
		}
		if file.headline.Match(slice) {
			break
		}
		if len(slice) > 0 {
			bytesBuffer.Write(slice)
			file.offset += int64(len(slice))
		}
	}
	return bytesBuffer.Bytes(), nil
}

func (file *generalFile) close() {
	file.reader.Close()
	file.offset = -1
	file.end = -1
	file.reader = nil
}

func (file *generalFile) String() string {
	return fmt.Sprintf("generalFile{start[%d],offset[%d],end[%d],index[%d],total[%d],file[%s],headline[%s]}",
		file.start, file.offset, file.end, file.index, file.total, file.reader.Name(), file.headline.String())
}
