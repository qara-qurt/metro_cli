package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"log"
	"metroCLI/consts"
	"metroCLI/parse"
	"os"
	"time"
)

func main() {
	log.Println(consts.Welcome)
	cliRun()
}

func cliRun() {
	app := &cli.App{
		Name:  "METRO CLI",
		Usage: "Расписание метро в Алматы",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "station",
				Aliases: []string{"s", "stat"},
				Value:   0,
				Usage:   "Указать станцию для расписание",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Показать полную таблицу распиание станций",
			},
		},
		Action: func(context *cli.Context) error {
			station := context.Int("station")
			all := context.Bool("all")

			if all && station == 0 {
				fmt.Println("Ошибка: Флаг --all можно использовать только с --station | -s ")
				return nil
			}

			for _, s := range consts.StationNames {
				stationName := consts.StationNames[station]
				if s == stationName {
					if all {
						return getScheduleByStation(stationName, station, all)
					}
					return getScheduleByStation(stationName, station, all)
				}
			}
			return getSchedule()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getSchedule() error {
	//schedule data
	data := parse.GetAllStateSchedule()
	if len(data) == 0 {
		fmt.Println("Расписание для станции не найдено")
		return nil
	}
	renderSchedule(data)
	return nil
}

func getScheduleByStation(stationName string, station int, all bool) error {
	data := parse.GetStatScheduleByName(stationName, station, all)
	if len(data) == 0 {
		fmt.Printf("Расписание для станции %s не найдено\n", stationName)
		return nil
	}
	renderStationSchedule(stationName, data)
	return nil
}

// render common schedule
func renderSchedule(data []parse.Schedule) {
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Станция", "Время прибытия(Райымбек.б - Б.Момышулы)", "Время прибытия(Б.Момышулы - Райымбек.Б)"}
	table.SetHeader(header)
	// add style to table
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoWrapText(false)
	// Add data to the table
	for i, entry := range data {
		name := fmt.Sprintf("%v. %s", i+1, entry.Station)
		table.Append([]string{name, entry.ArrivalTimeFromA, entry.ArrivalTimeFromB})
	}

	// Render the table
	table.Render()
}

// render station schedule
func renderStationSchedule(stationName string, data []parse.ScheduleStation) {
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)

	now := time.Now()
	stationHeader := []string{stationName, now.Format(time.TimeOnly)}
	stationsHeader := []string{data[0].PrevStation, data[0].NextStation}

	// add style to table
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoWrapText(false)

	table.Append(stationHeader)
	table.Append(stationsHeader)
	// Add data to the table
	for _, entry := range data {
		table.Append([]string{entry.PrevStationTime, entry.NextStationTime})
	}

	// Render the table
	table.Render()
}
