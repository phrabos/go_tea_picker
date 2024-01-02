package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
BorderStyle(lipgloss.NormalBorder()).
BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { 
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
			case "esc":
				if m.table.Focused() {
					m.table.Blur()
				} else {
					m.table.Focus()
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			case "e":
				return m, tea.Batch(
					tea.Printf("You selected %s", m.table.SelectedRow()[1]),
				)
			case "enter":
				return m, tea.Batch(
					tea.Printf("You selected %s", m.table.SelectedRow()[1]),
				)
		}
	}
	m.table, _ = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

type TeaInventory struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Category string `json:"category"`
	Subcategory string `json:"subCategory"`
	Weight int `json:"weight"`
	Province string `json:"province"`
}

func main() {

	teaInventory, decodeErr := decodeJson("./teas.json")
	if decodeErr != nil {
		fmt.Printf("error decoding json %v\n", decodeErr)
		os.Exit(1)
	}
	columnConfig := getColumnConfig()
	columns, rows := getTableData(teaInventory, columnConfig)

	// make the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// make table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("200")).
		Background(lipgloss.Color("57")).
		Bold(false)

	// apply the styles to the table
	t.SetStyles(s)
	
	m := model{t}

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Printf("error running tea app %v\n", err)
		os.Exit(1)
	}
}

func decodeJson(filePath string) ([]TeaInventory, error) {

	teaJson, osErr := os.ReadFile(filePath)
	if osErr != nil {
		return nil, fmt.Errorf("error reading file %v", osErr)
	}

	var teaData []TeaInventory
	
	jsonTeaData := []byte(teaJson)
	isJsonValid := json.Valid(jsonTeaData)
	if !isJsonValid {
		return nil, errors.New("invalid json")
	}
	jsonErr := json.Unmarshal(jsonTeaData, &teaData)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to unmarshall json: %v", jsonErr)
	}
	return teaData, nil
}

func getTableData(teaInventoryItems []TeaInventory, columnConfig map[string]ColumnConfig) ([]table.Column, []table.Row) {
	// use the first item to generate columns
	columns := getTableColumns(teaInventoryItems[0], columnConfig)
	rows := make([]table.Row, len(teaInventoryItems))
	for i, v := range teaInventoryItems {
		stringValues := make([]string, len(columnConfig))
		for header, config := range columnConfig {
			sortOrder := config.SortOrder
			fieldValue := getFieldByHeader(v, header);
			stringValues[sortOrder] = fieldValue
			rows[i] = stringValues
		}	
				
	}
	return columns, rows
}
type ColumnConfig struct {
	Width int
	SortOrder int
}
func getColumnConfig() map[string]ColumnConfig{
	return map[string]ColumnConfig{
		"Id": { Width: 5, SortOrder: 0},
		"Category": { Width: 8, SortOrder: 2},
		"Subcategory": { Width: 12, SortOrder: 3},
		"Name": { Width: 20, SortOrder: 1},
		"Weight": { Width: 10, SortOrder: 5},
		"Province": { Width: 20, SortOrder: 4},
	}
}

func getTableColumns(teaInventoryItem TeaInventory, columnConfig map[string]ColumnConfig) []table.Column {
	var columns = make([]table.Column, len(columnConfig))
	for header, config := range columnConfig {
		width := config.Width
		sortOrder := config.SortOrder
		if header == "Weight" {
			columns[sortOrder] = table.Column{Title: header + "(g)", Width: width}
		} else {
			columns[sortOrder] = table.Column{Title: header, Width: width}
		} 
	}
	return columns
}	

func getFieldByHeader(t TeaInventory, header string) string {
	switch header {
	case "Id":
		return t.Id
	case "Name":
		return t.Name
	case "Category":
		return t.Category
	case "SubCategory":
		return t.Subcategory
	case "Weight":
		return strconv.Itoa(t.Weight)
	case "Province":
		return t.Province
	default:
		return ""
	}
}

	// rowsStub := []table.Row{
	// 	{"101", "Li Shan", "oolong", "ball rolled", "55"},
	// 	{"102", "Dong Ding", "oolong", "ball rolled", "140"},
	// 	{"103", "Shui Xian", "oolong", "yencha", "75"},
	// 	{"104", "Sencha", "green", "steamed", "9"},
	// 	{"105", "Bai Mudan", "white", "", "50"},
	// 	{"106", "941", "puerh", "sheng", "7"},
	// 	{"107", "Charlie", "puerh", "shou", "155"},
	// 	{"108", "Da Hong Pao", "oolong", "yencha", "65"},
	// }