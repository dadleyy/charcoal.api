package defs

const (
	// ActivityStreamIdentifier name used in the hash of bg streams for the activity stream.
	ActivityStreamIdentifier = "activity"

	// SocketsStreamIdentifier name used in the hash of bg streams for the socket connections stream.
	SocketsStreamIdentifier = "sockets"

	// GamesStreamIdentifier name used in the hash of bg streams for the game stream.
	GamesStreamIdentifier = "games"

	// GamesStatsStreamIdentifier name used in the hash of bg streams for the game statistics stream.
	GamesStatsStreamIdentifier = "game-stats"

	// GameStatsRoundUpdate name used in the hash of bg streams for the game round updates stream.
	GameStatsRoundUpdate = "round-updated"

	// GameProcessorUserJoined name used in the verb of of activity messages when user joins a game.
	GameProcessorUserJoined = "user-joined"

	// GameProcessorUserLeft name used in the verb of of activity messages when user leaves a game.
	GameProcessorUserLeft = "user-left"

	// GameProcessorGameEnded name used in the verb of of activity messages when user leaves a game.
	GameProcessorGameEnded = "game-ended"
)
