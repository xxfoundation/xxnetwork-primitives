package rateLimiting

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/utils"
	"strings"
	"sync"
	"time"
)

// Whitelist structure contains a map of whitelist items and a file to update
// that list from.
type Whitelist struct {
	list         map[string]bool   // Contains the list of keys
	file         string            // Absolute URL for the whitelist file
	watcher      *fsnotify.Watcher // Watches for file changes
	sync.RWMutex                   // Only allows one writer at a time
}

// CreateWhitelistFile creates a whitelist file from hardcoded data so that the
// other functions can process the whitelist data.
// FIXME: HACK HACK HACK HACK
func CreateWhitelistFile(filePath string, whitelist []string) error {
	// Convert array of strings to single string with a decimeter of a new line
	dataString := strings.Join(whitelist, "\n")

	// Write the list to file
	err := utils.WriteFile(filePath, []byte(dataString), utils.FilePerms, utils.DirPerms)
	if err != nil {
		return errors.Errorf("failed to create whitelist file %s: %v", filePath, err)
	}

	return nil
}

// InitWhitelist initialises a map for the whitelist with the keys from the
// specified file path. The updateFinished channel receives a value when the map
// finises updating.
func InitWhitelist(filePath string, updateFinished chan bool) (*Whitelist, error) {
	newWhitelist := Whitelist{
		list: make(map[string]bool),
		file: filePath,
	}

	// Update the whitelist from the specified file
	err := newWhitelist.UpdateWhitelist()
	if err != nil {
		return nil, err
	}

	// Start watching the whitelist file for changes
	go newWhitelist.WhitelistWatcher(updateFinished)

	return &newWhitelist, nil
}

// UpdateWhitelist initialises a map for the whitelist with the keys found in
// the file.
func (wl *Whitelist) UpdateWhitelist() error {
	// Get list of strings from the whitelist file
	list, err := WhitelistFileParse(wl.file)
	if err != nil {
		return err
	}

	// If the file was read successfully, then update the list
	wl.Lock()
	defer wl.Unlock()

	// Reset the whitelist to empty
	wl.list = make(map[string]bool)

	// Add all the keys to the map
	for _, key := range list {
		wl.list[key] = true
	}

	return nil
}

// WhitelistFileParse parses the given file and stores each value in a slice.
// The file is expected to have value separated by new lines.
func WhitelistFileParse(filePath string) ([]string, error) {
	// Load file contents into memory
	data, err := utils.ReadFile(filePath)
	if err != nil {
		return []string{}, errors.Errorf("Failed to read whitelist file: %v", err)
	}

	// Convert the data to string, trim whitespace, and normalize new lines
	dataStr := strings.TrimSpace(string(normalizeNewlines(data)))

	// Return empty slice if the file is empty or only contains whitespace
	if dataStr == "" {
		return []string{}, nil
	}

	// Split the data at new lines and place in slice
	return strings.Split(dataStr, "\n"), nil
}

// WhitelistWatcher watches the specified whitelist file for changes. When
// changes occur, the file is parsed and the new whitelist is loaded into the
// whitelist map.
func (wl *Whitelist) WhitelistWatcher(updateFinished chan bool) {
	var err error

	wl.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		jww.ERROR.Printf("Failed to start new file watcher: %s", err)
	}
	defer func() {
		err = wl.WhitelistWatcherClose()
		if err != nil {
			jww.ERROR.Printf("Error running WhitelistWatcherClose(): %v", err)
		}
	}()

	done := make(chan bool)
	go func(updateFinished chan bool) {
		for {
			select {
			case event, ok := <-wl.watcher.Events:
				if !ok {
					return
				}
				jww.DEBUG.Printf("File watcher event: %v", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					jww.DEBUG.Printf("Watcher modified file: %v", event.Name)

					// Wait for write to end
					time.Sleep(5 * time.Second)

					// Update the whitelist from the new file
					err = wl.UpdateWhitelist()
					if err != nil {
						jww.ERROR.Printf("Error running WhitelistWatcherClose(): %v", err)
					}

					// Signify that the update is done
					if updateFinished != nil {
						updateFinished <- true
					}
				}

			case err1, ok := <-wl.watcher.Errors:
				if !ok {
					return
				}
				jww.DEBUG.Printf("File watcher error: %v", err1)
			}
		}
	}(updateFinished)

	// Add file watcher
	err = wl.watcher.Add(wl.file)
	if err != nil {
		jww.ERROR.Printf("Failed to add file watcher: %s", err)
	}

	<-done
}

// WhitelistWatcherClose closes the file watcher and closes the events channel.
func (wl *Whitelist) WhitelistWatcherClose() error {
	err := wl.watcher.Close()

	if err != nil {
		return errors.Errorf("Failed to close file watcher: %s", err)
	}

	return nil
}

// Exists searches if the specified key exists in the whitelist. Returns true if
// it exists and false otherwise.
func (wl *Whitelist) Exists(key string) bool {
	wl.RLock()
	defer wl.RUnlock()

	// Check if the key exists in the map
	_, ok := wl.list[key]

	return ok
}

// normalizeNewlines normalizes \r\n (Windows) and \r (Mac) into \n (UNIX).
func normalizeNewlines(d []byte) []byte {
	// Replace CR LF \r\n (Windows) with LF \n (UNIX)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)

	// Replace CF \r (Mac) with LF \n (UNIX)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)

	return d
}
