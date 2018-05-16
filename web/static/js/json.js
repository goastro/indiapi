function getJSON(url, success, error) {
    $.ajax({
        type: 'GET',
        url: url,
        dataType: 'json',
        async: false
    }).done(success).fail(error);
}

function postJSON(url, data, success, error) {
    $.ajax({
        type: 'POST',
        url: url,
        dataType: 'json',
        data: JSON.stringify(data),
        contentType: 'application/json',
        async: false
    }).done(success).fail(error);
}
