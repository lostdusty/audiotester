package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

//go:embed icon.png
var windowIcon []byte

func main() {
	appl := app.NewWithID("link.princessmortix.audiotester")
	win := appl.NewWindow("AudioTester")
	win.SetMaster()
	win.CenterOnScreen()
	win.RequestFocus()
	win.Resize(fyne.Size{Width: 800, Height: 400})
	win.SetIcon(fyne.NewStaticResource("icon", windowIcon))

	titleTextWelcome := widget.NewRichTextFromMarkdown("# Headphones tester\nSmall application created to test headphones or other audio equipament before buying, in-store.")
	titleTextWelcome.Wrapping = fyne.TextWrapWord

	btnLeftAudio := widget.NewButtonWithIcon("Test left channel", theme.NewThemedResource(resourceLeftSvg), nil)
	btnLeftAudio.IconPlacement = widget.ButtonIconLeadingText
	btnLeftAudio.Importance = widget.HighImportance

	btnRightAudio := widget.NewButtonWithIcon("Test right channel", theme.NewThemedResource(resourceRightSvg), nil)
	btnRightAudio.IconPlacement = widget.ButtonIconTrailingText
	btnRightAudio.Importance = widget.HighImportance

	leftAndRightAudio := container.NewGridWithColumns(2, btnLeftAudio, btnRightAudio)
	separatorAudio := canvas.NewLine(color.RGBA{150, 150, 150, 255})
	separatorAudio.StrokeWidth = 2.5

	textMoreAudio := widget.NewLabel("You can also test the bass, mid-range and high-range with the audio below")
	textMoreAudio.Wrapping = fyne.TextWrapWord

	btnShortAudioTest := widget.NewButtonWithIcon("Short audio test", theme.NewThemedResource(resourceShortSvg), nil)

	btnLongAudioTest := widget.NewButtonWithIcon("Long audio test", theme.NewThemedResource(resourceLongSvg), nil)

	btnAbout := widget.NewButtonWithIcon("About", theme.InfoIcon(), func() {
		infoTextTitle := widget.NewRichTextFromMarkdown("# AudioTester")
		infoExit := widget.NewButtonWithIcon("", theme.WindowCloseIcon(), nil)
		infoExit.Importance = widget.DangerImportance
		infoHeader := container.NewBorder(nil, nil, infoTextTitle, infoExit)
		infoText := widget.NewRichTextFromMarkdown("Test or diagnose audio from headphones.\n\n## Author\nPrincess Mortix - princessmortix.link\n\n## License\nBSD-3\n\n## Version\n1.0")
		info := widget.NewModalPopUp(container.NewBorder(infoHeader, nil, nil, nil, infoText), win.Canvas())
		infoExit.OnTapped = func() { info.Hide() }
		info.Show()
	})
	btnCredits := widget.NewButtonWithIcon("Credits", theme.NewThemedResource(resourceCreditsSvg), func() {
		creditsTitle := widget.NewRichTextFromMarkdown("# Credits")
		creditsExit := widget.NewButtonWithIcon("", theme.WindowCloseIcon(), nil)
		creditsExit.Importance = widget.DangerImportance
		creditsHeader := container.NewBorder(nil, nil, creditsTitle, creditsExit)
		creditsText := widget.NewRichTextFromMarkdown("## Right and Left audio\nExtracted from Realtek Audio Console, edited in Filmora\n\n## Short audio\n[By akelley6 @ freesound.org](https://freesound.org/people/akelley6/sounds/486456/) \n(CC-BY-4)\n\n## Long audio\n[By deleted user4397472 @ freesound.org](https://freesound.org/people/deleted_user_4397472/sounds/342608/) \n(Creative Commons 0)")
		creditsModal := widget.NewModalPopUp(container.NewBorder(creditsHeader, nil, nil, nil, creditsText), win.Canvas())
		creditsExit.OnTapped = func() { creditsModal.Hide() }
		creditsModal.Show()
	})
	btnCredits.Importance = widget.SuccessImportance
	creditsAndAbout := container.NewGridWithColumns(2, btnCredits, btnAbout)
	separatorCreditsAbout := canvas.NewLine(color.RGBA{150, 150, 150, 255})
	separatorCreditsAbout.StrokeWidth = 2.5

	disableAll := func() {
		btnLeftAudio.Disable()
		btnRightAudio.Disable()
		btnShortAudioTest.Disable()
		btnLongAudioTest.Disable()
	}
	enableAll := func() {
		btnLeftAudio.Enable()
		btnRightAudio.Enable()
		btnShortAudioTest.Enable()
		btnLongAudioTest.Enable()
	}

	windowContent := container.NewVBox(titleTextWelcome, leftAndRightAudio, separatorAudio, btnShortAudioTest, btnLongAudioTest, separatorCreditsAbout, creditsAndAbout)
	centeredContent := container.NewCenter(windowContent)
	win.SetContent(centeredContent)
	btnLeftAudio.OnTapped = func() {
		disableAll()
		play(1)
		enableAll()
	}
	btnRightAudio.OnTapped = func() {
		disableAll()
		play(2)
		enableAll()
	}
	btnShortAudioTest.OnTapped = func() {
		disableAll()
		play(3)
		enableAll()
	}
	btnLongAudioTest.OnTapped = func() {
		disableAll()
		play(4)
		enableAll()
	}

	win.Show()
	appl.Run()
}

func play(song int) {
	var playAudio []byte
	switch song {
	case 1:
		playAudio = resourceAudioLeftMp3.StaticContent
	case 2:
		playAudio = resourceAudioRightMp3.StaticContent
	case 3:
		playAudio = resourceAudioPreShortMp3.StaticContent
	case 4:
		playAudio = resourceAudioPreLongMp3.StaticContent
	}

	reader := bytes.NewReader(playAudio)
	audioReader := io.NopCloser(reader)
	streamer, format, err := mp3.Decode(audioReader)
	if err != nil {
		dialog.ShowInformation("Oh no, an error happened!", "This is a rare message,\nplease report this error to the developer\non https://github.com/princessmortix/audiotester", fyne.CurrentApp().Driver().AllWindows()[0])
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
		fmt.Println("Done playing")
	})))
	<-done
}
