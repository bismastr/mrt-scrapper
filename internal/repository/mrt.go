package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getScheduleById = `
SELECT
  	stations.id,
  	stations.name,
  	lines.id AS lines_id,
	lines.stations_id_start,
  	lines.stations_id_end,
  	end_stations.name AS stations_end_name,
	start_stations.name AS stations_start_name,
	schedules.time,
	schedules.is_holiday
FROM stations
LEFT JOIN lines ON stations.id = lines.stations_id_start
LEFT JOIN stations AS end_stations ON lines.stations_id_end = end_stations.id
LEFT JOIN stations AS start_stations ON lines.stations_id_start = start_stations.id
LEFT JOIN schedules ON lines.id = schedules.line_id
WHERE stations.id = $1 AND schedules.is_holiday = $2 AND lines.stations_id_end = $3
ORDER BY stations.id ASC, lines.id ASC, schedules.time ASC
`

func (q *Queries) GetScheduleById(ctx context.Context, id int64, isHoliday bool, directionStationId int64) ([]Schedule, error) {
	rows, err := q.db.Query(ctx, getScheduleById, id, isHoliday, directionStationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Schedule
	for rows.Next() {
		var i Schedule
		err := rows.Scan(&i.ID, &i.Name, &i.LinesID, &i.StationsIDStart, &i.StationsIDEnd, &i.StationsEndName, &i.StationsStartName, &i.Time, &i.IsHoliday)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

const getSchedule = `
SELECT
  	stations.id,
  	stations.name,
  	lines.id AS lines_id,
	lines.stations_id_start,
  	lines.stations_id_end,
  	end_stations.name AS stations_end_name,
	start_stations.name AS stations_start_name,
	schedules.time,
	schedules.is_holiday
FROM stations
LEFT JOIN lines ON stations.id = lines.stations_id_start
LEFT JOIN stations AS end_stations ON lines.stations_id_end = end_stations.id
LEFT JOIN stations AS start_stations ON lines.stations_id_start = start_stations.id
LEFT JOIN schedules ON lines.id = schedules.line_id
ORDER BY stations.id ASC, lines.id ASC, schedules.time ASC
`

func (q *Queries) GetSchedule(ctx context.Context) ([]Schedule, error) {
	rows, err := q.db.Query(ctx, getSchedule)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Schedule
	for rows.Next() {
		var i Schedule
		err := rows.Scan(&i.ID, &i.Name, &i.LinesID, &i.StationsIDStart, &i.StationsIDEnd, &i.StationsEndName, &i.StationsStartName, &i.Time, &i.IsHoliday)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const insertSchedule = `
INSERT INTO schedules (line_id, time, is_holiday)
VALUES ($1, $2, $3)
`

type InsertSchedule struct {
	LineID    int64
	Time      pgtype.Time
	IsHoliday pgtype.Bool
}

func (q *Queries) InsertSchedule(ctx context.Context, schedule InsertSchedule) error {
	_, err := q.db.Exec(ctx, insertSchedule, schedule.LineID, schedule.Time, schedule.IsHoliday)
	return err
}

const deleteSchedule = `
DELETE FROM schedules
WHERE line_id = $1
`

type DeleteSchedule struct {
	LineID int64
}

func (q *Queries) DeleteSchedule(ctx context.Context, schedule DeleteSchedule) error {
	_, err := q.db.Exec(ctx, deleteSchedule, schedule.LineID)
	return err
}

type Station struct {
	ID             int64
	StationName    string
	LaneID         int64
	StationStartID int64
	StationEndID   int64
}

const getLanes = `
SELECT 
	stations.id,
	stations.name,
	lines.id as lane_id,
	lines.stations_id_start,
	lines.stations_id_end
FROM stations
LEFT JOIN lines ON stations.id = lines.stations_id_start
ORDER BY stations.id
`

func (q *Queries) GetLanes(ctx context.Context) ([]Station, error) {
	rows, err := q.db.Query(ctx, getLanes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Station
	for rows.Next() {
		var i Station
		if err := rows.Scan(
			&i.ID,
			&i.StationName,
			&i.LaneID,
			&i.StationStartID,
			&i.StationEndID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
