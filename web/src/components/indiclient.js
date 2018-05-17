import React from 'react';
import { Accordion, Icon, Tab } from 'semantic-ui-react';
import Device from './device.js';

class INDIClient extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            devices: []
        };
    }

    componentDidMount() {
        const {clientId} = this.props;

        fetch('http://localhost:8080/api/' + clientId + '/devices')
            .then(res => res.json())
            .then((devices) => {
                console.log(devices);
                this.setState({
                    devices: devices
                });
            }, (error) => {
                this.setState({
                    devices: []
                });
            });
    }

    render() {
        const { active, index, onClick } = this.props;
        const { devices } = this.state;

        const panes = devices.map((device) => {
            return {
                menuItem: device.name,
                render: () => <Tab.Pane><Device device={device} /></Tab.Pane>
            }
        });

        return (
            <div>
                <Accordion.Title active={active} index={index} onClick={onClick}>
                    <Icon name='dropdown' />
                    INDI Client
                </Accordion.Title>
                <Accordion.Content active={active}>
                    <Tab panes={panes}  />
                </Accordion.Content>
            </div>);
    }
}

export default INDIClient;
