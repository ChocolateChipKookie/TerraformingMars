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

var externalData = null;

// Loads players and corporations from database
async function loadExternalData(){
    const octokit = new Octokit();
    // Fetch file from git
    var data = await getFile("data/games.json", octokit);
    // Decode base64 file content and parse json
    data = JSON.parse(b64_to_utf8(data.data.content));
    externalData = data;
    // Create datalist for HTML
    var datalist = data.map( game => `<option value="${game.name}">`).join('\n');
    document.getElementById("datalist-games").innerHTML = datalist;

    // Add buttons
    var buttons = data.map( game => `<button class="game-link-button" onclick="window.location='/games/${game.id}.html'">${game.name}</button>`).join('\n');
    document.getElementById("game-link-list").innerHTML = buttons;
}

loadExternalData();

window.redirectToGame = function(){
    var name = document.getElementById("game-search").value;
    var id = externalData.filter(elem => elem.name == name)[0].id;
    window.location = `/games/${id}.html`
}

//Add listener to password field
var game_search = document.getElementById("game-search");
game_search.addEventListener("keyup", function(event) {
    if (event.keyCode === 13) {
        event.preventDefault();
        redirectToGame();
    }
});
