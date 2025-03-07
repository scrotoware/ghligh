/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/scrotadamus/ghligh/document"
	"github.com/spf13/cobra"
)

type browsableDoc struct {
	doc         *document.GhlighDoc
	currentPage int

	title  string
	author string
	nPages int
}

func (d *browsableDoc) getCurrentPage() int {
	return d.currentPage
}
func (d *browsableDoc) setCurrentPage(i int) {
	switch {
	case i >= d.nPages:
		d.currentPage = d.nPages - 1
	case i < 0:
		d.currentPage = 0
	default:
		d.currentPage = i
	}
}

func openBrowsableDoc(path string) (*browsableDoc, error) {
	// open and init document
	doc, err := document.Open(path)
	if err != nil {
		return nil, err
	}

	b := &browsableDoc{
		doc:         doc,
		currentPage: 0,
		title:       doc.Info().Title,
		author:      doc.Info().Author,
		nPages:      doc.GetNPages(),
	}
	return b, nil
}

func (d *browsableDoc) header() string {
	return fmt.Sprintf("[yellow]%s - %s ~ Page %d of %d[-]",
		d.title,
		d.author,
		d.currentPage+1,
		d.nPages,
	)
}

func (d *browsableDoc) status(accumulator int) string {

	accumulated := ""
	if accumulator != 0 {
		accumulated = fmt.Sprintf("[%d] --", accumulator)
	}

	// TODO
	// return fmt.Sprintf("%s %d", accumulated, percetuale_cursore)
	return fmt.Sprintf("%s page %d of %d",
		accumulated,
		d.currentPage+1,
		d.nPages,
	)

}

func (d *browsableDoc) getCurrentPageText() string {
	text, _ := d.doc.GetPageText(d.currentPage)
	return text
}

type Browser struct {
	app         *tview.Application
	pageContent *tview.TextView
	header      *tview.TextView
	status      *tview.TextView
	layout      *tview.Flex

	docs    []*browsableDoc
	currDoc int

	accumulator int
}

func newBrowser(paths []string) *Browser {
	app := tview.NewApplication()
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	pageContent := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	status := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).     // TODO Fixed size header (1 line)
		AddItem(pageContent, 0, 1, true). // TODO Content takes all remaining space
		AddItem(status, 1, 0, false)      // TODO Fixed size footer (1 line)

	var docs []*browsableDoc
	for _, path := range paths {
		doc, err := openBrowsableDoc(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}
		docs = append(docs, doc)
	}

	return &Browser{
		app:         app,
		header:      header,
		pageContent: pageContent,
		status:      status,
		layout:      layout,
		docs:        docs,
	}

}

func (b *Browser) currentDoc() *browsableDoc {
	return b.docs[b.currDoc]
}

func (b *Browser) updateStatus() {
	b.status.SetText(b.currentDoc().status(b.accumulator))
}

func (b *Browser) updateHeader() {
	b.header.SetText(b.currentDoc().header())
}

func (b *Browser) updateContent() {
	text := b.currentDoc().getCurrentPageText()

	b.pageContent.
		SetTextAlign(tview.AlignLeft).
		SetText(text).
		SetBorder(true)
}

func (b *Browser) updatePage(newPage int) {
	// lock mutex here
	b.currentDoc().setCurrentPage(newPage)
	b.updateHeader()
	b.updateContent()
	// unlock mutex here

	b.updateStatus()
}
func (b *Browser) Run() {
	b.updatePage(0)
	b.app.SetInputCapture(b.handle)
	if err := b.app.SetRoot(b.layout, true).Run(); err != nil {
		panic(err)
	}
}

// if acc is equal to 0 return 1, return acc otherwise
func newAcc(acc int) int {
	if acc == 0 {
		return 1
	}
	return acc
}

func (b *Browser) nextDoc(acc int) {
	acc = newAcc(acc)
	// todo implement

}
func (b *Browser) prevDoc(acc int) {
	acc = newAcc(acc)
	// todo implement

}
func (b *Browser) nextPage(acc int) {
	b.resetAccumulator()
	acc = newAcc(acc)
	currentPage := b.currentDoc().getCurrentPage()
	b.updatePage(currentPage + acc)
}
func (b *Browser) prevPage(acc int) {
	b.resetAccumulator()
	acc = newAcc(acc)
	currentPage := b.currentDoc().getCurrentPage()
	b.updatePage(currentPage - acc)

}
func (b *Browser) scrollUp(acc int) {
	b.resetAccumulator()
	acc = newAcc(acc)
	b.pageContent.ScrollTo(acc, 0)
	b.updateStatus()

}
func (b *Browser) scrollDown(acc int) {
	b.resetAccumulator()
	acc = newAcc(acc)
	b.pageContent.ScrollTo(acc, 0)
	b.updateStatus()
}

func (b *Browser) resetAccumulator() {
	b.accumulator = 0
}
func (b *Browser) handle(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		b.accumulator = 0
		b.updateStatus()
	case tcell.KeyRune:
		switch event.Rune() {
		case 'Q':
			b.app.Stop()
			return event
			//TODO		case 'q': QUIT CURRENT DOCUMENT
			//			b.app.Stop()
			//			return event
		case 'N':
			b.nextDoc(b.accumulator)
		case 'P':
			b.prevDoc(b.accumulator)
		case 'n':
			b.nextPage(b.accumulator)
		case 'p':
			b.prevPage(b.accumulator)
		case 'j':
			b.scrollUp(b.accumulator)
		case 'k':
			b.scrollDown(b.accumulator)
		default:
			if event.Rune() >= '0' && event.Rune() <= '9' {
				b.accumulator = b.accumulator*10 + int(event.Rune()-'0')
				b.updateStatus()
			}
		}
	}
	return event
}

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse pdf file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var docs []*document.GhlighDoc
		for _, arg := range args {
			doc, err := document.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not open %s: %v\n", arg, err)
				continue
			}
			docs = append(docs, doc)
		}

		// if more then one file show menu
		browser := newBrowser(args)
		browser.Run()

	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
