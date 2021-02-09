window.utf8_to_b64 = function(str) {
    // first we use encodeURIComponent to get percent-encoded UTF-8,
    // then we convert the percent encodings into raw bytes which
    // can be fed into btoa.
    return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g,
        function toSolidBytes(match, p1) {
            return String.fromCharCode('0x' + p1);
    }));
}

window.b64_to_utf8 = function(str) {
    // Going backwards: from bytestream, to percent-encoding, to original string.
    return decodeURIComponent(atob(str).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
}

//Capitalizes given string
window.capitalize = function(str){
    return str.charAt(0).toUpperCase() + str.slice(1);
}

// Gets the child of a parent
// If the child doesnt exist in the parent node, returns false
window.getElementInsideContainer = function(containerID, childID) {
    var elm = document.getElementById(childID);
    var parent = elm ? elm.parentNode : {};
    return (parent.id && parent.id === containerID) ? elm : false;
}

//////////////
//    GIT
//////////////
window.owner = 'ChocolateChipKookie';
window.repo = "chocolatechipkookie.github.io";
window.encrypted_token = "U2FsdGVkX1+n1ehJgHqx60l9tKl1nu1zx0MlMiCXO+YDPnIW/5I0+1JboKey3qjNMo10biUocmSAMHrD0bwJ8Q==";
window.password_hash = "2127c97b1c21f675c8ea7c47ce5fffb827b15035aea988e525ab8a24fd8ad6d0"
//Loads git file
window.getFile = async function(filepath, octokit){
    return octokit.request(`GET https://api.github.com/repos/${owner}/${repo}/contents/${filepath}`);
}

//Checks if git file exists
window.checkFile = async function(filepath, octokit){
    try {
        getFile(filepath, octokit)
        return true;
    } catch (error) {
        if (error.status === 404) {
            return false;
        } else {
            throw error;
        }
    }

}

//Updates git file
window.updateFile = async function(content, filepath, message, octokit){
    var prev_blob = await getFile(filepath, octokit);
    
    const response = await octokit.request(`PUT https://api.github.com/repos/${owner}/${repo}/contents/${filepath}`, 
    {
        owner: owner,
        repo: repo,
        path: filepath,
        message: message,
        content: utf8_to_b64(content),
        sha: prev_blob.data.sha,
    });
    return response;
}

//Create git file
window.createFile = async function(content, filepath, message, octokit){
    const response = await octokit.request(`PUT https://api.github.com/repos/${owner}/${repo}/contents/${filepath}`,
    {
        owner: owner,
        repo: repo,
        path: filepath,
        message: message,
        content: utf8_to_b64(content),
    });
    return response;
}

//Create git file
window.pushFile = async function(content, filepath, message, octokit){
    if(checkFile(filepath, octokit)){
        return updateFile(content, filepath, message, octokit);
    }
    return createFile(content, filepath, message, octokit);
}

////////////
//    DATA
////////////

window.util = {
    "game_data":{
        "maps": [
            "Tharsis",
            "Hellas",
            "Elysium"
        ],
        "corporation_names": [
            "Beginner",
            "Credicor",
            "Ecoline",
            "Helion",
            "Interplanetary Cinematics",
            "Inventrix",
            "Mining Guild",
            "Phobolog",
            "Tharsis Republic",
            "Thorgate",
            "UNMI",
            "Saturn Systems",
            "Teractor",
            "Aridor",
            "Arklight",
            "Polyphemos",
            "Poseidon",
            "Stormcraft Incorporated",
            "Cheung Shing Mars",
            "Point Luna",
            "Robinson Industries",
            "Valley Trust",
            "Vitor",
            "Aphrodite",
            "Celestic",
            "Manutech",
            "Morning Star INC",
            "Viron",
            "Arcadian Communities",
            "Recyclon",
            "Splice",
            "Lakefront Resorts",
            "Pristar",
            "Septem Tribus",
            "Terralabs Research",
            "Utopia Invest",
            "Factorum",
            "Mons Insurance",
            "Philares"
          ],
        "milestones":{
            "Tharsis": [
                "Builder",
                "Gardener",
                "Mayor",
                "Planer",
                "Terraformer"
            ],
            "Hellas":[
                "Diversifier",
                "Tactician",
                "Polar Explorer",
                "Energizer",
                "Rim Settler"
            ],
            "Elysium":[
                "Generalist",
                "Specialist",
                "Ecologist",
                "Tycoon",
                "Legend"
            ],
            "Venus":[
                "Hoverlord"
            ],
            "all":[
                "Builder",
                "Gardener",
                "Mayor",
                "Planer",
                "Terraformer",
                "Diversifier",
                "Tactician",
                "Polar Explorer",
                "Energizer",
                "Rim Settler",
                "Generalist",
                "Specialist",
                "Ecologist",
                "Tycoon",
                "Legend",
                "Hoverlord"
            ]
          },
        "awards":{
            "Tharsis":[
                "Landlord",
                "Scientist",
                "Banker",
                "Thermalist",
                "Miner"
            ],
            "Hellas":[
                "Cultivator",
                "Magnate",
                "Space Baron",
                "Excentric",
                "Contractor"
            ],
            "Elysium":[
                "Celebrity",
                "Industrialist",
                "Desert Settler",
                "Estate Dealer",
                "Benefactor"
            ],
            "Venus":[
                "Venuphile"
            ],
            "all":[
                "Landlord",
                "Scientist",
                "Banker",
                "Thermalist",
                "Miner",
                "Cultivator",
                "Magnate",
                "Space Baron",
                "Excentric",
                "Contractor",
                "Celebrity",
                "Industrialist",
                "Desert Settler",
                "Estate Dealer",
                "Benefactor",
                "Venuphile"
            ]
        }
    }
};
