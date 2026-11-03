// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package e2e

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"image"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgb/xtest"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

// AT-SPI state bits, see [AtspiStateType] and its definition in
// [atspi-constants.h].
//
// [AtspiStateType]: https://docs.gtk.org/atspi2/enum.StateType.html
// [atspi-constants.h]: https://gitlab.gnome.org/GNOME/at-spi2-core/-/blob/main/atspi/atspi-constants.h
const (
	stateChecked  = 4  // check box ticked
	stateEnabled  = 8  // control accepts user input
	stateSelected = 23 // tab or list item currently selected
	stateShowing  = 25 // on screen, taking hidden ancestors into account
)

// desktop is a headless X11 display with its own D-Bus session and AT-SPI
// accessibility bus, for driving and inspecting Qt applications and Firefox.
//
// Applications started with [desktop.env] render on the Xvfb display, and
// expose their UI elements via AT-SPI on the desktop's private accessibility
// bus. Specs then inspect and operate these elements using [accessible], and
// click on them via the XTest extension where AT-SPI actions don't suffice.
// Desktops are isolated from each other and the developer's session, with
// their own home, XDG and runtime directories, as well as extcap directory
// containing only the cshargextcap plugin under test.
type desktop struct {
	display    string     // X11 DISPLAY, such as ":42"
	xauthority string     // Xauthority file granting access to the display
	sessionBus string     // D-Bus session bus address
	runtimeDir string     // XDG_RUNTIME_DIR, also identifying desktop processes
	home       string     // HOME of applications on the desktop
	extcapDir  string     // WIRESHARK_EXTCAP_DIR with only cshargextcap
	session    *dbus.Conn // session bus connection
	a11y       *dbus.Conn // accessibility bus connection
	x          *xgb.Conn  // X11 connection for synthesizing input
}

// desktopMarker is an environment variable identifying the processes of a
// desktop, including those that URL handlers detach.
const desktopMarker = "CSHARGEXTCAP_E2E_DESKTOP"

// newDesktop starts a new desktop that gets torn down after the current spec,
// killing all applications started on it.
func newDesktop() *desktop {
	GinkgoHelper()
	d := &desktop{
		runtimeDir: tempDir(),
		home:       tempDir(),
		extcapDir:  tempDir(),
		xauthority: filepath.Join(tempDir(), "Xauthority"),
	}
	Expect(os.Symlink(extcapBin, filepath.Join(d.extcapDir, "cshargextcap"))).To(Succeed())

	// A wildcard Xauthority entry for any display, with a random cookie. Each
	// field is a big-endian length-prefixed byte string, following the family
	// (FamilyWild): address, display number, auth name, and auth data.
	cookie := make([]byte, 16)
	_, _ = rand.Read(cookie)
	var xauth bytes.Buffer
	_ = binary.Write(&xauth, binary.BigEndian, uint16(0xffff))
	for _, field := range [][]byte{nil, nil, []byte("MIT-MAGIC-COOKIE-1"), cookie} {
		_ = binary.Write(&xauth, binary.BigEndian, uint16(len(field)))
		xauth.Write(field)
	}
	Expect(os.WriteFile(d.xauthority, xauth.Bytes(), 0o600)).To(Succeed())

	// Xvfb picks a free display number and reports it on stdout.
	xvfb := Successful(gexec.Start(
		exec.Command("Xvfb", "-displayfd", "1", "-auth", d.xauthority,
			"-screen", "0", "1280x1024x24", "-nolisten", "tcp"),
		nil, GinkgoWriter))
	DeferCleanup(func() { xvfb.Terminate().Wait(5 * time.Second) })
	Eventually(xvfb.Out).Within(10 * time.Second).Should(gbytes.Say(`\d+\n`))
	displayNo := strings.TrimSpace(string(xvfb.Out.Contents()))
	d.display = ":" + displayNo

	cmd := exec.Command("dbus-daemon", "--session", "--nofork", "--print-address")
	// The activated AT-SPI bus launcher inherits this environment.
	cmd.Env = append(os.Environ(),
		"DISPLAY="+d.display, "XAUTHORITY="+d.xauthority, "XDG_RUNTIME_DIR="+d.runtimeDir)
	bus := Successful(gexec.Start(cmd, nil, GinkgoWriter))
	DeferCleanup(func() { bus.Terminate().Wait(5 * time.Second) })
	Eventually(bus.Out).Within(10 * time.Second).Should(gbytes.Say(`\n`))
	d.sessionBus = strings.TrimSpace(string(bus.Out.Contents()))

	d.session = Successful(dbus.Connect(d.sessionBus))
	DeferCleanup(d.session.Close)
	// Asking for the accessibility bus activates the AT-SPI bus launcher.
	var a11yAddr string
	Eventually(func() error {
		return d.session.Object("org.a11y.Bus", "/org/a11y/bus").
			Call("org.a11y.Bus.GetAddress", 0).Store(&a11yAddr)
	}).Within(10 * time.Second).ProbeEvery(250 * time.Millisecond).Should(Succeed())
	d.a11y = Successful(dbus.Connect(a11yAddr))
	DeferCleanup(d.a11y.Close)

	// Passing the cookie directly keeps xgb from logging via its global logger.
	d.x = Successful(xgb.NewConnNetWithCookieHex(
		Successful(net.Dial("unix", "/tmp/.X11-unix/X"+displayNo)), hex.EncodeToString(cookie)))
	DeferCleanup(d.x.Close)
	Expect(xtest.Init(d.x)).To(Succeed())

	DeferCleanup(d.killProcesses)
	return d
}

// env returns the environment variables for applications on the desktop, to be
// added to the test process environment.
func (d *desktop) env() []string {
	return []string{
		desktopMarker + "=" + d.runtimeDir,
		"DISPLAY=" + d.display,
		"XAUTHORITY=" + d.xauthority,
		"DBUS_SESSION_BUS_ADDRESS=" + d.sessionBus,
		"XDG_RUNTIME_DIR=" + d.runtimeDir,
		"HOME=" + d.home,
		"XDG_CONFIG_HOME=" + filepath.Join(d.home, ".config"),
		"XDG_DATA_HOME=" + filepath.Join(d.home, ".local", "share"),
		"WIRESHARK_EXTCAP_DIR=" + d.extcapDir,
		// Never forward URLs to a real browser, such as VS Code's host browser helper.
		"BROWSER=true",
		"DEFAULT_BROWSER=true",
		"QT_QPA_PLATFORM=xcb",
		"QT_LINUX_ACCESSIBILITY_ALWAYS_ON=1",
		"GNOME_ACCESSIBILITY=1",
	}
}

// killProcesses kills all processes started on the desktop, identified by the
// desktop marker in their environment. This also catches processes detached
// from their starters, such as Wireshark launched by Firefox via a URL handler,
// and the extcap plugin and dumpcap processes Wireshark starts in turn.
func (d *desktop) killProcesses() {
	marker := []byte("\x00" + desktopMarker + "=" + d.runtimeDir + "\x00")
	environs, _ := filepath.Glob("/proc/[0-9]*/environ")
	for _, environ := range environs {
		env, err := os.ReadFile(environ)
		if err != nil || !bytes.Contains(append([]byte{0}, env...), marker) {
			continue
		}
		if pid, err := strconv.Atoi(filepath.Base(filepath.Dir(environ))); err == nil {
			_ = syscall.Kill(pid, syscall.SIGKILL)
		}
	}
}

// root returns the root of all accessible applications on the desktop; its
// children are the applications.
func (d *desktop) root() accessible {
	return accessible{conn: d.a11y, accessibleRef: accessibleRef{
		Bus: "org.a11y.atspi.Registry", Path: "/org/a11y/atspi/accessible/root"}}
}

// app returns the named application on the desktop, without searching other
// applications' (potentially large) subtrees.
func (d *desktop) app(name string) (accessible, bool) {
	for _, app := range d.root().children() {
		if app.name() == name {
			return app, true
		}
	}
	return accessible{}, false
}

// clickOn clicks the left mouse button on the center of a, once a is showing
// and has stopped moving, as stray clicks might hit links. Use it where
// applications ignore AT-SPI actions or require a real user gesture, such as
// Wireshark's extcap gear and combo box popups, or links in web pages.
func (d *desktop) clickOn(a accessible) {
	GinkgoHelper()
	var rect image.Rectangle
	Eventually(func() image.Rectangle {
		prev := rect
		rect = image.Rectangle{}
		// Firefox reports -1 extents until it has cached them.
		if r := a.extents(); a.state(stateShowing) && r.Min.X >= 0 && r.Min.Y >= 0 {
			rect = r
		}
		if rect != prev {
			return image.Rectangle{}
		}
		return rect
	}).Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).ShouldNot(BeZero())
	at := image.Pt((rect.Min.X+rect.Max.X)/2, (rect.Min.Y+rect.Max.Y)/2)

	root := xproto.Setup(d.x).DefaultScreen(d.x).Root
	Expect(xtest.FakeInputChecked(d.x, xproto.MotionNotify, 0, 0, root,
		int16(at.X), int16(at.Y), 0).Check()).To(Succeed())
	for _, ev := range []byte{xproto.ButtonPress, xproto.ButtonRelease} {
		Expect(xtest.FakeInputChecked(d.x, ev, 1, 0, root, 0, 0, 0).Check()).To(Succeed())
	}
}

// accessibleRef is the D-Bus wire form of an AT-SPI object reference.
type accessibleRef struct {
	Bus  string
	Path dbus.ObjectPath
}

// accessible is a UI element as seen through AT-SPI. Queries on vanished
// elements return zero values, so that probes in Eventually simply retry.
// Roles and names differ between toolkits: for instance, Qt reports "push
// button", while Firefox reports "button".
type accessible struct {
	conn *dbus.Conn
	accessibleRef
}

// call calls the specified AT-SPI method, such as "Accessible.GetRoleName", on
// this UI element.
func (a accessible) call(method string, args ...any) *dbus.Call {
	return a.conn.Object(a.Bus, a.Path).Call("org.a11y.atspi."+method, 0, args...)
}

// wrap returns the element referenced by ref on the same accessibility bus.
func (a accessible) wrap(ref accessibleRef) accessible {
	return accessible{conn: a.conn, accessibleRef: ref}
}

// name returns a's accessible name, such as a label text or button caption.
func (a accessible) name() string {
	v, _ := a.conn.Object(a.Bus, a.Path).GetProperty("org.a11y.atspi.Accessible.Name")
	s, _ := v.Value().(string)
	return s
}

// role returns a's role name, such as "label" or "check box".
func (a accessible) role() string {
	var role string
	_ = a.call("Accessible.GetRoleName").Store(&role)
	return role
}

// children returns a's child elements, in UI order.
func (a accessible) children() []accessible {
	var refs []accessibleRef
	_ = a.call("Accessible.GetChildren").Store(&refs)
	children := make([]accessible, 0, len(refs))
	for _, ref := range refs {
		children = append(children, a.wrap(ref))
	}
	return children
}

// child returns a's child element at the specified index.
func (a accessible) child(index int32) accessible {
	var ref accessibleRef
	_ = a.call("Accessible.GetChildAtIndex", index).Store(&ref)
	return a.wrap(ref)
}

// parent returns a's parent element.
func (a accessible) parent() accessible {
	v, _ := a.conn.Object(a.Bus, a.Path).GetProperty("org.a11y.atspi.Accessible.Parent")
	var ref accessibleRef
	_ = dbus.Store([]any{v.Value()}, &ref)
	return a.wrap(ref)
}

// indexInParent returns a's index among its parent's children.
func (a accessible) indexInParent() int32 {
	var index int32
	_ = a.call("Accessible.GetIndexInParent").Store(&index)
	return index
}

// next returns the next sibling, such as the field following a label.
func (a accessible) next() accessible {
	return a.parent().child(a.indexInParent() + 1)
}

// state returns whether the specified AT-SPI state bit of a is set, see the
// state* constants.
func (a accessible) state(bit uint) bool {
	var states []uint32
	if a.call("Accessible.GetState").Store(&states) != nil || len(states) != 2 {
		return false
	}
	return (uint64(states[1])<<32|uint64(states[0]))&(1<<bit) != 0
}

// extents returns the screen area of a.
func (a accessible) extents() image.Rectangle {
	var ext struct{ X, Y, W, H int32 }
	_ = a.call("Component.GetExtents", uint32(0)).Store(&ext)
	return image.Rect(int(ext.X), int(ext.Y), int(ext.X+ext.W), int(ext.Y+ext.H))
}

// press triggers a's (first) action, such as pressing a button or tab.
func (a accessible) press() error {
	return a.call("Action.DoAction", int32(0)).Err
}

// setText replaces the text of a's editable text, such as a line edit.
func (a accessible) setText(text string) error {
	return a.call("EditableText.SetTextContents", text).Err
}

// uri returns the target of link a.
func (a accessible) uri() string {
	var uri string
	_ = a.call("Hyperlink.GetURI", int32(0)).Store(&uri)
	return uri
}

// cell returns the cell of table a at the specified row and column.
func (a accessible) cell(row, column int32) accessible {
	var ref accessibleRef
	_ = a.call("Table.GetAccessibleAt", row, column).Store(&ref)
	return a.wrap(ref)
}

// columns returns the number of columns of table a.
func (a accessible) columns() int32 {
	v, _ := a.conn.Object(a.Bus, a.Path).GetProperty("org.a11y.atspi.Table.NColumns")
	n, _ := v.Value().(int32)
	return n
}

// find returns the first element in a's subtree, a included, with the
// specified role and a name satisfying the matcher.
func (a accessible) find(role string, name types.GomegaMatcher) (accessible, bool) {
	if a.role() == role {
		if ok, _ := name.Match(a.name()); ok {
			return a, true
		}
	}
	for _, child := range a.children() {
		if found, ok := child.find(role, name); ok {
			return found, true
		}
	}
	return accessible{}, false
}

// mustFind is like find, but fails the current spec if there's no match.
func (a accessible) mustFind(role string, name types.GomegaMatcher) accessible {
	GinkgoHelper()
	found, ok := a.find(role, name)
	Expect(ok).To(BeTrue(), "no %s with name %s", role, name)
	return found
}

// findAll returns all elements in a's subtree, a included, with the specified
// role, in depth-first order.
func (a accessible) findAll(role string) []accessible {
	var found []accessible
	if a.role() == role {
		found = append(found, a)
	}
	for _, child := range a.children() {
		found = append(found, child.findAll(role)...)
	}
	return found
}

// names returns the names of the specified elements.
func names(as []accessible) []string {
	names := make([]string, 0, len(as))
	for _, a := range as {
		names = append(names, a.name())
	}
	return names
}
