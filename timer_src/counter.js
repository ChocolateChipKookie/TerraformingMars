var canvas;
var MAX_PLAYERS = 5;
var no_players = 3;

var clicks = [];
var start=false;
var pause = false;

var colors = ["#444444","#B90E0A","#11AA11","#3344AA","#CCCC11"];
var text_color = "#000"
var button_color = "#DDD"
var passed_color = "#999";
var background_color = "#1E1E1E";
var font = null;

var max_time = 5;
var gen_first = 0;
var generation = 1;
var active_player = 0;

var timers = [];
var players = null;


function parseTime(time){
    var negative = false;

    if (time < 0){
        time = -time;
        negative = true;
    }

    var hours = Math.floor(time/3600000)
    var minutes = Math.floor(time%3600000/60000)
    var seconds = Math.floor(time%60000/1000)

    return `${negative ? "-":""}${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
}

function timerSetup(){
    // Set start flag to true
    start = true;

    var unitWidth = 0.8 * windowWidth;
    var unitHeight = (windowHeight - (10*no_players) - 2 * 30 - 15) / (no_players + 3);
    
    // Create players
    players = new Array(no_players).fill().map(Object);
    players.forEach(function(elem){
        elem.time = max_time * 60 * 1000;
        elem.big_pass = false;
        elem.timer_start = Date.now();
    });

    // Function executed when the main timer is clicked
    function nextPlayer(){
        // If the turn is paused, do nothing
        if(pause) return;

        // Update the timer for the active player
        players[active_player].time -= Date.now() - players[active_player].timer_start;
        
        // Update the current active player
        do {
            // Increase the counter by one untill a player has not passed
            active_player = (active_player + 1) % no_players;
        } while(players[active_player].big_pass);
        // Set the timer start to now
        players[active_player].timer_start = Date.now();
        mainTimer.stroke = colors[active_player];
    }

    // Function executed when all the players have ended their turn
    function nextGen(){
        // Set the active and generation first player to the new value
        active_player = gen_first = (gen_first + 1) % no_players;
        // Set generation first player colour
        mainTimer.stroke = colors[active_player];
        // Increment generation
        generation += 1;

        // Reset button variables
        players.forEach(player=>player.big_pass = false);
        timers.forEach(timer => timer.color = button_color);
    }

    // Function to enter pause
    function enterPause(){
        pause = !pause;
        if(pause){
            // Entering pause
            players[active_player].time -= Date.now() - players[active_player].timer_start;
            pauseButton.text = "Unpause"
        }
        else{
            // Exiting pause
            players[active_player].timer_start = Date.now();
            pauseButton.text = "Pause"
        }
    }

    // Create a timer box for every player
    for(var i = 0; i < no_players; ++i){
        var timer = new Clickable();
        timer.locate(0.1 * windowWidth, 10 + (unitHeight + 10) * i);
        timer.resize(unitWidth, unitHeight);
        timer.color = button_color;
        timer.textColor = text_color;
        timer.strokeWeight = 10;
        timer.stroke = colors[i];
        timer.id = i;
        timer.textSize = 20;
        timer.textFont = font;
        timer.onPress = function(){
            // If paused, the player cannot pass
            if(pause) return;
            // If big pass, able to unpass
            if(players[this.id].big_pass){
                players[this.id].big_pass = false;
                this.color = button_color;
                return;
            }
            // Player can only pass when he is active
            if(active_player == this.id){
                // Update player to big pass
                players[this.id].big_pass = true;
                this.color = passed_color;
                // Check if the generation is finished
                var finished = players.reduce((acc, elem)=> acc && elem.big_pass, true);
                if (finished){
                    // Go to next generation and pause
                    enterPause();
                    nextGen();
                }
                else{
                    // Go to next player
                    nextPlayer();
                }
            }
        }
        timer.updateText = function () {
            var time_delta = 0;
            if ((this.id == active_player)&&!pause){
                time_delta = Date.now() - players[this.id].timer_start;
            }
            var time = parseTime(players[this.id].time - time_delta);
            this.text = this.id == gen_first ? `[${time}]` : time;
            return this;
        }

        timers.push(timer);
    }

    // Create the main timer box
    mainTimer = new Clickable();
    mainTimer.locate(0.1 * windowWidth, 40 + (unitHeight + 10) * no_players);
    mainTimer.resize(unitWidth, unitHeight * 2);
    mainTimer.stroke = colors[active_player];
    mainTimer.strokeWeight = 5;
    mainTimer.color = button_color;
    mainTimer.textColor = text_color;
    mainTimer.textSize = 40;
    mainTimer.textFont = font;
    // Update text to the current player time
    mainTimer.updateText = function(){
        var time_delta = !pause ? Date.now() - players[active_player].timer_start : 0;
        this.text = `${parseTime(players[active_player].time - time_delta)}\nGen: ${generation}`;
        return this;
    };
    // On press change player
    mainTimer.onPress = nextPlayer;

    // Create button for pause
    pauseButton = new Clickable();
    pauseButton.locate(0.1 * windowWidth, 70 + (unitHeight + 10) * no_players + 2*unitHeight);
    pauseButton.resize(unitWidth, unitHeight);
    pauseButton.textColor = text_color;
    pauseButton.color = button_color;
    pauseButton.text = 'Pause';
    pauseButton.onPress = enterPause;
    pauseButton.textSize = 40;
    pauseButton.textFont = font;
}

function colourRotate (index){
    var tmp = colors[MAX_PLAYERS-1];

    for(var i = MAX_PLAYERS-1; i >= index; --i){
        colors[i]=colors[i-1];        
    }

    colors[index]=tmp;

    for(var i = 0; i < MAX_PLAYERS; ++i){
        clicks[i].color = colors[i];        
    }
}

function setup() {
    // Set graphics globals
    createCanvas(windowWidth, windowHeight);
    frameRate(30);
    font = loadFont('../resources/fonts/Prototype.ttf');

    // The unit height is window height - 8 spaces 15 pixels tall / seven possible tiles
    var unitHeight = (windowHeight - 8 * 15) / 7;

    //Create player counter
    playersSelector = new Clickable();
    playersSelector.locate(0.1*windowWidth, 15);
    playersSelector.resize(0.35*windowWidth, unitHeight);
    playersSelector.color = button_color;
    playersSelector.text = no_players;
    playersSelector.textFont = font;
    playersSelector.textSize = 40;
    playersSelector.onPress = function () {
        no_players += 1;
        if (no_players > 5) 
            no_players = 1;
        this.text = no_players;
    };

    // Create time counter
    timeSelector = new Clickable();
    timeSelector.locate(0.55*windowWidth, 15);
    timeSelector.resize(0.35*windowWidth, unitHeight);
    timeSelector.color = button_color;
    timeSelector.text = max_time;
    timeSelector.textFont = font;
    timeSelector.textSize = 40;
    timeSelector.onPress = function () {
        max_time = max_time + 5;
        if(max_time > 120){
            max_time = 5;
        }
        this.text = max_time;
    }

    startButton = new Clickable();
    startButton.locate(0.1*windowWidth, unitHeight + 30);
    startButton.resize(0.8*windowWidth, unitHeight);
    startButton.updatePosition = function(){
        this.locate(0.1*windowWidth, 15 + (unitHeight + 15) * (no_players + 1));
        return this;
    }
    startButton.color = button_color;
    startButton.textScaled = true;
    startButton.text = "Start";
    startButton.textFont = font;
    startButton.textSize = 40;
    startButton.onPress = timerSetup;
  
    for(var i = 0; i < MAX_PLAYERS; ++i){
        var click = new Clickable();
        click.id = i;
        click.locate(0.1*windowWidth, 15 + (unitHeight + 15) * (i+1));
        click.resize(0.8*windowWidth, unitHeight);
        click.text = "";
        click.color = colors[i];
        click.onPress = function () {
            colourRotate(this.id);
        }
        clicks.push(click);
    }
}

function draw() {
    if (!start){
        background(background_color);
        playersSelector.draw();
        startButton.updatePosition().draw();
        timeSelector.draw();

        for(var i = 0; i < no_players; ++i){
            clicks[i].draw();
        }
    }
    else{
//        background(colors[active_player]);
        background(background_color);
        mainTimer.updateText().draw();
        pauseButton.draw();
        timers.forEach(timer => timer.updateText().draw());
    }
}
