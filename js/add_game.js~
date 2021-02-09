import { Octokit } from "https://cdn.skypack.dev/@octokit/core";

///////////////
//    UTIL
///////////////

// Gets active game modes
function getActiveModes(){
    var mode_selection = document.getElementById("expansion-checkbox-container");
    var mode = Array.from(mode_selection.querySelectorAll("input"))
        .filter(function(elem){return elem.checked;})
        .map(function(elem){return elem.value;});
    return mode;
}

////////////////////
//     WARNING
////////////////////

function updateWarning(){
    // Fetch colony list element and number of selected colonies
    var colony_list = document.getElementById("colony-list")
    var current_colonies = colony_list.querySelectorAll(".colony-container").length;

    // Update the style of the container depending if any of the colonies is selected
    colony_list.style.display = current_colonies == 0 ? "none" : "block";

    // Get the number of players playing
    var players = parseInt(document.getElementById("player-number").value, 10);


    // Set the expected players value
    var expected_colonies = 5;
    if(players != 2){
        expected_colonies = players + 2;
    }

    // Add the exception for Aridor
    var colonies_in_play = Array.from(document.getElementsByClassName("player-colony-select")).map(function(elem){return elem.value;});
    if (colonies_in_play.includes("Aridor")){
        expected_colonies += 1;
    }

    // Set warning if necessary
    var warning = document.getElementById("colony-warning")

    if (expected_colonies != current_colonies){
        
        warning.innerHTML = 
        `WARNING! Wrong number of colonies, ${expected_colonies} expected!`
        warning.style.display = 'block'
    }
    else{
        warning.style.display = 'none'
    }
}

//////////////////////
//    CHECKBOXES
//////////////////////


//Updates the list of colonies when a checkbox is clicked
window.checkboxFunction = function (elem){

    //Gets colony
    var colony = getElementInsideContainer("colony-list", elem.value + '-icon');
    var div_list = document.getElementById("colony-list");

    if (elem.checked && colony == false){
        //Get colony name
        var colony_name = elem.value;
        // Add colony to the colony list
        div_list.innerHTML += 
        `
        <div class="colony-container" id="${colony_name}-icon" onclick="increaseColonyNumber(this)">
            <img src="resources/icons/colonies/${colony_name}.png" class="colony-icon">
            <input type="hidden" value="0">
            <div class="colony-text">
                ${capitalize(colony_name)}<br>(0)
            </div>
        </div>
        `;
    } else if(!elem.checked && colony != false){
        // If enement is unchecked and colony exists remove it
        colony.remove();
    }

    // Updates warning
    updateWarning();
}

// Hides the given checkbox container
window.hideCheckboxContainer = function (elem, checkbox_container){
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

///////////////
//    DATE
///////////////

//Sets the default date
function setDefaultDate(){
    var today = new Date();
    const dateTimeFormat = new Intl.DateTimeFormat('en', { year: 'numeric', month: '2-digit', day: '2-digit' });
    const [{ value: month },,{ value: day },,{ value: year }] = dateTimeFormat .formatToParts(today);
    document.getElementById("datepicker").value = `${year}-${month}-${day}`;
}

setDefaultDate()

/////////////////////
//     PLAYERS
/////////////////////

// Loads players and corporations from database
async function loadPlayerNames(){
    const octokit = new Octokit();
    //Fetch file from git
    var data = await getFile("data/data.json", octokit);
    //Decode base64 file content and parse json
    data = JSON.parse(b64_to_utf8(data.data.content));

    //Add names to datalist
    var names = data.player_names;
    //Create datalist for HTML
    var datalist = names.map(function(name){return `<option value="${name}">`}).join('\n');
    document.getElementById("datalist-players").innerHTML = datalist;
}

loadPlayerNames();

//Generates player input fields
window.generatePlayerInputs = async function (){
    // Fetch number of players
    var number_of_players_input = document.getElementById("player-number");
    // Check and update the range
    if (number_of_players_input.value > 5)
        number_of_players_input.value = 5;
    else if (number_of_players_input.value < 1)
             number_of_players_input.value = 1;

    // Fetch number of players and player input list
    var players = number_of_players_input.value;
    var players_list = document.getElementById("player-list");


    //Initial loading of player names and corporations

    //Add corporations
    var corporations = util.game_data.corporation_names;
    corporations.sort();
    var corporation_list = corporations.map(corp => `<option value="${corp}">${corp}</option>`).join("\n");
    
    // Create input fields
    players_list.innerHTML = 
    `
    <div class="player-input-div">
        <input  class="player-input" id="player" list="datalist-players">
        <select class="player-input"> ${corporation_list}</select>
    </div>
    `.repeat(players);
    
    // Update warning
    updateWarning();

    // Set players input list to visible
    players_list.style.display = "block"
    document.getElementById("base-option-generate").style.display = "block";
}

generatePlayerInputs();


//////////////////
//      MAPS
//////////////////

function setMaps(){
    var map_select = document.getElementById("map");
    var maps = util.game_data.maps;
    map_select.innerHTML = maps.map(elem => `<option value="${elem}">${elem}</option>`).join("\n");;
}

setMaps()


//////////////////////
//     COLONIES
//////////////////////

// Activates colonies expansion
window.activateColonies = function (elem){
    // Fetches the colonies option field
    var colony_options = document.getElementById("colonies-options");
    colony_options.style.display = elem.checked ? "block" : "none";
}

// Increases the colony counter text by one in modulo 4 arithmetic
window.increaseColonyNumber = function (elem){
    // Gets the name of the colony
    var colony = elem.id.substring(0, elem.id.length - 5);

    // Gets the current number of colonies on a moon
    var current_colonies = elem.querySelector("input");
    current_colonies.value = (parseInt(current_colonies.value, 10) + 1) % 4;

    // Gets the text div element of elem
    var text = elem.querySelector("div");
    text.innerHTML = `${capitalize(colony)}<br>(${current_colonies.value})`
}

///////////////////////
//      MILESTONES
///////////////////////

// Checks if there are more than one of the same milestone entered
window.checkMilestonesValidity = function(){
    var active_milestones = Array.from(document.getElementsByClassName("milestone-select"))
        .map(milestone => milestone.value);
    return (new Set(active_milestones)).size == active_milestones.length;
}

window.parseMilestones = function(){
    var players = Array.from(document.getElementsByClassName("milestone-select-player")).map(elem => elem.value);
    var milestones = Array.from(document.getElementsByClassName("milestone-select")).map(elem => elem.value);
    var result = {};

    for(var i = 0; i < milestones.length; ++i){
        result[milestones[i]] = players[i];
    }
    
    return result;
}

//Update scores for milestone
window.updateMilestoneScores = function (){
    var milestones = Array.from(document.getElementsByClassName("milestone-select-player")).map(elem => elem.value);
    var player_counts = {}
    milestones.forEach(function(milestone){
        player_counts[milestone] = player_counts[milestone] ? player_counts[milestone] + 1 : 1;
    });

    var player_list = Array.from(document.getElementById("player-list").querySelectorAll("input")).map(elem => elem.value);
    var milestone_points = Array.from(document.getElementById("milestone-points").querySelectorAll("input"));

    for(var i = 0; i < player_list.length; ++i){
        var no_milestones = player_counts[player_list[i]] ? player_counts[player_list[i]] : 0;
        milestone_points[i].value =  5 * no_milestones;
    }

    updateTotals();
}

//Generates player input fields
window.updateMilestoneInputs = function (){
    // Fetch number of milestones
    var number_of_milestones_input = document.getElementById("milestone-number");
    // Check and update the range
    if (number_of_milestones_input.value > 3)
        number_of_milestones_input.value = 3;
    else if (number_of_milestones_input.value < 0)
        number_of_milestones_input.value = 0;

    // Fetch number of milestones and player input list
    var milestones = number_of_milestones_input.value;
    var milestone_list = document.getElementById("milestone-list");

    //Initial loading of milestones names and corporations

    function getMilestones(){
        var active_milestones = [];
        var active_map = document.getElementById("map").value;
        active_milestones.push.apply(active_milestones, util.game_data.milestones[active_map]);
        if(document.getElementById("venus-checkbox").checked){
            active_milestones.push.apply(active_milestones, util.game_data.milestones["Venus"]);
        }
        active_milestones.sort();
        return active_milestones;
    }

    function getAllMilestones(){
        var active_milestones = [];
        active_milestones.push.apply(active_milestones, util.game_data.milestones["all"]);
        active_milestones.sort();
        return active_milestones;
    }

    //var milestones_select = getMilestones().map(milestone => `<option value="${milestone}">${milestone}</option>`).join("\n");
    var milestones_select = getAllMilestones().map(milestone => `<option value="${milestone}">${milestone}</option>`).join("\n");
    var active_players_select = util.game_data.active_players.map(player => `<option value="${player}">${player}</option>`).join("\n");
    
    // Create input fields
    milestone_list.innerHTML = 
        `<div class="milestone-container">
            <select class="milestone-select"> ${milestones_select}</select>
            <select class="milestone-select-player" onchange="updateMilestoneScores()"> ${active_players_select}</select>
        </div>`.repeat(milestones);
    
    updateMilestoneScores();

    // Set milestone input list to visible
    document.getElementById("milestones").style.display = milestones == 0 ? "none":"block";
}

//////////////////////
//     AWARDS
//////////////////////

// Checks if there are more than 1 of the same awards entered
window.checkAwardsValidity = function(){
    var active_awards = Array.from(document.getElementsByClassName("award-points-row"))
        .map(award => award.children[0].children[0].value);
    return (new Set(active_awards)).size == active_awards.length;
}

// Parse award points
window.parseAwardPoints = function(){
    var player_names = Array.from(document.getElementById("player-list").children).map( element => element.children[0].value);
    var awards = Array.from(document.getElementsByClassName("award-points-row"));
    var awards_points = {}

    // For each award:
    awards.forEach(function(award_row){
        var award_name = award_row.children[0].children[0].value;
        awards_points[award_name] = {};

        //  Get all the scores:
        var scores = Array.from(award_row.querySelectorAll("input"))
            .map(award => parseInt(award.value, 10))
            .map(score => isNaN(score) ? 0:score);

        for(var i = 0; i < player_names.length; ++i){
            awards_points[award_name][player_names[i]] = scores[i]
        }

    });
    return awards_points;
}

// Calculates who won what awards
window.calculateAwards = function(){
    var player_names = Array.from(document.getElementById("player-list").children).map( element => element.children[0].value);
    var awards = Array.from(document.getElementsByClassName("award-points-row"));
    var player_awards = {}

    // For each award:
    awards.forEach(function(award_row){
        var award_name = award_row.children[0].children[0].value;
        player_awards[award_name] = {};
        //  Get all the scores:
        var scores = Array.from(award_row.querySelectorAll("input"))
            .map(award => parseInt(award.value, 10))
            .map(score => isNaN(score) ? 0:score);

        //  Get the max score
        var maxes = Array.from(new Set(scores)).sort(function(a, b){return b - a;});

        //  If only one player with max score
        if (scores.filter(score => score == maxes[0]).length == 1){
            // Give that player first place for the award
            var first_index = scores.findIndex(score => score == maxes[0]);
            player_awards[award_name].first = [player_names[first_index]];
            player_awards[award_name].second = [];
            // Check the second highest score
            // Give all players with second highest score the second place
            for (var i = 0; i < player_names.length; ++i){
                if (scores[i] == maxes[1]){
                    player_awards[award_name].second.push(player_names[i]);
                }
            }
        }
        else{
            player_awards[award_name].first = [];
            player_awards[award_name].second = [];
            // Give all players with highest score the first place
            for (var i = 0; i < player_names.length; ++i){
                if (scores[i] == maxes[0]){
                    player_awards[award_name].first.push(player_names[i]);
                }
            }
        }
    });
    return player_awards;
}

//Update scores for awards
window.updateAwardScores = function (){
    var award_scores = calculateAwards();
    
    var player_names = Array.from(document.getElementById("player-list").children).map( element => element.children[0].value);
    var totals = Array.from(document.getElementById("award-table-points").querySelectorAll("input"))
    var totals_point_table = Array.from(document.getElementById("award-points").querySelectorAll("input"))
    
    for(var i = 0; i < player_names.length; ++i){
        var entries = Object.entries(award_scores);
        var total_score = 0;

        entries.forEach(function(entry){
            if(entry[1].first.includes(player_names[i]))        total_score += 5;
            else if(entry[1].second.includes(player_names[i]))  total_score += 2;
        });
        totals[i].value = total_score;
        totals_point_table[i].value = total_score;
    }

    updateTotals();
}

//Generates player input fields
window.updateAwardInputs = function (){
    // Gets the points div, points table and player list nodes
    var awards_div = document.getElementById("awards");
    var award_table = document.getElementById("award-table");
    var player_list = document.getElementById("player-list");
    var award_number = document.getElementById("award-number").value;

    if(award_number == 0){
        awards_div.style.display = "none";
        return;
    }
    
    // Fetches the names of the players
    var player_names = Array.from(player_list.children).map(
        function(element) { return element.children[0].value; }
    );

    // Creates the name row for the table
    var name_row = 
    `<tr id="award-names">
        <td class="table-cell">Name</td>
        ${player_names.map(element => `<td class="table-cell">${element}</td>`).join("\n")}
    </tr>`;

    // Creates the inputs for the scores
    var createInputs = (type, additional = "") => `<td class="table-cell"><input class="table-input" type="${type}" ${additional}></td>`.repeat(player_names.length);

    function getAwards(){
        var active_awards = [];
        var active_map = document.getElementById("map").value;
        active_awards.push.apply(active_awards, util.game_data.awards[active_map]);
        if(document.getElementById("venus-checkbox").checked){
            active_awards.push.apply(active_awards, util.game_data.awards["Venus"]);
        }
        active_awards.sort();
        return active_awards;
    }

    function getAllAwards(){
        var active_awards = [];
        active_awards.push.apply(active_awards, util.game_data.awards["all"]);
        active_awards.sort();
        return active_awards;
    }

//    var awards_select = getAwards().map(award => `<option value="${award}">${award}</option>`).join("\n");
    var awards_select = getAllAwards().map(award => `<option value="${award}">${award}</option>`).join("\n");

    // Adds all point rows
    var awards = `
    <tr class="award-points-row">
        <td class="table-cell"><select class="award-select" >${awards_select}</select></td>
        ${createInputs("number", 'onchange="updateAwardScores()" min="-5" value="0"')}
    </tr>
    `.repeat(award_number)

    var awards_points = `
    <tr id="award-table-points">
        <td class="table-cell">Total points</td>
        ${createInputs("number", 'disabled value="0"')}
    </tr>
    `
    
    // Creates the table and sets the style to visible
    award_table.innerHTML = name_row + awards + awards_points;

    document.getElementById("awards").style.display = "block";

    // Show points div
    awards_div.style.display = 'block';
}

//////////////////
//     TABLE
//////////////////

// Uncheck the checkbox
document.getElementById("gold-lead-input").checked = false;

// Updates totals and ranks
window.updateTotals = function(){
    // Define keys of scores that are tracked
    var categories = [
        "tr-points",
        "award-points",
        "milestone-points",
        "greenery-points",
        "city-points",
        "card-points",
        "gold-lead-points"
    ]

    // Case turmoil is in game
    if (getActiveModes().includes("Turmoil")){
        categories.push("turmoil-points")
    }

    // Fetch number of players
    var no_players = parseInt(document.getElementById("player-number").value, 10);
    // Init scores
    var scores = new Array(no_players).fill(0);
    // Calculate totals, iterating over all scores
    categories.forEach(
        function(key){
            // Get all points from the category defined by the category key
            var points = document.getElementById(key).querySelectorAll("input");
            for(var i = 0; i < no_players; ++i){
                // Try parsing the key, in case it is not a number, ignore 
                var value = parseInt(points[i].value, 10);
                if(!isNaN(value)){
                    scores[i] += value;
                }
            }
        }
    );

    // Define Decorate-Sort-Undecorate
    const dsu = (arr1, arr2) => arr1
        .map((item, index) => [arr2[index], item]) 
        .sort(([arg1], [arg2]) => arg2 - arg1) 
        .map(([, item]) => item);

    // Create arranged array, and calculate indexes
    const iota = Array.from(Array(no_players).keys());
    const result = dsu(iota, scores);

    // Update ranks
    var ranks = document
        .getElementById("ranks")
        .querySelectorAll("input");

    for (var i = 0; i < no_players; ++i){
        ranks[result[i]].value = i + 1;
    }

    // Update totals
    var totals = document
        .getElementById("total-points")
        .querySelectorAll("input");
    // Get gold lead points
    var gold_lead = document.getElementById("gold-lead-points").querySelectorAll("input");
    for (var i = 0; i < no_players; ++i){
        // Check if gold lead points are present
        var value = parseInt(gold_lead[i].value, 10);
        if(!isNaN(value)){
            //Remove gold lead from total
            scores[i] -= value;
        }
        // Update total
        totals[i].value = scores[i];
    }

}

// Update the visibility of the gold lead field
window.goldLeadFunction = function(){
    var gold_lead = document.getElementById("gold-lead-input").checked;
    document.getElementById("gold-lead-points").style.display = gold_lead ? "" : "none";
}

// Update the visibility of the turmoil field
window.updateTurmoil = function(){
    var gold_lead = document.getElementById("turmoil-checkbox").checked;
    try{
        document.getElementById("turmoil-points").style.display = gold_lead ? "" : "none";
    } catch {}
}

// Generates the point table for the scores
window.generatePointTable = function(){
    // Gets the points div, points table and player list nodes
    var points_div = document.getElementById("points");
    var table = document.getElementById("points-table");
    var player_list = document.getElementById("player-list");

    // Fetches the names of the players
    var player_names = Array.from(player_list.children).map(
        function(element) { return element.children[0].value; }
    );

    util.game_data.active_players = player_names;

    // Creates the name row for the table
    var name_row = 
    `<tr id="player-names">
        <td class="table-cell">Name</td>
        ${player_names.map(element => `<td class="table-cell">${element}</td>`).join("\n")}
    </tr>`;

    // Constant value for the score categories
    var categories = [
        ["TR", "tr-points"], 
        ["Awards", "award-points"], 
        ["Milestones", "milestone-points"], 
        ["Greenery", "greenery-points"], 
        ["Cities", "city-points"], 
        ["Cards", "card-points"],
        ["Turmoil", "turmoil-points"],
        ["Gold lead", "gold-lead-points"]
    ];

    // Creates the inputs for the scores
    var createInputs = (type, additional = "") => `<td class="table-cell"><input class="table-input" type="${type}" ${additional}></td>`.repeat(player_names.length);

    // Adds all point rows
    var points = categories.map( option =>
    `<tr id="${option[1]}">
        <td class="table-cell">${option[0]}</td>
        ${createInputs("number", 'onchange="updateTotals()"')}
    </tr>
    `)

    document.getElementsByClassName("")

    // Adds total points rows
    var total_points = 
        `<tr id="total-points">
            <td class="table-cell">Total</td>
            ${createInputs("text", 'value="0" disabled')}
        </tr>`

    // Adds ranks rows
    var ranks = 
        `<tr id="ranks">
            <td class="table-cell">Rank</td>
            ${createInputs("text", 'value="1" disabled')}
        </tr>`

    // Adds possibility for notes
    var note_input = 
        `<tr id="player-note">
            <td class="table-cell">Notes</td>
            ${createInputs("text")}
        </tr>`
    
    // Creates the table and sets the style to visible
    table.innerHTML = name_row + points.join('\n') + total_points + ranks + note_input;

    //Disable awards
    ["award-points", "milestone-points"].forEach(key =>
        Array.from(document.getElementById(key).querySelectorAll("input")).forEach(function(elem){
            elem.disabled = true;
            elem.value = 0;
        })
    );
    



    // Update display of gold lead
    goldLeadFunction();

    updateTurmoil();

    updateAwardInputs();
    updateMilestoneInputs();
    document.getElementById("award-options").style.display = "block";    
    document.getElementById("milestone-options").style.display = "block";    

    // Show points div
    points_div.style.display = 'block';
}

//////////////////
//     SUBMIT
//////////////////

function createCode(modes){
    var mode_dict = {
        "Base": "B",
        "Corporate Era": "CE",
        "Prelude": "Pr",
        "Colonies": "Col",
        "Draft": "Dr",
        "Venus": "Ve",
        "Turmoil": "Tu",
    }
    return modes.map( mode => mode_dict[mode]).join("+");
}

// Generates the html for the page
function generateGameSite(game){

    function transformDate(date){
        return `${date.substring(8)}/${date.substring(5, 7)}/${date.substring(0, 4)}`
    }

    var names = game.scores.map(score => `                    <td class="table-cell" style="font-weight: bolder; font-family: 'Courier New', Courier, monospace;">${score.player}</td>`).join('\n');
    var corp = game.scores.map(score => `                    <td class="table-cell">${score.corporation}</td>`).join('\n');
    var tr = game.scores.map(score => `                    <td class="table-cell">${score.tr}</td>`).join('\n');
    var award = game.scores.map(score => `                    <td class="table-cell">${score.awards}</td>`).join('\n');
    var milestones = game.scores.map(score => `                    <td class="table-cell">${score.milestones}</td>`).join('\n');
    var greenery = game.scores.map(score => `                    <td class="table-cell">${score.greenery}</td>`).join('\n');
    var cities = game.scores.map(score => `                    <td class="table-cell">${score.cities}</td>`).join('\n');
    var cards = game.scores.map(score => `                    <td class="table-cell">${score.cards}</td>`).join('\n');
    var lead = game.scores.map(score => `                    <td class="table-cell">${score.lead}</td>`).join('\n');
    var total = game.scores.map(score => `                    <td class="table-cell">${score.total}</td>`).join('\n');
    var rank = game.scores.map(score => `                    <td class="table-cell">${score.rank}</td>`).join('\n');
    var note = game.scores.map(score => `                    <td class="table-cell" style="text-transform: none; font-weight: bolder; font-family: 'Courier New', Courier, monospace;">${score.note}</td>`).join('\n');

    lead = 
        `
                <tr id="gold-lead">
                    <td class="table-cell">Gold lead</td>
${lead}
                </tr>
        `

    if(game.scores.map(score => score.lead).every(item => item === 0)){
        lead = ""
    }

    var turmoil = ""

    if(game.mode.includes("Turmoil")){

    turmoil = 
        `
                <tr id="turmoil-points">
                    <td class="table-cell">Turmoil</td>
${game.scores.map(score => `                    <td class="table-cell">${score.turmoil}</td>`).join('\n')}
                </tr>
        `
    }

    // TODO: Add milestones and awards if they are part of the entry
    // if (entry[milestones]) ...
    // if (entry[awards]) ...

    return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <title>${game.name}</title>
    <link href="../css/meta.css" rel="stylesheet" type="text/css">    
    <link href="../css/game.css" rel="stylesheet" type="text/css">    
    <link href="../css/table.css" rel="stylesheet" type="text/css">
</head>
<body>
    <div class="container page-title">
        ${game.name}
    </div>

    <div class="container">
        <div class="container-banner">
            Game outcome
        </div>
        <div class="sub-container">
            <div class="sub-container--element">
                <label for="name-input">Name:</label>
                <input type="text" class="option-input" id="name-input" disabled value="${game.name}">
            </div>
            <div class="sub-container--element">
                <label for="date-input">Date:</label>
                <input type="text" class="option-input" id="date-input" disabled value="${transformDate(game.date)}" required="required">
            </div>
            <div class="sub-container--element">
                <label for="mode-input">Mode:</label>
                <input type="text" class="option-input" id="mode-input" disabled value="${game.mode_code}">
            </div>
            <div class="sub-container--element">
                <label for="map-input">Map:</label>
                <input type="text" class="option-input" id="map-input" disabled value="${game.map}">
            </div>
            <div class="sub-container--element">
                <label for="players-input">Players:</label>
                <input type="text" class="option-input" id="players-input" disabled value="${game.players}">
            </div>
            <div class="sub-container--element">
                <label for="generations-input">Generations:</label>
                <input type="text" class="option-input" id="generations-input" disabled value="${game.generation}">
            </div>
            <div class="sub-container--element">
                <label for="winner-input">Winner:</label>
                <input type="text" class="option-input" id="winner-input" style="font-weight: bolder; font-family: 'Courier New', Courier, monospace;" disabled value="${game.winner}">
            </div>
            <div class="sub-container--element">
                <label for="win-corp-input">Winner corporation:</label>
                <input type="text" class="option-input" id="win-corp-input" disabled value="${game.win_corp}">
            </div>
            <div class="sub-container--element">
                <label for="win-score-input">Winner points:</label>
                <input type="text" class="option-input" id="win-score-input" disabled value="${game.win_score}">
            </div>
        </div>

        <div class="container-banner">
            Points
        </div>
        <div class="sub-container sub-container--element" id="points-table-div">
            <table class="table" id="points-table">
                <tr id="player-names">                    
                    <td class="table-cell">Name</td>
${names}
                </tr>
                <tr id="player-corporations">                    
                    <td class="table-cell">Corporation</td>
${corp}
                </tr>
                <tr id="tr-points">
                    <td class="table-cell">TR</td>
${tr}
                </tr>
                <tr id="award-points">
                    <td class="table-cell">Awards</td>
${award}
                </tr>
                <tr id="milestone-points">
                    <td class="table-cell">Milestones</td>
${milestones}
                </tr>
                <tr id="greenery-points">
                    <td class="table-cell">Greenery</td>
${greenery}
                </tr>
                <tr id="city-points">
                    <td class="table-cell">Cities</td>
${cities}
                </tr>
                <tr id="card-points">
                    <td class="table-cell">Cards</td>
${cards}
                </tr>
${turmoil}
                <tr id="total-points">
                    <td class="table-cell">Total</td>
${total}
                </tr>
${lead}
                <tr id="ranks">
                    <td class="table-cell">Rank</td>
${rank}
                </tr>
                <tr id="player-note">
                    <td class="table-cell">Notes</td>
${note}
                </tr>
            </table>
        </div>
    </div>
    <div class="link-button-wrapper">
        <button class="link-button-half" onclick="window.location='/games/${game.id - 1}.html';">Previous</button>
        <button class="link-button-half" onclick="window.location='/games/${game.id + 1}.html';">Next</button>
    </div>
    <div class="link-button-wrapper">
        <button class="link-button" onclick="window.location='../index.html';">Main page</button>
    </div>
</body>
<footer>
</footer>
`
}

window.calculateEloParameters = function(games){
    var means = {}
    games.forEach(function(game){
        if(!means[game.players]){
            means[game.players] = {total:0, count:0};
        }
        game.scores.forEach(function(score){
            means[game.players].total += score.total;
            means[game.players].count += 1;
        });
    });
    Object.keys(means).forEach(key => means[key].mean = means[key].total / means[key].count);
    
    var dev = {}
    games.forEach(function(game){
        if(!dev[game.players]){
            dev[game.players] = {total:0, count:0};
        }
        game.scores.forEach(function(score){
            var diff = (score.total - means[game.players].mean);
            dev[game.players].total += diff*diff;
            dev[game.players].count += 1;
        });
    });
    
    Object.keys(dev).forEach(key => dev[key].stddev = Math.sqrt(dev[key].total / dev[key].count));
    
    return {means: means, stddev:dev};
}

function calculateEloChange(game, player_stats, means, stddevs){
    var mean = means[game.players].mean;
    var stddev = stddevs[game.players].stddev;
    
    var players = game.scores.map(elem => elem.player);
    const pool_percentage = 0.05;
    var pool = players.reduce((acc, player) => acc + (player_stats[player].elo.value * pool_percentage), 0);


    function getFactor(score){
        function erf(x) {
            var z;
            const ERF_A = 0.147; 
            var the_sign_of_x;
            if(0==x) {
                the_sign_of_x = 0;
                return 0;
            } else if(x>0){
                the_sign_of_x = 1;
            } else {
                the_sign_of_x = -1;
            }

            var one_plus_axsqrd = 1 + ERF_A * x * x;
            var four_ovr_pi_etc = 4/Math.PI + ERF_A * x * x;
            var ratio = four_ovr_pi_etc / one_plus_axsqrd;
            ratio *= x * -x;
            var expofun = Math.exp(ratio);
            var radical = Math.sqrt(1-expofun);
            z = radical * the_sign_of_x;
            return z;
        }

        var normalized = (score - mean) / stddev;
        const sqrt2 = 1.41421356237;

        return (erf(normalized/sqrt2) + 1)/2;
    }

    var scores = {}
    game.scores.forEach(elem => scores[elem.player] = 
        {
            score: elem.total, 
            elo: player_stats[elem.player].elo.value,
            factor: getFactor(elem.total)
        });

    var total_factors = Object.keys(scores).reduce((acc, key) => scores[key].factor + acc, 0);
    players.forEach(function(player) {
        scores[player].factor /= total_factors;
        scores[player].change = pool * scores[player].factor - scores[player].elo * pool_percentage;
    });

    return scores;
}

function recalculateStats(player, game, player_stats){
    // Check if in data stats
    // If not add entry
    if (!player_stats[player]){
        var stats = game.scores.find(elem => elem.player == player);
        var wins = stats.rank == 1 ? 1 : 0;
        var game_stats = {
            games: [
                {
                    id: game.id,
                    corporation: stats.corporation,
                    rank: stats.rank,
                    points: stats.total,
                    players: game.players
                }
            ],
            best_games: [
                {
                    players: game.players,
                    id: game.id,
                    score: stats.total
                }
            ],
            fave_corp: {
                name: stats.corporation,
                total: 1,
                wins: wins
            },
            best_corp: {
                name: stats.corporation,
                total: 1,
                wins: wins
            },
            most_wins: {
                name: stats.corporation,
                total: 1,
                wins: wins
            },
            elo: {
                value: 1000,
                history: [
//                    {
//                        id: 0,
//                        delta: 0
//                    }
                ]
            }
        }
        player_stats[player] = game_stats;
        return player_stats;
    }

    // Else
    var game_stats = game.scores.find(elem => elem.player == player);
    var player_stats = player_stats[player];

    // Add game to games
    player_stats.games.push({
        id: game.id,
        corporation: game_stats.corporation,
        rank: game_stats.rank,
        points: game_stats.total,
        players: game.players
    });

    // Check player number, and update best game if applicable
    var best_game_index = player_stats.best_games.findIndex(elem => elem.players == game.players);
    if(best_game_index == -1){
        player_stats.best_games.push({
                players: game.players,
                id: game.id,
                score: game_stats.total
            });
    } else if (player_stats.best_games[best_game_index].score < game_stats.total){
        player_stats.best_games[best_game_index].id = game.id;
        player_stats.best_games[best_game_index].score = game_stats.total;
    }
        
    var occurances = {}
    player_stats.games.forEach(
        function (game){
            if(occurances[game.corporation]){
                occurances[game.corporation].total += 1;
                occurances[game.corporation].wins += game.rank == 1 ? 1:0;
            }
            else{
                occurances[game.corporation] = {total: 1, wins: (game.rank == 1 ? 1:0)};
            }
        }
    )

    occurances = Object.entries(occurances);
    
    function maxElem(array, key){
        var tmp = array.map(key);
        var max_value = Math.max(...tmp);
        return array[tmp.findIndex(elem => elem == max_value)];
    }

    function maxElems(array, key){
        var max_value = Math.max(...array.map(key));
        return array.filter(elem => key(elem) == max_value);
    }

    // If fave corp, increment
    if (player_stats.fave_corp.name == game_stats.corporation){
        player_stats.fave_corp.total += 1;
    } else{
        var fave = maxElems(occurances, corp => corp[1].total);
        fave = maxElem(fave, corp => corp[1].wins / corp[1].total);
        player_stats.fave_corp = {name: fave[0], total: fave[1].total, wins: fave[1].wins};
    }

    // Update best corp if applicable
    var best = maxElems(occurances, corp => corp[1].wins / corp[1].total);
    best = maxElem(best, corp => corp[1].total);
    player_stats.best_corp = {name: best[0], total: best[1].total, wins: best[1].wins};

    // Update most wins if applicablee
    var most_wins = maxElems(occurances, corp => corp[1].wins);
    most_wins = maxElem(most_wins, corp => -corp[1].total);
    player_stats.most_wins = {name: most_wins[0], total: most_wins[1].total, wins: most_wins[1].wins};

    return player_stats;
}

window.submitForm = async function(){
    // Fetch submit button and set it to disabled
    var submit_button = document.getElementById("submit-button");
    submit_button.disabled = true;

    // Update banner to given value and enable submit button
    function updateBanner(value, isWarning){
        var banner = document.getElementById("game-added-banner");
        banner.style.backgroundColor = isWarning ? "black" : "rgb(240, 240, 240)";
        banner.style.color = isWarning ? "white" : "black";
        banner.style.display = "block";
        banner.innerHTML = value;
        submit_button.disabled = false;
    }

    // Passeword management
    var password_field = document.getElementById("password-field");
    var password = password_field.value;
    password_field.value = "";

    ////////////////////////
    //      Program exits
    ////////////////////////

    // Check if password is valid
    if (CryptoJS.SHA256(password).toString() != password_hash){
        updateBanner("Wrong password!", true);
        return;
    }

    // Check if milestones is valid
    if (!checkMilestonesValidity()){
        updateBanner("Repeating milestones not allowed!", true);
        return;
    }

    // Check if award is valid
    if (!checkAwardsValidity()){
        updateBanner("Repeating awards not allowed!", true);
        return;
    }

    // Fetch player names
    var player_names = Array.from(document.getElementsByClassName("player-input-div"))
        .map( elem => elem.children[0].value);

    // Check if there are several same names or "" in player names    
    if (player_names.includes("")){
        updateBanner("Empty name not allowed", true);
        return;
    } 
    else if (player_names.length != (new Set(player_names)).size){
        updateBanner("Players with same name not allowed", true);
        return;
    }


    // Get GitHub api token by decrypting encrypted token
    var decrypted = CryptoJS.AES.decrypt(encrypted_token, password);

    // Init Octokit
    const octokit = new Octokit({auth: decrypted.toString(CryptoJS.enc.Utf8)});

    // Fetch global values
    var name = document.getElementById("name-input").value;
    var note = document.getElementById("game-notes").value;
    var date = document.getElementById("datepicker").value;

    var mode = getActiveModes();
    var mode_code = createCode(mode);

    var map = document.getElementById("map").value;
    var no_players = parseInt(document.getElementById("player-number").value, 10);
    var generations = parseInt(document.getElementById("generations").value, 10);

    // Fetch player specific inputs
    var players = Array.from(document.getElementsByClassName("player-input-div"))
        .map(function(elem){
            return {player: elem.children[0].value, corporation: elem.children[1].value};
        });

    var categories = [
        ["tr-points", "tr"],
        ["award-points", "awards"],
        ["milestone-points", "milestones"],
        ["greenery-points", "greenery"],
        ["city-points", "cities"],
        ["card-points", "cards"],
        ["gold-lead-points", "lead"],
        ["total-points", "total"],
        ["ranks", "rank"],
        ["player-note", "note"],
    ];

    if (mode.includes("Turmoil")){
        categories.push(["turmoil-points", "turmoil"]);
    }

    // Iterate over the categories and add them to player stats
    categories.forEach(
        function(elem){
            var key = elem[0];
            var category = elem[1];
            var points = document
                .getElementById(key)
                .querySelectorAll("input");

            for(var i = 0; i < no_players; ++i){
                //Points vs note in the end
                if (category == "note"){
                    players[i]["note"] = points[i].value;
                }
                else{
                    players[i][category] = points[i].value == "" ? 0 : parseInt(points[i].value, 10);
                }

            }
        }
    )

    // Get winner
    var winner = players.find(player => player.rank == 1);

    // Get data
    var metadata = JSON.parse(b64_to_utf8((await getFile("data/data.json", octokit)).data.content));;
    var games_data = JSON.parse(b64_to_utf8((await getFile("data/games.json", octokit)).data.content));;
    var player_stats = JSON.parse(b64_to_utf8((await getFile("data/player_stats.json", octokit)).data.content));;
 
    // Check if game with same name exists
    if (games_data.map(game => game.name).includes(name)){
        updateBanner("Game with same name already exists", true);
        return;
    } 
    
    // Create entry
    var entry = {
        id: metadata.current_id,
        name: name,
        note: note,
        date: date,
        mode: mode,
        mode_code: mode_code,
        map: map,
        winner: winner.player,
        win_corp: winner.corporation,
        win_score: winner.total,
        players: no_players,
        generation: generations,
        // TODO: Add image loading
        image:"",
        scores: players
    }

    // Increase id counter
    metadata.current_id++;

    // Add milestones
    entry.milestones = parseMilestones();

    // Add awards
    entry.awards = {winners: calculateAwards(), points: parseAwardPoints()};

    // Add colonies
    if(mode.includes("Colonies")){
        var colonies = document.getElementsByClassName("colony-container");
        colonies = Array.from(colonies).map(function(entry){
                    return {
                        name: entry.children[2].childNodes[0].nodeValue.trim(),
                        count: parseInt(entry.children[1].value, 10)
                    };
                }
            );
        entry.colonies = colonies;
    }

    // Add game
    games_data.push(entry);

    // Add players
    player_names.forEach(function (name){
        if (!metadata.player_names.includes(name)){
            metadata.player_names.push(name);
        }
    });

    // Update stats
    player_names.forEach(function(player){
        recalculateStats(player, entry, player_stats);
    });

    // Create game site
    var game_site = generateGameSite(entry);
    await createFile(game_site, `games/${entry.id}.html`, `Created site for game "${name}" id:${entry.id}`, octokit)
    updateBanner("Game added!", true);

    // Push files
    await updateFile(JSON.stringify(player_stats, null, 2), "data/player_stats.json", `Added game "${name}"`, octokit);
    await updateFile(JSON.stringify(metadata, null, 2), "data/data.json", `Added game "${name}"`, octokit);
    await updateFile(JSON.stringify(games_data, null, 2), "data/games.json", `Added game "${name}"`, octokit);

    //Update elo
    await recalculateAllEloScores(password);

    window.location = '/index.html'
}

//Add listener to password field
var password_field = document.getElementById("password-field");
password_field.addEventListener("keyup", function(event) {
    if (event.keyCode === 13) {
        event.preventDefault();
        submitForm();
    }
});

///////////////////
//   ADMIN UTIL
///////////////////

//Used to update the encrypted token in function above
window.encryptToken = function (token, password){
    var encrypted = CryptoJS.AES.encrypt(token, password);
    console.log("Encrypted", encrypted.toString());
    console.log("SHA", CryptoJS.SHA256(password).toString());
}

window.calculateAllStats = async function (password){
    // Get GitHub api token by decrypting encrypted token
    var decrypted = CryptoJS.AES.decrypt(encrypted_token, password);

    // Init Octokit
    const octokit = new Octokit({auth: decrypted.toString(CryptoJS.enc.Utf8)});

    // Add all players to the database
    var games = JSON.parse(b64_to_utf8((await getFile("data/games.json", octokit)).data.content));
    var player_stats = JSON.parse(b64_to_utf8((await getFile("data/player_stats.json", octokit)).data.content));

    for(var i=0; i < games.length; ++i){
        games[i].scores.forEach(function(elem){
            recalculateStats(elem.player, games[i], player_stats)
        });
    }

    await updateFile(JSON.stringify(player_stats, null, 2), "data/player_stats.json", "Added player_stats to data/data.json", octokit)
    console.log("Updated stats");
}

window.createSite = async function(password, game_id){
    //Check password
    if (CryptoJS.SHA256(password).toString() != password_hash){
        console.log("Wrong password");
        return;
    }

    // Get GitHub api token by decrypting encrypted token
    var decrypted = CryptoJS.AES.decrypt(encrypted_token, password);

    // Init Octokit
    const octokit = new Octokit({auth: decrypted.toString(CryptoJS.enc.Utf8)});

    // Add all players to the database
    var games = JSON.parse(b64_to_utf8((await getFile("data/games.json", octokit)).data.content));

    try{
        var entry = games.find(game => game.id == game_id);
        var game_site = generateGameSite(entry);
        var response = await pushFile(game_site, `games/${entry.id}.html`, `Created site for game "${entry.name}" id:${entry.id}`, octokit);
        console.log(`Created site ${entry.id}.html`);
    } catch (error){
        console.log(error);
    }

}

// Script for recreating all sites, requires password to be set to the password
window.createAllSites = async function (password){
    //Check password
    if (CryptoJS.SHA256(password).toString() != password_hash){
        console.log("Wrong password");
        return;
    }
    
    // Get GitHub api token by decrypting encrypted token
    var decrypted = CryptoJS.AES.decrypt(encrypted_token, password);

    // Init Octokit
    const octokit = new Octokit({auth: decrypted.toString(CryptoJS.enc.Utf8)});

    // Add all players to the database
    var data = await getFile("data/games.json", octokit);
    var games = JSON.parse(b64_to_utf8(data.data.content));

    for(var i=0; i < games.length; ++i){
        try{
            var entry = games[i];
            var game_site = generateGameSite(entry);
            var response = await pushFile(game_site, `games/${entry.id}.html`, `Created site for game "${entry.name}" id:${entry.id}`, octokit);
            console.log(`Created site ${games[i].id}.html`);
        } catch (error){
            console.log(error);
        }
    }
}

// Script for recreating all sites, requires password to be set to the password
window.recalculateAllEloScores = async function (password){
    //Check password
    if (CryptoJS.SHA256(password).toString() != password_hash){
        console.log("Wrong password");
        return;
    }
    
    // Get GitHub api token by decrypting encrypted token
    var decrypted = CryptoJS.AES.decrypt(encrypted_token, password);

    // Init Octokit
    const octokit = new Octokit({auth: decrypted.toString(CryptoJS.enc.Utf8)});

    // Add all players to the database
    var player_stats = JSON.parse(b64_to_utf8((await getFile("data/player_stats.json", octokit)).data.content));
    var games = JSON.parse(b64_to_utf8((await getFile("data/games.json", octokit)).data.content));

    var elo_params = calculateEloParameters(games);

    // Reset elo values
    Object.keys(player_stats).forEach(key => player_stats[key].elo = {value: 1000, history:[]});

    games.forEach(function(game){
        var means = elo_params.means;
        var stddevs = elo_params.stddev;
        var changes = calculateEloChange(game, player_stats, means, stddevs)
        game.scores
            .map(score=>score.player)
            .forEach(function(player){
                var change = changes[player].change;
                player_stats[player].elo.value += change;
                player_stats[player].elo.history.push({id:game.id, value:player_stats[player].elo.value, change:change});
            });
    });
    
    await updateFile(JSON.stringify(player_stats, null, 2), "data/player_stats.json", `Updated elo "${name}"`, octokit);
}

