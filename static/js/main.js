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
        if (xmlhttp.readyState === 4) {
            document.getElementById("results-container").innerHTML = "<h3>Results:</h3>";
            if (xmlhttp.status === 200) {
                var jsonResponse = JSON.parse(xmlhttp.responseText);
                successHandler(jsonResponse);
            }
            else {
                errorHandler(xmlhttp.status);
            }
            document.getElementById("pn").value = "";
            document.getElementById("fn").value = "";
        }
    };
    var number =  document.getElementById("pn").value;
    if (!isValidNumber(number, "US")) {
        alert("Please enter a valid US phone number");
        return
    }
    var name = document.getElementById("fn").value;
    if (!isValidName(name)) {
        return
    }
    var params = "pn=" + cleanPhone(number) + "&fn=" + cleanName(name);
    xmlhttp.open("POST","/api/search/", true);
    xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xmlhttp.send(params);
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
    errorDiv.className += "error";
    errorDiv.innerHTML = err;
}

function successHandler(response) {
    var html = "";
    var resultsDiv = createDiv();

    if ("error" in response) {
        resultsDiv.className += "error";
        html += response.error;
    } else {
        html +=
            "<h4>Name and Number</h4>" +
            "<div class='result-labels'>Phone Number: </div><div class='result-values'>" + formatLocal("US", response.pn) + "</div><br/>" +
            "<div class='result-labels'>Full Name: </div><div class='result-values'>" + response.Matches[0].fn + "</div><br>" +
            "<h4>Address</h4>" +
            "<div class='result-labels'>Street: </div><div class='result-values'>" + response.Matches[0].Address.Street + "</div><br/>" +
            "<div class='result-labels'>City: </div><div class='result-values'>" + response.Matches[0].Address.City + "</div><br/>" +
            "<div class='result-labels'>State: </div><div class='result-values'>" + response.Matches[0].Address.State + "</div><br/>" +
            "<div class='result-labels'>Zip: </div><div class='result-values'>" + response.Matches[0].Address.Zip + "</div>";
    }
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

function cleanName(fullName) {
    return fullName.replace(/[^0-9a-z]/gi, '');
}
