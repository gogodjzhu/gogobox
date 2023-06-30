package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"gogobox/pkg/cmdutil/tui/tui_list"
	"gogobox/pkg/cmdutil/tui/tui_result"
	"gogobox/pkg/cmdutil/tui/tui_textinput"
	"os"
)

func main() {
	//TestResult()
	//TestTextInput()
	TestList()
}

func TestResult() {
	red := color.New(color.FgRed).SprintFunc()
	p := tea.NewProgram(tui_result.NewModel([]string{"a", "b", "c"}, red("please select..."), func(s string) {
		fmt.Println("selected:", s)
	}))

	// Run returns the model as a tea.Model.
	_, err := p.Run()
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}

func TestTextInput() {
	p := tea.NewProgram(tui_textinput.NewModel("Please input...", "placeholder", func(s string) {
		fmt.Println("You input:", s)
	}))
	if _, err := p.Run(); err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}

func TestList() {
	p := tea.NewProgram(tui_list.NewModel("", []tui_list.Option{
		tui_list.NewOption("apple", "apple is good\napple is very good"),
		tui_list.NewOption("banana", "banana is good"),
		tui_list.NewOption("orange", "orange is good"),
	}, []tui_list.CallbackFunc{
		{
			Keys: []string{"h"},
			Callback: func(option tui_list.Option) []tui_list.Option {
				fmt.Println("Press key:h, You select:", option.Title())
				return nil
			},
			ShortDescription: "hKey",
			FullDescription:  "h for a key full description",
		},
		{
			Keys: []string{"n"},
			Callback: func(option tui_list.Option) []tui_list.Option {
				fmt.Println("Press key:n, You select:", option.Title())
				return nil
			},
			ShortDescription: "nKey",
			FullDescription:  "n for a key full description",
		},
	}))
	if _, err := p.Run(); err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}
