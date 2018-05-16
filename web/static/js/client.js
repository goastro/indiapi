var pendingRefresh = null;

function checkClientsIn(delay) {
    clearTimeout(pendingRefresh);

    pendingRefresh = setTimeout(checkClients, delay);
}

function checkClients() {
    clearTimeout(pendingRefresh);

    getJSON('/api/clients', function (clients) {
        for (var i = 0; i < clients.clients.length; i++) {
            var c = clients.clients[i];

            if (c.connected) {
                getJSON('/api/' + c.clientId + '/devices', function (devices) {
                    renderClient(c.clientId, devices);
                });
            }
        }
    });

    pendingRefresh = setTimeout(checkClients, 10000);
}

function setSwitchValue(clientId, device, prop, value, switchValue) {
    postJSON('/api/' + clientId + '/devices/switches/set', {
        device: device,
        property: prop,
        switch: value,
        value: switchValue
    }, function() {
        checkClientsIn(1000);
    }, function(xhr) {
        console.log(xhr.jsonResponse.error);
    })
}