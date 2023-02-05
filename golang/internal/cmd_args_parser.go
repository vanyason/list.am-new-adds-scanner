package internal

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type CmdArguments struct {
	Price        uint `short:"p" long:"price" description:"price in drams" required:"true"`
	Rooms        uint `short:"r" long:"rooms" description:"amount of rooms. 0 means do not care" default:"0"`
	ErrorCounter uint `short:"e" long:"errorcounter" description:"amount of errors to stop execution. 0 means never" default:"15"`
	LoopPause    uint `short:"t" long:"looppause" description:"time in minutes to wait before next loop. 0 means never" default:"15"`
	LogFileName  string
	DBFileName   string
}

func ParseCmdLineArgs() (CmdArguments, error) {
	var args CmdArguments
	if _, err := flags.NewParser(&args, flags.HelpFlag|flags.PassDoubleDash).Parse(); err != nil {
		return args, err
	}

	const maxPrice uint = 10000000
	const maxRooms uint = 20
	const maxErrors uint = 50
	const maxPause uint = 3 * 60

	if args.Price > maxPrice {
		return CmdArguments{}, fmt.Errorf("invalid price: %d. max: %d", args.Price, maxPrice)
	}
	if args.Rooms > maxRooms {
		return CmdArguments{}, fmt.Errorf("invalid rooms amount counter: %d. max: %d", args.Rooms, maxRooms)
	}
	if args.ErrorCounter > maxErrors {
		return CmdArguments{}, fmt.Errorf("invalid error counter: %d. max: %d", args.ErrorCounter, maxErrors)
	}
	if args.LoopPause > maxPause {
		return CmdArguments{}, fmt.Errorf("invalid loop pause: %d. max: %d", args.LoopPause, maxPause)
	}

	args.LogFileName = fmt.Sprintf("%d_%d.log", args.Rooms, args.Price)
	args.DBFileName = fmt.Sprintf("%d_%d.json", args.Rooms, args.Price)

	return args, nil
}
