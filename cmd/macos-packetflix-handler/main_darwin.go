// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// Apple event four-char codes for "open URL" events.
const (
	kInternetEventClass = 'G'<<24 | 'U'<<16 | 'R'<<8 | 'L'
	kAEGetURL           = 'G'<<24 | 'U'<<16 | 'R'<<8 | 'L'
	keyDirectObject     = '-'<<24 | '-'<<16 | '-'<<8 | '-'
)

// idleTimeout is how long to wait for (further) URLs before quitting.
const idleTimeout = 5 * time.Second

// Cocoa must run on the main thread.
func init() { runtime.LockOSThread() }

func main() {
	openLog()
	log.Printf("started, pid %d", os.Getpid())
	if _, err := purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit",
		purego.RTLD_GLOBAL|purego.RTLD_NOW); err != nil {
		alert(err.Error())
		os.Exit(1)
	}
	idle := time.AfterFunc(idleTimeout, func() {
		log.Print("idle, quitting")
		os.Exit(0)
	})

	sel := objc.RegisterName
	handleGetURL := sel("handleGetURLEvent:withReplyEvent:")
	handlerClass, err := objc.RegisterClass("PacketflixURLHandler", objc.GetClass("NSObject"), nil, nil,
		[]objc.MethodDef{{
			Cmd: handleGetURL,
			Fn: func(_ objc.ID, _ objc.SEL, event objc.ID, _ objc.ID) {
				idle.Stop()
				defer idle.Reset(idleTimeout)
				nsURL := event.Send(sel("paramDescriptorForKeyword:"), uint32(keyDirectObject)).
					Send(sel("stringValue"))
				if nsURL == 0 {
					log.Print("URL event without URL")
					return
				}
				pfURL := objc.Send[string](nsURL, sel("UTF8String"))
				log.Printf("received %s", pfURL)
				if err := startWireshark(pfURL); err != nil {
					alert(err.Error())
				}
			},
		}})
	if err != nil {
		alert(err.Error())
		os.Exit(1)
	}

	app := objc.ID(objc.GetClass("NSApplication")).Send(sel("sharedApplication"))
	// Must be registered before the app finishes launching, as otherwise the
	// URL event that launched us gets lost.
	objc.ID(objc.GetClass("NSAppleEventManager")).Send(sel("sharedAppleEventManager")).
		Send(sel("setEventHandler:andSelector:forEventClass:andEventID:"),
			objc.ID(handlerClass).Send(sel("new")), handleGetURL,
			uint32(kInternetEventClass), uint32(kAEGetURL))
	log.Print("waiting for URL events")
	app.Send(sel("run"))
}

// openLog logs to ~/Library/Logs/packetflix-handler.log, overwriting the log of
// the previous launch; this includes Go runtime crash output on stderr.
func openLog() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	f, err := os.Create(filepath.Join(home, "Library", "Logs", "packetflix-handler.log"))
	if err != nil {
		return
	}
	_ = syscall.Dup2(int(f.Fd()), 2)
	log.SetOutput(f)
}

// startWireshark starts a new Wireshark instance capturing from the specified
// packetflix URL. No shell is involved, so the URL cannot inject commands.
func startWireshark(pfURL string) error {
	args, err := openArgs(pfURL)
	if err != nil {
		return err
	}
	log.Printf("running open %q", args)
	if out, err := exec.Command("/usr/bin/open", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("cannot start Wireshark: %w\n%s", err, bytes.TrimSpace(out))
	}
	return nil
}

// alert shows the message in a modal dialog, passing it as an argument instead
// of as part of the script in order to avoid script injection.
func alert(msg string) {
	log.Printf("error: %s", msg)
	_ = exec.Command("/usr/bin/osascript",
		"-e", "on run argv",
		"-e", `display alert "Edgeshark packetflix:// Handler" message (item 1 of argv) as critical`,
		"-e", "end run",
		msg).Run()
}
