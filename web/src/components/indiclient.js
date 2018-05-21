import React from 'react';
import { Accordion, Icon, Tab } from 'semantic-ui-react';
import Device from './device.js';

class INDIClient extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            devices: []
        };

        this.pendingRequest = null;
    }

    componentDidMount() {
        this.getDevices();
    }

    getDevices(delay) {
        const { clientId } = this.props;

        clearTimeout(this.pendingRequest);

        if (delay) {
            this.pendingRequest = setTimeout(this.getDevices.bind(this), delay);
            return;
        }

        fetch('http://localhost:8080/api/' + clientId + '/devices')
            .then(res => res.json())
            .then((devices) => {
                this.setState({
                    devices: devices.sort((a, b) => a.name.localeCompare(b.name))
                });
            }, (error) => {
                this.setState({
                    devices: []
                });
            });

        this.pendingRequest = setTimeout(this.getDevices.bind(this), 10000);
    }

    updateValue(device, property, valueName, value, type) {
        const { clientId } = this.props;

        var endpoint = '';
        var body = null;

        if (type === 'text') {
            endpoint = 'texts';
            body = {
                device: device,
                property: property,
                text: valueName,
                value: value
            };
        } else if (type === 'switch') {
            endpoint = 'switches';
            body = {
                device: device,
                property: property,
                switch: valueName,
                value: value
            };
        } else if (type === 'number') {
            endpoint = 'numbers';
            body = {
                device: device,
                property: property,
                number: valueName,
                value: value
            };
        } else {
            return;
        }

        fetch('http://localhost:8080/api/' + clientId + '/devices/' + endpoint + '/set', {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {
                'content-type': 'application/json'
            },
            mode: 'cors'
        }).then(() => {
            this.getDevices(1000);
        });
    }

    render() {
        const { clientId, active, index, onClick } = this.props;
        const { devices } = this.state;

        const panes = devices.map((device) => {
            return {
                menuItem: device.name,
                render: () => <Tab.Pane><Device device={device} onUpdate={this.updateValue.bind(this)} clientId={clientId} /></Tab.Pane>
            }
        }, this);

        return (
            <div>
                <Accordion.Title active={active} index={index} onClick={onClick}>
                    <Icon name='dropdown' />
                    INDI Client
                </Accordion.Title>
                <Accordion.Content active={active}>
                    <Tab panes={panes} />
                </Accordion.Content>
            </div>
        );
    }
}

export default INDIClient;
