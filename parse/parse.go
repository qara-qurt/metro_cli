package parse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"log"
	"strings"
	"time"
)

type Schedule struct {
	Station          string // current station name
	ArrivalTimeFromA string // arrival time from a
	ArrivalTimeFromB string // arrival time from b
}

type ScheduleStation struct {
	PrevStation     string // previous station name
	PrevStationTime string // previous station times
	NextStation     string // next station name
	NextStationTime string // next station times
}

func GetAllStateSchedule() []Schedule {
	var scheduleData []Schedule

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"http://metroalmaty.kz/?q=ru/schedule-list"},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			body := r.HTMLDoc.Find("tbody")
			body.Find("tr").Each(func(i int, s *goquery.Selection) {
				// Information
				elem := s.Find("td")
				// Station
				station := elem.Find("a").Text()
				// Arrival time from Райымбек
				time1 := s.Find("td:nth-child(2)").Text()
				if time1 == "" {
					time1 = "🏁"
				}
				// Arrival time from Б.Момышулы
				time2 := s.Find("td:nth-child(3)").Text()
				if time2 == "" {
					time2 = "🏁"
				}

				scheduleEntry := Schedule{
					Station:          station,
					ArrivalTimeFromA: color.GreenString(" ˅ " + time1),
					ArrivalTimeFromB: color.GreenString(" ˄ " + time2),
				}
				scheduleData = append(scheduleData, scheduleEntry)
			})
		},
		LogDisabled: true,
	}).Start()

	return scheduleData
}

func GetStatScheduleByName(stationName string, station int, all bool) []ScheduleStation {
	var scheduleData []ScheduleStation
	url := fmt.Sprintf("http://metroalmaty.kz/?q=ru/schedule-list-view/%v", station)

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			table := r.HTMLDoc.Find("table")
			// previous station name
			prevStation := table.Find("th:nth-child(1)").Text()
			// next station name
			nextStation := table.Find("th:nth-child(2)").Text()

			body := r.HTMLDoc.Find("tbody")
			body.Find("tr").Each(func(i int, s *goquery.Selection) {
				// previous station time
				prevStationTime := s.Find("td:nth-child(1)").Text()
				// next station time
				nextStationTime := s.Find("td:nth-child(2)").Text()

				var scheduleEntry ScheduleStation
				// all schedule
				if all {
					coloredPrevTime := getColorString(prevStationTime)
					coloredNextTime := getColorString(nextStationTime)

					scheduleEntry = ScheduleStation{
						PrevStation:     prevStation,
						NextStation:     nextStation,
						PrevStationTime: coloredPrevTime,
						NextStationTime: coloredNextTime,
					}
					scheduleData = append(scheduleData, scheduleEntry)
				} else {
					// current hour
					currentHour := fmt.Sprintf("%02d", time.Now().Hour())
					// next hour +1
					nextHour := fmt.Sprintf("%02d", time.Now().Hour()+1)
					if strings.HasPrefix(prevStationTime, currentHour) || strings.HasPrefix(prevStationTime, nextHour) {

						coloredPrevTime := getColorString(prevStationTime)
						coloredNextTime := getColorString(nextStationTime)

						scheduleEntry = ScheduleStation{
							PrevStation:     prevStation,
							NextStation:     nextStation,
							PrevStationTime: coloredPrevTime,
							NextStationTime: coloredNextTime,
						}
						scheduleData = append(scheduleData, scheduleEntry)
					}
				}

			})
		},
		LogDisabled: true,
	}).Start()

	return scheduleData
}

// miss the train
func isMissedTime(t string) bool {
	if strings.HasPrefix(t, "24") || t == "" {
		return false
	}
	now, err := time.Parse("15:04:05", time.Now().Format("15:04:05"))
	prevTime, err := time.Parse("15:04:05", t)
	if err != nil {
		log.Fatal(fmt.Sprintf("invalid time %v", err))
	}

	if now.Before(prevTime) {
		return false
	}
	return true
}

// getColorString returns a colored time string
func getColorString(stationTime string) string {
	isMissed := isMissedTime(stationTime)
	if isMissed {
		return color.RedString(" ˅ " + stationTime)
	}
	return color.GreenString(" ˄ " + stationTime)
}
