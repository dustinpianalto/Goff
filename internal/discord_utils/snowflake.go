package discord_utils

import "time"

type Snowflake struct {
	CreationTime time.Time
	WorkerID     int8
	ProcessID    int8
	Increment    int16
}

func ParseSnowflake(s int64) Snowflake {
	const (
		DISCORD_EPOCH   = 1420070400000
		TIME_BITS_LOC   = 22
		WORKER_ID_LOC   = 17
		WORKER_ID_MASK  = 0x3E0000
		PROCESS_ID_LOC  = 12
		PROCESS_ID_MASK = 0x1F000
		INCREMENT_MASK  = 0xFFF
	)
	t := (s >> TIME_BITS_LOC) + DISCORD_EPOCH
	sec := t / 1000.0
	nano := t % 1000.0 * 1000000.0
	creationTime := time.Unix(sec, nano)
	workerID := (s & WORKER_ID_MASK) >> WORKER_ID_LOC
	processID := (s & PROCESS_ID_MASK) >> PROCESS_ID_LOC
	increment := s & INCREMENT_MASK
	return Snowflake{
		CreationTime: creationTime,
		WorkerID:     int8(workerID),
		ProcessID:    int8(processID),
		Increment:    int16(increment),
	}
}
