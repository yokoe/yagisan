package yagisan

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"golang.org/x/xerrors"
)

// Run starts observing files.
func Run() error {
	return watchFileWrite(func(path string) {
		log.Printf("File change: %v\n", path)

		msgs, err := runTest()
		if len(msgs) > 0 {
			for _, m := range msgs {
				log.Printf("%v\n", m)
			}
		} else if err != nil {
			log.Printf("Error: %v\n", err)
		}
	})
}

type fileChangeHandler func(string)

func watchFileWrite(handler fileChangeHandler) error {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	r := regexp.MustCompile("(.*).go")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				handler(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("."); err != nil {
		return xerrors.Errorf("failed to add . recursively: %w", err)
	}

	log.Println("start watching files...")
	if err := w.Start(time.Millisecond * 100); err != nil {
		return err
	}

	return nil
}

func runTest() ([]string, error) {
	out, err := exec.Command("go", "test", "./...").Output()
	s := string(out)

	errorMsgs := []string{}
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, "FAIL\t") {
			errorMsgs = append(errorMsgs, l)
		}
	}
	return errorMsgs, err
}
