package commands

import (
	"fmt"
	"strconv"

	"treds/resp"
	"treds/store"
)

const ZRangeCommand = "ZRANGE"

func RegisterZRangeCommand(r CommandRegistry) {
	r.Add(&CommandRegistration{
		Name:     ZRangeCommand,
		Validate: validateZRange(),
		Execute:  executeZRangeCommand(),
		IsWrite:  false,
	})
}

func validateZRange() ValidationHook {
	return func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("expected 3 or multiple of 2 arguments, got %d", len(args))
		}
		_, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid start index: %s", args[0])
		}
		_, err = strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid end index: %s", args[1])
		}
		return nil
	}
}

func executeZRangeCommand() ExecutionHook {
	return func(args []string, store store.Store) string {
		startIndex, _ := strconv.Atoi(args[1])
		endIndex, _ := strconv.Atoi(args[2])
		withScore := true
		if len(args) > 3 {
			includeScore, err := strconv.ParseBool(args[3])
			if err != nil {
				return err.Error()
			}
			withScore = includeScore
		}
		v, err := store.ZRange(args[0], startIndex, endIndex, withScore)
		if err != nil {
			return resp.EncodeError(err.Error())
		}
		return resp.EncodeStringArray(v)
	}
}
