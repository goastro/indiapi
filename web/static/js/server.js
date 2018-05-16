function startServer() {
    $.post("/api/server/start").done(function () {
        $('#server_message p').text("indiserver started successfully");
        $('#server_message').removeClass('warning').show(500);
    }).fail(function () {
        $('#server_message p').text("indiserver was unable to start. check the log.");
        $('#server_message').addClass('warning').show(500);
    });
}

function stopServer() {
    $.post("/api/server/stop").done(function () {
        $('#server_message p').text("indiserver stopped successfully");
        $('#server_message').removeClass('warning').show(500);
    }).fail(function () {
        $('#server_message p').text("indiserver was unable to stop. check the log.");
        $('#server_message').addClass('warning').show(500);
    });
}

function startDriver(driver, name) {
    postJSON('/api/server/drivers/start', {
        "driver": driver,
        "name": name
    }, function () {
        $('#server_message p').text("driver started successfully");
        $('#server_message').removeClass('warning').show(500);
    }, function () {
        $('#server_message p').text("driver was unable to start. check the log.");
        $('#server_message').addClass('warning').show(500);
    });
}

function stopDriver(driver, name) {
    postJSON('/api/server/drivers/stop', {
        "driver": driver,
        "name": name
    }, function () {
        $('#server_message p').text("driver stopped successfully");
        $('#server_message').removeClass('warning').show(500);
    }, function () {
        $('#server_message p').text("driver was unable to stop. check the log.");
        $('#server_message').addClass('warning').show(500);
    });
}
