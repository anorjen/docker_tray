package window

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

var password string
var cancel bool

func GetPass() (string, bool) {

	if password == "" {
		gtk.Init(nil)

		win := newPassWindow()
		win.ShowAll()
		win.SetKeepAbove(true)

		gtk.Main()
	}

	return password, cancel
}

func newPassWindow() *gtk.Window {
	var err error

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle("Password")
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(200, 100)
	win.SetBorderWidth(10)
	win.SetResizable(false)

	win.Connect("destroy", func() {
		win.Destroy()
	})

	passwordLabel, err := gtk.LabelNew("Password: ")
	if err != nil {
		log.Fatal("Unable to create password label:", err)
	}
	passwordLabel.SetWidthChars(12)

	messageLabel, err := gtk.LabelNew("Wrong password")
	if err != nil {
		log.Fatal("Unable to create message label:", err)
	}
	messageLabel.SetWidthChars(5)
	messageLabel.SetNoShowAll(true)

	passwordEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}
	passwordEntry.SetVisibility(false)

	okButton, err := gtk.ButtonNewWithLabel("OK")
	if err != nil {
		log.Fatal("Unable to create OK button:", err)
	}
	cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		log.Fatal("Unable to create Cancel button:", err)
	}

	closeFunc := func() {
		win.Close()
		win.Destroy()
		gtk.MainQuit()
	}

	okActivateFunc := func() {
		password, err = passwordEntry.GetText()
		password = strings.TrimSpace(password)

		ok := checkPassword(password)
		if !ok {
			messageLabel.Show()
			passwordEntry.SetText("")
		} else {
			closeFunc()
		}
	}

	okButton.Connect("clicked", func() {
		okActivateFunc()
	})

	passwordEntry.Connect("activate", func() {
		okActivateFunc()
	})

	cancelButton.Connect("clicked", func() {
		log.Println("cancel")
		cancel = true
		closeFunc()
	})

	hBox1, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	hBox2, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)

	vBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	hBox1.PackStart(passwordLabel, true, false, 5)
	hBox1.PackStart(passwordEntry, true, false, 5)

	hBox2.PackStart(cancelButton, true, false, 5)
	hBox2.PackStart(okButton, true, false, 5)

	vBox.PackStart(hBox1, false, false, 5)
	vBox.PackStart(messageLabel, false, false, 0)
	vBox.PackStart(hBox2, false, false, 5)

	win.Add(vBox)

	return win
}

func checkPassword(password string) bool {
	cmd := exec.Command("sudo", "-lkS")
	cmd.Stdin = strings.NewReader(password)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	code := cmd.ProcessState.ExitCode()
	log.Println("Check password code: ", code)

	return code == 0
}
