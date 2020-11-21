import { Octokit } from "https://cdn.skypack.dev/@octokit/core";

// Updates the visibility of the container
window.hideContainer = function (elem, checkbox_container){
    //Fetches the checkbox container
    var x = document.getElementById(checkbox_container);

    // Updates the style and button text
    if (x.style.display === "none") {
        x.style.display = "block";
        elem.value = "Hide";
    } else {
        x.style.display = "none";
        elem.value = "Show";
    }
}

var metadata = null;
var player_stats = null;
var all_games = null;

// Loads players and corporations from database
async function loadExternalData(){
    const octokit = new Octokit();
    // Decode base64 file content and parse json
    all_games = JSON.parse(b64_to_utf8((await getFile("data/games.json", octokit)).data.content));
    metadata = JSON.parse(b64_to_utf8((await getFile("data/data.json", octokit)).data.content));
    player_stats = JSON.parse(b64_to_utf8((await getFile("data/player_stats.json", octokit)).data.content));
    // Create datalist for HTML
    var datalist = metadata.player_names.map( player => `<option value="${player}">`).join('\n');
    document.getElementById("datalist-players").innerHTML = datalist;

    var placement_games = all_games.length / 10;

    // Add buttons
    var buttons = metadata.player_names
        .map( player => `
            <button class="player-link-button" onclick="displayPlayerStatsToggle('${player}');">${player}</button>
            <button class="player-link-button" disabled style="width:125px">${Math.floor(player_stats[player].elo.value)}${player_stats[player].games.length < placement_games ? "?":""}</button>
        `)
        .join('\n');
    document.getElementById("player-list").innerHTML = buttons;
}

loadExternalData();

//Add listener to password field
var game_search = document.getElementById("player-search");
game_search.addEventListener("keyup", function(event) {
    if (event.keyCode === 13) {
        event.preventDefault();
        displayPlayerStats(game_search.value);
    }
});

window.playerSearch = function(){
    displayPlayerStats(document.getElementById("player-search").value);
}

window.displayPlayerStatsToggle = function(name){
    hideContainer( document.getElementById("hide-players"), "player-list");
    
    // TODO: REMOVE
    console.log(player_stats[name]);

    displayPlayerStats(name);
}

window.displayPlayerStats = function (name){
    var player_data = player_stats[name];
    document.getElementById("player-name").value = name;    
    document.getElementById("player-elo").value = Math.floor(player_data.elo.value);

    document.getElementById("player-total-games").value = player_data.games.length;
    document.getElementById("fave-corp-number").value = `${player_data.fave_corp.total} games / ${player_data.fave_corp.wins} wins`;
    document.getElementById("fave-corp").value = player_data.fave_corp.name;

    document.getElementById("best-corp-number").value = `${player_data.best_corp.wins/player_data.best_corp.total * 100}% ( ${player_data.best_corp.wins} / ${player_data.best_corp.total} )`;
    document.getElementById("best-corp").value = player_data.best_corp.name;

    document.getElementById("win-corp-number").value = `${player_data.most_wins.wins} of ${player_data.most_wins.total} games`;
    document.getElementById("win-corp").value = player_data.most_wins.name;

    document.getElementById("stat-container").style.display = "block";

    for(var i = 1; i <= 5; ++i){
        document.getElementById(`${i}-player`).style.display = "none";
    }

    player_data.best_games.forEach(function(game){
        var all = player_data.games.filter(elem => elem.players == game.players);
        var won = all.filter(elem => elem.rank == 1);

        document.getElementById(`${game.players}-total`).innerHTML = all.length;
        document.getElementById(`${game.players}-won`).innerHTML = won.length
        document.getElementById(`${game.players}-winrate`).innerHTML = `${(won.length/all.length*100).toFixed(2)} %`
        document.getElementById(`${game.players}-best`).innerHTML = game.score;
        document.getElementById(`${game.players}-link`).onclick = function(){window.location = `/games/${game.id}.html`};
        document.getElementById(`${game.players}-player`).style.display = "";
    });

    var all = player_data.games;
    var won = all.filter(elem => elem.rank == 1);
    var max_points = Math.max(...player_data.best_games.map(game => game.score));
    var max_id = player_data.best_games.find(game => game.score == max_points).id;

    document.getElementById("all-total").innerHTML = all.length;
    document.getElementById("all-won").innerHTML = won.length
    document.getElementById("all-winrate").innerHTML = `${(won.length/all.length*100).toFixed(2)} %`
    document.getElementById("all-best").innerHTML = max_points;
    document.getElementById("all-link").onclick = function(){window.location = `/games/${max_id}.html`};


    // Point table

    var table = document.getElementById("game-table");

    var table_header = `
    <tr>
        <td class="table-cell">Game (link)</td>
        <td class="table-cell">Score</td>
        <td class="table-cell">Rank</td>
        <td class="table-cell">Elo</td>
    </tr>
    `
    var table_contents = ""

    var games = {};
    player_data.games.forEach(game => games[game.id] = {score:game.points, rank:game.rank});
    Object.keys(games).forEach(function(id){
        var game = all_games.find(game=>game.id == id);
        
        // REMOVE
        console.log(game.id);
        
        var elo = player_data.elo.history.find(game=>game.id == id);
        table_contents += `
        <tr>
            <td class="table-cell table-button" onclick="window.location = '/games/${id}.html'">${game.name}</td>
            <td class="table-cell">${games[id].score}</td>
            <td class="table-cell">${games[id].rank}</td>
            <td class="table-cell">${Math.floor(elo.value)} (${elo.change < 0 ? "":"+"}${Math.floor(elo.change)})</td>
        </tr>`
    });

    table.innerHTML = table_header + table_contents;

    // TODO: Tablica rezultata https://stackoverflow.com/questions/45857682/interpolation-of-colors
    document.getElementById("game-container").style.display = "block";
    return;
}
