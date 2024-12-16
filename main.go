package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/tjarratt/babble"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyState int

const (
	KeyUnused KeyState = iota
	KeyWrong
	KeyPresent
	KeyCorrect
)

type State struct {
	word       string
	guess      string
	message    string
	display    []string
	correct    map[rune]bool
	guessCount int
	exit       bool
	guessBoard [6]string
	keyboard   map[rune]KeyState

	input textinput.Model

	flexBox *flexbox.FlexBox
}

func initialState() State {
	ti := textinput.New()
	ti.Placeholder = "guess"
	ti.Width = 5
	ti.CharLimit = 5
	ti.Focus()
	ti.Validate = func(s string) error {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return fmt.Errorf("letters only")
			}
		}
		return nil
	}

	s := State{
		correct:  make(map[rune]bool),
		keyboard: make(map[rune]KeyState),
		input:    ti,
		flexBox:  flexbox.New(0, 0),
	}

	// Create two cells for board and keyboard
	row := s.flexBox.NewRow().AddCells(
		flexbox.NewCell(1, 1),
		flexbox.NewCell(5, 1),
	)
	s.flexBox.AddRows([]*flexbox.Row{row})

	babbler := babble.NewBabbler()
	babbler.Count = 1

	// damn this is funny xd
	// have to do this couse babble doesn't let you return specific length
	// and i don't feel like making my own version
gotWord:
	s.word = babbler.Babble()
	if len(s.word) != 5 {
		goto gotWord
	}
	s.word = strings.ToLower(s.word)

	for _, r := range s.word {
		s.correct[r] = true
	}

	s.keyboard = make(map[rune]KeyState)
	for r := 'a'; r <= 'z'; r++ {
		s.keyboard[r] = KeyUnused
	}

	s.display = make([]string, 5)

	return s
}

func (s State) Init() tea.Cmd {
	return textinput.Blink
}

func (s State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var currentGuess strings.Builder

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.flexBox.SetWidth(msg.Width)
		s.flexBox.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return s, tea.Quit
		case tea.KeyEnter:
			if len(s.input.Value()) == 5 {
				s.guess = strings.ToLower(s.input.Value())

				if s.guessCount == 5 {
					s.message = "You clearly don't play scrabble, you lost!\nThe word was: " + s.word
					s.exit = true
					// goto display
					return s, tea.Quit
				}

				if s.guess == s.word {
					s.message = "Good job, you opened the dictionary, you won!\nThe word was: " + s.word
					s.exit = true
					// goto display
					return s, tea.Quit
				}

				// damn this is ugly but does the job
				for i := range s.guess {
					if s.correct[rune(s.guess[i])] {
						if s.guess[i] == s.word[i] {
							s.display[i] = rightPlace.Render(string(s.guess[i]))
							s.keyboard[rune(s.guess[i])] = KeyCorrect
						} else {
							s.display[i] = containsLetter.Render(string(s.guess[i]))
							// so we don't reset correct to just present
							if s.keyboard[rune(s.guess[i])] != KeyCorrect {
								s.keyboard[rune(s.guess[i])] = KeyPresent
							}
						}
					} else {
						s.display[i] = wrong.Render(string(s.guess[i]))
						s.keyboard[rune(s.guess[i])] = KeyWrong
					}
				}

				for _, d := range s.display {
					currentGuess.WriteString(d)
				}
				s.guessBoard[s.guessCount] = currentGuess.String()

				s.display = make([]string, 5)

				s.input.Reset()
				s.guessCount++
			}
		}
	}

	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

func (s State) View() string {
	if s.message != "" {
		return s.message + "\n"
	}

	var boardView strings.Builder
	var keyboardView strings.Builder

	// First add all previous guesses
	for i := 0; i < s.guessCount; i++ {
		boardView.WriteString(s.guessBoard[i] + "\n")
	}

	// Then build and add current guess display
	var currentGuess strings.Builder
	for _, d := range s.display {
		currentGuess.WriteString(d)
	}

	// Add current guess to view
	boardView.WriteString(currentGuess.String())

	// Save current guess for next round
	if len(s.display) > 0 {
		s.guessBoard[s.guessCount] = currentGuess.String()
	}

	// keyboard builder
	rows := []string{
		"QWERTYUIOP",
		"ASDFGHJKL",
		"ZXCVBNM",
	}

	for _, row := range rows {
		// Add proper spacing for keyboard layout
		if row == "ASDFGHJKL" {
			keyboardView.WriteString(" ")
		} else if row == "ZXCVBNM" {
			keyboardView.WriteString("  ")
		}

		for _, r := range row {
			switch s.keyboard[unicode.ToLower(r)] {
			case KeyUnused:
				keyboardView.WriteString(string(r) + " ")
			case KeyWrong:
				keyboardView.WriteString(wrong.Render(string(r)) + " ")
			case KeyPresent:
				keyboardView.WriteString(containsLetter.Render(string(r)) + " ")
			case KeyCorrect:
				keyboardView.WriteString(rightPlace.Render(string(r)) + " ")
			}
		}
		keyboardView.WriteString("\n")
	}

	// Add input view to board content
	// boardView.WriteString("\n" + s.input.View())
	board := lipgloss.JoinVertical(
		lipgloss.Left,
		boardView.String(),
		lipgloss.NewStyle().PaddingTop(6-s.guessCount).Render(s.input.View()),
	)

	// Set content for both cells
	s.flexBox.GetRow(0).GetCell(0).SetContent(board)
	s.flexBox.GetRow(0).GetCell(1).SetContent(keyboardView.String())

	// out := lipgloss.JoinVertical(
	// 	lipgloss.Left,
	// 	s.flexBox.Render(),
	// 	s.input.View(),
	// )

	return s.flexBox.Render()
}

func main() {
	p := tea.NewProgram(initialState(), tea.WithFPS(60))
	if _, err := p.Run(); err != nil {
		fmt.Printf("We fucked up, error: %v", err)
		os.Exit(1)
	}
	// fmt.Println(word)
	// loop()
}
