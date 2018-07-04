"use strict";

var members = {};
/**
 * Random name, generated on client start
 * @type {string}
 */
var myName = null;
var myId = null;
/**
 * @type {WsClient}
 */
var client = null;
var syncTimeTimer = null;
var map = null;
var gameTime = 0;
var lastServerState = null;

function onmessage(messageType, data) {
    switch (messageType) {
        case MessageType.TEXT:
            if (data.sender === 0) {
                showMessage(data.text, "text-primary");
            }
            else if (data.Sender === myId) {
                showMessage(data.text, "text-success");
            }
            else {
                var username = data.sender in members ? members[data.sender].name : 'Unknown sender';
                showMessage(username + ": " + data.text);
            }
            break;
        case MessageType.WELCOME:
            myId = data.id;
            newMember(myId, myName);
            break;
        case MessageType.USER_LIST:
            if ('users' in data && data.users != null) {
                for (var i = 0, user; user = data.users[i]; i++) {
                    newMember(user.id, user.name);
                }
            }
            break;
        case MessageType.USER_LOGGEDIN:
            showMessage(data.name + " logged in", "text-muted");
            newMember(data.id, data.name);
            break;
        case MessageType.USER_LOGGEDOUT:
            removeMember(data.id);
            break;
        case MessageType.SYNC_OBJECTS_POSITIONS:
            updateObjectsPositions(data.positions, data.time);
            if (!map.isAnimating) {
                map.forceDraw();
            }
            break;
        case MessageType.ERROR:
            showMessage('Error: ' + data.description);
            break;
        case MessageType.CHANGE_SIMULATION_MODE:
            // сервер вообще не должен присылать такого сообщения
            console.warn("TODO");
            break;
        case MessageType.SERVER_STATE:
            map.updateServerState(data);
            break;
        default:
            for (var key in MessageType) {
                if (MessageType.hasOwnProperty(key)) {
                    if (MessageType[key] === messageType) {
                        showMessage('Not implemented ' + key.toUpperCase());
                        return;
                    }
                }
            }
            showMessage("Unknown message type: " + messageType + data, "text-danger");
            break;
    }
}

function onclose() {
    if (syncTimeTimer != null) {
        clearInterval(syncTimeTimer);
        syncTimeTimer = null;
    }

    $('.chat_members').empty();
    members = {};
    map.clear();
    showMessage('disconnected');
}

function updateObjectsPositions(positions, time) {
    for (var objectId in positions) {
        if (positions.hasOwnProperty(objectId)) {
            map.updateObjectPosition(positions[objectId], time);
        }
    }
}

/**
 *
 * @param {Number} id
 * @param {String} name
 * @returns {Player}
 */
function newMember(id, name) {
    if (!(id in members)) {
        // console.log('adding player #' + id + ' (' + name + ')');
        var member = new Player(id, name);
        member.anchor = $('<div class="member" aria-hidden="true" data-id="' + id + '">' + name + '</div>');
        $('.chat_members').append(member.anchor);
        if (id === myId) {
            member.isMe = true;
            member.anchor.css('font-weight', 'bold');
        } else {
            member.isMe = false;
        }
        members[id] = member;

        map.addPlayer(member);

        return member;
    }
    else {
        return members[id];
    }
}

function removeMember(id) {
    if (id in members) {
        showMessage(members[id].name + " logged out");
        members[id].anchor.remove();
        map.removePlayer(id);
        delete members[id];
    }
}

function showMessage(text, messageType) {
    if (typeof messageType === 'undefined' || messageType == null) {
        messageType = "";
    }

    $('.chat_window').append('<div class="message ' + messageType + '">' + text + '</div>');
}

/**
 *
 * @param id
 * @param name
 * @constructor
 */
var Player = function (id, name) {
    this.id = id;
    this.name = name;
    /**
     *
     * @type {boolean}
     */
    this.isMe = false;
    /**
     *
     * @type {HTMLElement}
     */
    this.anchor = null;

    this.state = {
        position: {x: 0, y: 0}
    };
};