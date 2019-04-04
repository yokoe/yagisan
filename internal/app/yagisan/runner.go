package yagisan

import (
	"log"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
	"golang.org/x/xerrors"
)

// Run starts observing files.
func Run() error {
	return watchFileWrite(func(path string) {
		log.Printf("File change: %v\n", path)
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
