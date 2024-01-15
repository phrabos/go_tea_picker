package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var baseStyle = lipgloss.NewStyle().
BorderStyle(lipgloss.NormalBorder()).
BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
	selectedTea string
	spinner  spinner.Model
	loading bool
	quitting bool
	teaInventory TeaInventory
}

func (m model) Init() tea.Cmd { 
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					m.loading = true
					return m, tea.Quit
				case "e":
					m.selectedTea = m.table.SelectedRow()[1]
					return m, nil 
				case "enter":
					return m, tea.Batch(
						tea.Printf("You selected %s", m.table.SelectedRow()[1]),
					)
				}
		default:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
	}
	m.table, _ = m.table.Update(msg)
	return m, nil
}

func (m model) View() string {
	str := fmt.Sprintf("\n\n   %s Loading Tea Inventory\n\n", m.spinner.View())
	if m.loading {
		return str
	}
	if m.quitting {
		return str + "\n"
	}
	if !m.loading && m.selectedTea == "" {
		return baseStyle.Render(m.table.View()) + "\n"
		
	}
		return fmt.Sprintf("You selected %v", m.selectedTea)
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, loading: true}
}

type TeaInventory struct {
	Id string `bson:"_id"`
	Name string `bson:"name"`
	Category string `bson:"category"`
	Subcategory string `bson:"subCategory"`
	Weight int `bson:"weight"`
	Province string `bson:"province"`
}

func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(viper.GetString("MONGO_URI")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
  client, mongoErr := mongo.Connect(context.TODO(), opts)
  if mongoErr != nil {
    panic(mongoErr)
  }
  defer func() {
    if mongoErr = client.Disconnect(context.TODO()); mongoErr != nil {
      panic(mongoErr)
    }
  }()
  // Send a ping to confirm a successful connection
  if connectionErr := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); connectionErr != nil {
    panic(connectionErr)
  }
  // fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	collection := client.Database("TeaCo").Collection("Inventory")
	cursor, collectionErr := collection.Find(context.Background(), bson.D{{}})
	if collectionErr != nil {
		log.Fatal(collectionErr)
	}

	var teaInventorySlice []TeaInventory

	for cursor.Next(context.Background()) {
		var teaInventoryItem TeaInventory
		decodeErr := cursor.Decode(&teaInventoryItem)
		if decodeErr != nil {
			log.Fatal(decodeErr)
		} 
		teaInventorySlice = append(teaInventorySlice, teaInventoryItem)
	}

	// teaInventory, decodeErr := decodeJson("./teas.json")
	// if decodeErr != nil {
	// 	fmt.Printf("error decoding json %v\n", decodeErr)
	// 	os.Exit(1)
	// }
	columnConfig := getColumnConfig()
	columns, rows := getTableData(teaInventorySlice, columnConfig)

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
	
	// m := model{table: t, selectedTea: ""}

	_, err := tea.NewProgram(initialModel()).Run()
	if err != nil {
		fmt.Printf("error running tea app %v\n", err)
		os.Exit(1)
	}
}

// func decodeJson(filePath string) ([]TeaInventory, error) {

// 	teaJson, osErr := os.ReadFile(filePath)
// 	if osErr != nil {
// 		return nil, fmt.Errorf("error reading file %v", osErr)
// 	}

// 	var teaData []TeaInventory
	
// 	jsonTeaData := []byte(teaJson)
// 	isJsonValid := json.Valid(jsonTeaData)
// 	if !isJsonValid {
// 		return nil, errors.New("invalid json")
// 	}
// 	jsonErr := json.Unmarshal(jsonTeaData, &teaData)
// 	if jsonErr != nil {
// 		return nil, fmt.Errorf("failed to unmarshall json: %v", jsonErr)
// 	}
// 	return teaData, nil
// }

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
