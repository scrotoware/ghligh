/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/scrotadamus/ghligh/document"
	"github.com/spf13/cobra"
)

var recursive bool

type resolver struct {
	paths   []string
	recurse bool
	ctx     context.Context
	ch      chan<- string
	wg      sync.WaitGroup
}

func (r *resolver) resolve() {
	for _, path := range r.paths {
		r.wg.Add(1)
		go r.resolvePath(path)
	}

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()
}

func (r *resolver) resolvePath(path string) {
	defer r.wg.Done()
	if err := r.ctx.Err(); err != nil {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		str, err := filepath.Abs(path)
		if err != nil {
			return
		}

		r.ch <- str
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving info for %s: %v\n", entry.Name(), err)
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		if info.IsDir() {
			if r.recurse {
				if err := r.ctx.Err(); err != nil {
					return
				}
				r.wg.Add(1)
				go r.resolvePath(fullPath)
			}
		} else if info.Mode().IsRegular() {
			r.ch <- fullPath
		}
	}
}

func lsArgs(args []string) []string {
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
			return nil
		}
		return []string{cwd}

	}
	return args
}

func checkFile(path string) bool {
	doc, err := document.Open(path)
	if err != nil {
		return false
	}
	defer doc.Close()

	return doc.HasHighlights()
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "show files with highlights or tagged with 'ls' [unix]",
	Long: `
	ghligh ls file1.pdf directory  [-R] [-c]

	will show every file inside directory that contains highlights or it is marked
	ls with the ghligh tag add command

	ghligh ls # show files in current dir
	ghligh ls file1.pdf # if it outpus file1.pdf it means that file1.pdf contains highlights
	ghligh ls -c file1.pdf # same as ghligh ls file1.pdf but exit status will not be if file1.pdf doesnt
				contains highlights

	ghligh ls -R # do it recursively, be careful with symlink dir cycles, as I am to lazy to
			address that particular issue
`,
	Run: func(cmd *cobra.Command, args []string) {
		files := lsArgs(args)
		ch := make(chan string)
		ctx := context.Background()

		res := resolver{
			paths:   files,
			recurse: recursive,
			ctx:     ctx,
			ch:      ch,
		}

		go res.resolve()

		var wg sync.WaitGroup
		var found bool
		for file := range ch {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				if checkFile(f) {
					found = true
					fmt.Printf("%s\n", f)
				}
			}(file)
		}
		wg.Wait()

		check, err := cmd.Flags().GetBool("check")
		if err != nil {
			cmd.Help()
			return
		}
		if check && !found {
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&recursive, "recursive", "R", false, "List recursively")
	lsCmd.Flags().BoolP("check", "c", false, "exit status is 1 if no file its found")
	// order pdf by time of something (modification / creation) ???
	//lsCmd.Flags().BoolP("time", "t", false, "ls by time")
}
