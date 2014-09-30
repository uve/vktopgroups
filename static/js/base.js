var mindale = mindale || {};


/**
 * Client ID of the application (from the APIs Console), e.g.
 * 123.apps.googleusercontent.com
 * @type {string}
 */

if (window.location.host == "mindale.com"){

    mindale.CLIENT_ID = "882975820932-q34i2m1lklcmv8kqqrcleumtdhe4qbhk.apps.googleusercontent.com"
    mindale.HOST = "https://mindale-com.appspot.com/_ah/api";
}
else{

    mindale.CLIENT_ID = "882975820932-ltfu1oa31f80o4v9tqfp8k5ghe45oibu.apps.googleusercontent.com";
    mindale.HOST = "http://localhost:8080/_ah/api";
}


function init() {
    mindale.init(mindale.HOST);
}

/**
 * Scopes used by the application.
 * @type {string}
 */
mindale.SCOPES =
    'https://www.googleapis.com/auth/userinfo.email';


/**
 * Whether or not the user is signed in.
 * @type {boolean}
 */
mindale.signedIn = false;


/**
 * Response type of the auth token.
 * @type {string}
 */
mindale.RESPONSE_TYPE = 'token id_token';

/**
 * Status for an unfinished game.
 * @type {number}
 */
mindale.NOT_DONE = 0;

/**
 * Status for a victory.
 * @type {number}
 */
mindale.WON = 1;

/**
 * Status for a loss.
 * @type {number}
 */
mindale.LOST = 2;

/**
 * Status for a tie.
 * @type {number}
 */
mindale.TIE = 3;

/**
 * Strings for each numerical status.
 * @type {Array.number}
 */
mindale.STATUS_STRINGS = [
    "NOT_DONE",
    "WON",
    "LOST",
    "TIE"
];


/**
 * Whether or not the game is waiting for a user's move.
 * @type {boolean}
 */
mindale.waitingForMove = true;

/**
 * Loads the application UI after the user has completed auth.
 */
mindale.userAuthed = function() {
  var request = gapi.client.oauth2.userinfo.get().execute(function(resp) {
    if (!resp.code) {
      var token = gapi.auth.getToken();
      // Use id_token instead of bearer token
      //token.access_token = token.id_token;
      gapi.auth.setToken(token);
      mindale.signedIn = true;
      document.getElementById('userLabel').innerHTML = resp.email;
      document.getElementById('signinButton').innerHTML = 'Sign out';
      //mindale.setBoardEnablement(true);
      //mindale.queryScores();

      mindale.payments();

    }
  });
};

/**
 * Handles the auth flow, with the given value for immediate mode.
 * @param {boolean} mode Whether or not to use immediate mode.
 * @param {Function} callback Callback to call on completion.
 */
mindale.signin = function(mode, callback) {
  gapi.auth.authorize({client_id: mindale.CLIENT_ID,
      scope: mindale.SCOPES, immediate: mode,
      response_type: mindale.RESPONSE_TYPE},
      callback);
};

/**
 * Presents the user with the authorization popup.
 */
mindale.auth = function() {
  if (!mindale.signedIn) {
    mindale.signin(false,
        mindale.userAuthed);
  } else {
    mindale.signedIn = false;
    document.getElementById('userLabel').innerHTML = '(not signed in)';
    document.getElementById('signinButton').innerHTML = 'Sign in';



    //mindale.setBoardEnablement(false);
  }


};



/**
 * Queries for results of previous games.
 */
mindale.payments = function() {

    console.log('Loading payments');


    gapi.client.mindale.scores.list().execute(function(resp) {


        console.log(resp);
    });
};



/**
 * Handles a square click.
 * @param {MouseEvent} e Mouse click event.
 */
mindale.clickSquare = function(e) {
  if (mindale.waitingForMove) {
    var button = e.target;
    button.innerHTML = 'X';
    button.removeEventListener('click', mindale.clickSquare);
    mindale.waitingForMove = false;

    var boardString = mindale.getBoardString();
    var status = mindale.checkForVictory(boardString);
    if (status == mindale.NOT_DONE) {
      mindale.getComputerMove(boardString);
    } else {
      mindale.handleFinish(status);
    }
  }
};

/**
 * Resets the game board.
 */
mindale.resetGame = function() {
  var buttons = document.querySelectorAll('td');
  for (var i = 0; i < buttons.length; i++) {
    var button = buttons[i];
    button.removeEventListener('click', mindale.clickSquare);
    button.addEventListener('click', mindale.clickSquare);
    button.innerHTML = '-';
  }
  document.getElementById('victory').innerHTML = '';
  mindale.waitingForMove = true;
};

/**
 * Gets the computer's move.
 * @param {string} boardString Current state of the board.
 */
mindale.getComputerMove = function(boardString) {
  gapi.client.mindale.board.getmove({'state': boardString}).execute(
      function(resp) {
    mindale.setBoardFilling(resp.state);
    var status = mindale.checkForVictory(resp.state);
    if (status != mindale.NOT_DONE) {
      mindale.handleFinish(status);
    } else {
      mindale.waitingForMove = true;
    }
  });
};

/**
 * Sends the result of the game to the server.
 * @param {number} status Result of the game.
 */
mindale.sendResultToServer = function(status) {
  gapi.client.mindale.scores.insert({'outcome':
      mindale.STATUS_STRINGS[status]}).execute(
      function(resp) {
    mindale.queryScores();
  });
};

/**
 * Queries for results of previous games.
 */
mindale.queryScores = function() {
  gapi.client.mindale.scores.list().execute(function(resp) {
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
mindale.setBoardEnablement = function(state) {
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
mindale.setBoardFilling = function(boardString) {
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
mindale.checkForVictory = function(boardString) {
  var status = mindale.NOT_DONE;

  // Checks rows and columns.
  for (var i = 0; i < 3; i++) {
    var rowString = mindale.getStringsAtPositions(
        boardString, i*3, (i*3)+1, (i*3)+2);
    status |= mindale.checkSectionVictory(rowString);

    var colString = mindale.getStringsAtPositions(
      boardString, i, i+3, i+6);
    status |= mindale.checkSectionVictory(colString);
  }

  // Check top-left to bottom-right.
  var diagonal = mindale.getStringsAtPositions(boardString,
      0, 4, 8);
  status |= mindale.checkSectionVictory(diagonal);

  // Check top-right to bottom-left.
  diagonal = mindale.getStringsAtPositions(boardString, 2,
      4, 6);
  status |= mindale.checkSectionVictory(diagonal);

  if (status == mindale.NOT_DONE) {
    if (boardString.indexOf('-') == -1) {
      return mindale.TIE;
    }
  }

  return status;
};

/**
 * Checks whether a set of three squares are identical.
 * @param {string} section Set of three squares to check.
 * @return {number} Status code for the victory state.
 */
mindale.checkSectionVictory = function(section) {
  var a = section.charAt(0);
  var b = section.charAt(1);
  var c = section.charAt(2);
  if (a == b && a == c) {
    if (a == 'X') {
      return mindale.WON;
    } else if (a == 'O') {
      return mindale.LOST
    }
  }
  return mindale.NOT_DONE;
};

/**
 * Handles the end of the game.
 * @param {number} status Status code for the victory state.
 */
mindale.handleFinish = function(status) {
  var victory = document.getElementById('victory');
  if (status == mindale.WON) {
    victory.innerHTML = 'You win!';
  } else if (status == mindale.LOST) {
    victory.innerHTML = 'You lost!';
  } else {
    victory.innerHTML = 'You tied!';
  }
  mindale.sendResultToServer(status);
};

/**
 * Gets the current representation of the board.
 * @return {string} Current state of the board.
 */
mindale.getBoardString = function() {
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
mindale.getStringsAtPositions = function(boardString, first,
    second, third) {
  return [boardString.charAt(first),
          boardString.charAt(second),
          boardString.charAt(third)].join('');
};

/**
 * Initializes the application.
 * @param {string} apiRoot Root of the API's path.
 */
mindale.init = function(apiRoot) {
  // Loads the OAuth and Tic Tac Toe APIs asynchronously, and triggers login
  // when they have completed.
  var apisToLoad;
  var callback = function() {
    if (--apisToLoad == 0) {
      mindale.signin(true,
          mindale.userAuthed);
    }
  }

  apisToLoad = 2; // must match number of calls to gapi.client.load()
  gapi.client.load('mindale', 'v1', callback, apiRoot);
  gapi.client.load('oauth2', 'v2', callback);



    /*
  var buttons = document.querySelectorAll('td');
  for (var i = 0; i < buttons.length; i++) {
    var button = buttons[i];
    button.addEventListener('click', mindale.clickSquare);
  }
*/
   /*
  var reset = document.querySelector('#restartButton');
  reset.addEventListener('click', mindale.resetGame);
  */
};
