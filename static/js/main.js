"use strict";

var form = document.getElementById("search-form");
form.addEventListener('submit', function (event) {
    event.preventDefault();
    sendData();
});

function sendData() {
    var xmlhttp = window.XMLHttpRequest ?
        new XMLHttpRequest() : new ActiveXObject("Microsoft.XMLHTTP");

    xmlhttp.onreadystatechange = function() {
        console.error("onreadystate");
        if (xmlhttp.readyState === 4) {
            console.error("ready state");
            console.error(xmlhttp.responseText);
            console.error("any request");

            console.error(jsonResponse);
            var containerHTML = "<h3>Results:</h3>";
            document.getElementById("results-container").innerHTML = containerHTML;
            if (xmlhttp.status === 200) {
                var jsonResponse = JSON.parse(xmlhttp.responseText);
                successHandler(jsonResponse);
            }
            else {
                errorHandler(xmlhttp.status);
            }
        }
    };
    var number =  document.getElementById("pn").value;
    console.error(number);
    if (!isValidNumber(number, "US")) {
        alert("Please enter a valid US phone number");
        return
    }
    var name = document.getElementById("fn").value;
    if (!isValidName(name)) {
        return
    }
    var params = "pn=" + number + "&fn=" + name;
    console.error(params);
    xmlhttp.open("POST","/api/search/", true);
    console.error("post");
    xmlhttp.setRequestHeader("Content-Type", "application/json");
    // xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    console.error("request header");
    xmlhttp.send(params);
    console.error("params sent");
}

// creates a div to be append to the main results container div
function createDiv() {
    var newDiv = document.createElement("div");
    var containerDiv = document.getElementById("results-container");
    containerDiv.appendChild(newDiv);

    return newDiv
}

function errorHandler(err) {
    var errorDiv = createDiv();
    var html = "<h3>" + err + "</h3>";
    errorDiv.innerHTML = html;
}

function successHandler(response) {
    var html = "";
    if ("error" in response) {
        html += "<p>" + response.error + "</p>";
        console.error(html);
    } else {
        html +=
            "<p>" + response.pn + "</p>" +
            '<p>' + response.name + '</p><br>' +
            '<p>' + response.Matches[0].Address.street + '</p>' +
            '<p>' + response.Matches[0].Address.city + '</p>' +
            '<p>' + response.Matches[0].Address.state + '</p>' +
            '<p>' + response.Matches[0].Address.zip + '</p>';
    }
    var resultsDiv = createDiv();
    resultsDiv.innerHTML = html;
}

function isValidName(fullName) {
    // just ensure something is entered
    if (fullName.length < 1) {
        alert("Please enter a name");
        return false
    }
    return true
}
