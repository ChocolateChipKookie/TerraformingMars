- Migrate data from mega json to database.

- The expansions entry in the db can also be just a set of values instead of map value: bool

- Clean up codebase, specifically go through the whole go codebase and manually check if everything makes sense

- When number of players is reduced all the game data referencing the deleted player has to be removed (eg for milestones and awards)
- Make the typescript part use the same data in game-data.json
