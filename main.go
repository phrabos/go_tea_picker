package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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
		case "enter":
			return m, tea.Batch(
				tea.Printf("You selected %s", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
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
	Weight string `json:"weight"`
	Province string `json:"province"`
}

func main() {

	teaInventory := decodeJson("./teas.json")
	columnConfig := getColumnConfig()
	columns, rows := getTableData(teaInventory, columnConfig)

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

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

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
	t.SetStyles(s)
	
	m := model{t}

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running tea app", err)
		os.Exit(1)
	}
}

	func decodeJson(filePath string) []TeaInventory {

		teaJson, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		var teaData []TeaInventory
		
		jsonTeaData := []byte(teaJson)
		isJsonValid := json.Valid(jsonTeaData)
		
		if isJsonValid {
			json.Unmarshal(jsonTeaData, &teaData)
			return teaData
		} else {
			panic("INVALID JSON")
		}
	}

	// func getTableColumns(teaInventoryItem TeaInventory, columnConfig map[string]ColumnConfig) []table.Column {
	// 	values := reflect.ValueOf(teaInventoryItem)
	// 	types := values.Type()
	// 	var columns = make([]table.Column, values.NumField())
	// 	for i := 0; i < values.NumField(); i++ {
	// 		header := types.Field(i).Name
	// 		config := columnConfig[header]
	// 		width := config.Width
	// 		sortOrder := config.SortOrder
	// 		columns[sortOrder] = table.Column{Title: header, Width: width}
	// 	}
	// 	return columns

	// }	
	func getTableData(teaInventoryItems []TeaInventory, columnConfig map[string]ColumnConfig) ([]table.Column, []table.Row) {
		var columns []table.Column
		rows := make([]table.Row, len(teaInventoryItems))
		for i, v := range teaInventoryItems {
			values := reflect.ValueOf(v)
			stringValues := make([]string, values.NumField())
			// use the first item to generate columns
			if i == 0 {
				columns = make([]table.Column, values.NumField())
				types := values.Type()
				for j := 0; j < values.NumField(); j++ {
					header := types.Field(j).Name
					config := columnConfig[header]
					width := config.Width
					sortOrder := config.SortOrder
					if header == "Weight" {
						columns[sortOrder] = table.Column{Title: header + "(g)", Width: width}
					} else {
						columns[sortOrder] = table.Column{Title: header, Width: width}
					} 
				}
				} else {
				types := values.Type()
				for k := 0; k < values.NumField(); k++ {
					header := types.Field(k).Name
					config := columnConfig[header]
					sortOrder := config.SortOrder
					fieldValue := values.Field(k);
					stringValues[sortOrder] = fieldValue.String()
					rows[i] = stringValues
				}
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
