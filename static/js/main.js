document.getElementById("numbersearch").addEventListener("submit", function() {
    var xmlhttp= window.XMLHttpRequest ?
        new XMLHttpRequest() : new ActiveXObject("Microsoft.XMLHTTP");

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == 4 && xmlhttp.status == 200)
            alert(xmlhttp.responseText);
    };

    var number = cleanAndValidateNumber(document.getElementById("number").value);
    var name = validateName(document.getElementById('fullname').value);
    xmlhttp.open("GET","/search/?pn=" + number + "&fn=" + name, true);
    xmlhttp.send();
});

function cleanAndValidateNumber(value) {
    var cleanedValue = value.match(/\d+/g);

    if (cleanedValue === null || len(cleanedValue) > 10) {
        alert("Please enter a valid phone number");
        return
    }
    return cleanedValue
}

function validateName(value) {
    // just ensure something is entered
    if (len(value) < 1) {
        alert("Please enter a name");
        return false
    }
    return true
}
