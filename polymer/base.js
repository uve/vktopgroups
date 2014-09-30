// Copyright 2012 Google Inc. All Rights Reserved.

/**
 * @fileoverview
 * Provides methods for the TicTacToe sample UI and interaction with the
 * TicTacToe API.
 *
 * @author danielholevoet@google.com (Dan Holevoet)
 */

/** google global namespace for Google projects. */
var google = google || {};

/** TicTacToe namespace for this sample. */
google.projects = google.projects || {};

/**
 * Client ID of the application (from the APIs Console), e.g.
 * 123.apps.googleusercontent.com
 * @type {string}
 */
google.projects.CLIENT_ID =
    'YOUR-CLIENT-ID';

/**
 * Scopes used by the application.
 * @type {string}
 */
google.projects.SCOPES =
    'https://www.googleapis.com/auth/userinfo.email';

/**
 * Response type of the auth token.
 * @type {string}
 */
google.projects.RESPONSE_TYPE = 'token id_token';

/**
 * Status for an unfinished game.
 * @type {number}
 */
google.projects.NOT_DONE = 0;

/**
 * Status for a victory.
 * @type {number}
 */
google.projects.WON = 1;

/**
 * Status for a loss.
 * @type {number}
 */
google.projects.LOST = 2;

/**
 * Status for a tie.
 * @type {number}
 */
google.projects.TIE = 3;

/**
 * Strings for each numerical status.
 * @type {Array.number}
 */
google.projects.STATUS_STRINGS = [
    "NOT_DONE",
    "WON",
    "LOST",
    "TIE"
];

/**
 * Whether or not the user is signed in.
 * @type {boolean}
 */
google.projects.signedIn = false;

/**
 * Whether or not the game is waiting for a user's move.
 * @type {boolean}
 */
google.projects.waitingForMove = true;

/**
 * Loads the application UI after the user has completed auth.
 */
google.projects.userAuthed = function() {
  var request = gapi.client.oauth2.userinfo.get().execute(function(resp) {
    if (!resp.code) {
      var token = gapi.auth.getToken();
      // Use id_token instead of bearer token
      token.access_token = token.id_token;
      gapi.auth.setToken(token);
      google.projects.signedIn = true;

      console.log("Auth Token: " + token.access_token);
      
      /*
      document.getElementById('userLabel').innerHTML = resp.email;
      document.getElementById('signinButton').innerHTML = 'Sign out';
      */
      /*
      google.projects.setBoardEnablement(true);      
      google.projects.queryScores();
      */


      google.projects.loginSuccess(resp);
      //google.projects.loginSuccess(resp);
      

      //google.projects.ProjectsList();
    }
  });
};

/**
 * Handles the auth flow, with the given value for immediate mode.
 * @param {boolean} mode Whether or not to use immediate mode.
 * @param {Function} callback Callback to call on completion.
 */
google.projects.signin = function(mode, callback) {
  gapi.auth.authorize({client_id: google.projects.CLIENT_ID,
      scope: google.projects.SCOPES, immediate: mode,
      response_type: google.projects.RESPONSE_TYPE},
      callback);
};

/**
 * Presents the user with the authorization popup.
 */
google.projects.auth = function() {
  if (!google.projects.signedIn) {
    google.projects.signin(false,
        google.projects.userAuthed);
  } else {
    google.projects.signedIn = false;
    /*
    document.getElementById('userLabel').innerHTML = '(not signed in)';
    document.getElementById('signinButton').innerHTML = 'Sign in';
    */
    google.projects.setBoardEnablement(false);
  }
};

/**
 * Handles a square click.
 * @param {MouseEvent} e Mouse click event.
 */
google.projects.clickSquare = function(e) {
  if (google.projects.waitingForMove) {
    var button = e.target;
    button.innerHTML = 'X';
    button.removeEventListener('click', google.projects.clickSquare);
    google.projects.waitingForMove = false;

    var boardString = google.projects.getBoardString();
    var status = google.projects.checkForVictory(boardString);
    if (status == google.projects.NOT_DONE) {
      google.projects.getComputerMove(boardString);
    } else {
      google.projects.handleFinish(status);
    }
  }
};

/**
 * Resets the game board.
 */
google.projects.resetGame = function() {
  var buttons = document.querySelectorAll('td');
  for (var i = 0; i < buttons.length; i++) {
    var button = buttons[i];
    button.removeEventListener('click', google.projects.clickSquare);
    button.addEventListener('click', google.projects.clickSquare);
    button.innerHTML = '-';
  }
  document.getElementById('victory').innerHTML = '';
  google.projects.waitingForMove = true;
};

/**
 * Gets the computer's move.
 * @param {string} boardString Current state of the board.
 */
google.projects.getComputerMove = function(boardString) {
  gapi.client.tictactoe.board.getmove({'state': boardString}).execute(
      function(resp) {
    google.projects.setBoardFilling(resp.state);
    var status = google.projects.checkForVictory(resp.state);
    if (status != google.projects.NOT_DONE) {
      google.projects.handleFinish(status);
    } else {
      google.projects.waitingForMove = true;
    }
  });
};

/**
 * Sends the result of the game to the server.
 * @param {number} status Result of the game.
 */
google.projects.sendResultToServer = function(status) {
  gapi.client.tictactoe.scores.insert({'outcome':
      google.projects.STATUS_STRINGS[status]}).execute(
      function(resp) {
    google.projects.queryScores();
  });
};

/**
 * Queries for results of previous games.
 */
google.projects.queryScores = function() {
  gapi.client.tictactoe.scores.list().execute(function(resp) {
    var history = document.getElementById('gameHistory');
    history.innerHTML = '';
    if (resp.items) {
      for (var i = 0; i < resp.items.length; i++) {
        var score = document.createElement('li');
        score.innerHTML = resp.items[i].outcome;
        history.appendChild(score);
      }
    }
  });
};

/**
 * Shows or hides the board and game elements.
 * @param {boolean} state Whether to show or hide the board elements.
 */
google.projects.setBoardEnablement = function(state) {
  if (!state) {
    document.getElementById('board').classList.add('hidden');
    document.getElementById('gameHistoryWrapper').classList.add('hidden');
    document.getElementById('warning').classList.remove('hidden');
  } else {
    document.getElementById('board').classList.remove('hidden');
    document.getElementById('gameHistoryWrapper').classList.remove('hidden');
    document.getElementById('warning').classList.add('hidden');
  }
};

/**
 * Sets the filling of the squares of the board.
 * @param {string} boardString Current state of the board.
 */
google.projects.setBoardFilling = function(boardString) {
  var buttons = document.querySelectorAll('td');
  for (var i = 0; i < buttons.length; i++) {
    var button = buttons[i];
    button.innerHTML = boardString.charAt(i);
  }
};

/**
 * Checks for a victory condition.
 * @param {string} boardString Current state of the board.
 * @return {number} Status code for the victory state.
 */
google.projects.checkForVictory = function(boardString) {
  var status = google.projects.NOT_DONE;

  // Checks rows and columns.
  for (var i = 0; i < 3; i++) {
    var rowString = google.projects.getStringsAtPositions(
        boardString, i*3, (i*3)+1, (i*3)+2);
    status |= google.projects.checkSectionVictory(rowString);

    var colString = google.projects.getStringsAtPositions(
      boardString, i, i+3, i+6);
    status |= google.projects.checkSectionVictory(colString);
  }

  // Check top-left to bottom-right.
  var diagonal = google.projects.getStringsAtPositions(boardString,
      0, 4, 8);
  status |= google.projects.checkSectionVictory(diagonal);

  // Check top-right to bottom-left.
  diagonal = google.projects.getStringsAtPositions(boardString, 2,
      4, 6);
  status |= google.projects.checkSectionVictory(diagonal);

  if (status == google.projects.NOT_DONE) {
    if (boardString.indexOf('-') == -1) {
      return google.projects.TIE;
    }
  }

  return status;
};

/**
 * Checks whether a set of three squares are identical.
 * @param {string} section Set of three squares to check.
 * @return {number} Status code for the victory state.
 */
google.projects.checkSectionVictory = function(section) {
  var a = section.charAt(0);
  var b = section.charAt(1);
  var c = section.charAt(2);
  if (a == b && a == c) {
    if (a == 'X') {
      return google.projects.WON;
    } else if (a == 'O') {
      return google.projects.LOST
    }
  }
  return google.projects.NOT_DONE;
};

/**
 * Handles the end of the game.
 * @param {number} status Status code for the victory state.
 */
google.projects.handleFinish = function(status) {
  var victory = document.getElementById('victory');
  if (status == google.projects.WON) {
    victory.innerHTML = 'You win!';
  } else if (status == google.projects.LOST) {
    victory.innerHTML = 'You lost!';
  } else {
    victory.innerHTML = 'You tied!';
  }
  google.projects.sendResultToServer(status);
}; 

/**
 * Gets the current representation of the board.
 * @return {string} Current state of the board.
 */
google.projects.getBoardString = function() {
  var boardStrings = [];
  var buttons = document.querySelectorAll('td');
  for (var i = 0; i < buttons.length; i++) {
    boardStrings.push(buttons[i].innerHTML);
  }
  return boardStrings.join('');
};

/**
 * Gets the values of the board at the given positions.
 * @param {string} boardString Current state of the board.
 * @param {number} first First element to retrieve.
 * @param {number} second Second element to retrieve.
 * @param {number} third Third element to retrieve.
 */
google.projects.getStringsAtPositions = function(boardString, first,
    second, third) {
  return [boardString.charAt(first),
          boardString.charAt(second),
          boardString.charAt(third)].join('');
};

/**
 * Initializes the application.
 * @param {string} apiRoot Root of the API's path.
 */
google.projects.init = function(apiRoot, client_id ) {
  // Loads the OAuth and Tic Tac Toe APIs asynchronously, and triggers login
  // when they have completed.

  google.projects.CLIENT_ID = client_id;


  var apisToLoad;
  var callback = function() {
    if (--apisToLoad == 0) {
      google.projects.signin(true,
          google.projects.userAuthed);
    }
  }

  console.log("apiRoot: " + apiRoot);

  apisToLoad = 2; // must match number of calls to gapi.client.load()
  //gapi.client.load('vktopgroups', 'v1', callback, apiRoot);
  //gapi.client.load('oauth2', 'v2', callback);
  

/*
  var buttons = document.querySelectorAll('td');
  for (var i = 0; i < buttons.length; i++) {
    var button = buttons[i];
    button.addEventListener('click', google.projects.clickSquare);
  }

  var reset = document.querySelector('#restartButton');
  reset.addEventListener('click', google.projects.resetGame);
  */
};


/**
 * Queries for results of previous games.
 */
google.projects.ProjectsList = function() {
  gapi.client.vktopgroups.projects.list().execute(function(resp) {

    console.log(resp);
    /*
    var history = document.getElementById('gameHistory');
    history.innerHTML = '';
    if (resp.items) {
      for (var i = 0; i < resp.items.length; i++) {
        var score = document.createElement('li');
        score.innerHTML = resp.items[i].outcome;
        history.appendChild(score);
      }
    }
    */
  });
};


