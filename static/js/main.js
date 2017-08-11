document.getElementById("numbersearch").addEventListener("submit", function() {
    var xmlhttp= window.XMLHttpRequest ?
        new XMLHttpRequest() : new ActiveXObject("Microsoft.XMLHTTP");

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState === 4 && xmlhttp.status === 200)
            alert(xmlhttp.responseText);
    };
    var number = cleanAndValidateNumber(document.getElementById("pn").value);
    var name = validateName(document.getElementById('fn').value);
    var formData = new FormData();
    formData.append("pn", number);
    formData.append("fn", name);
    aler(formData);
    xmlhttp.open("GET","/search/", false);
    xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xmlhttp.send(formData);
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
