import React from 'react';
import { Container, Accordion, Divider, List } from 'semantic-ui-react'
import Group from './group.js';

function getGroups(device) {
    var groups = {};

    for (var x in device.textProperties) {
        if (!groups[device.textProperties[x].group]) {
            groups[device.textProperties[x].group] = [];
        }

        groups[device.textProperties[x].group].push({
            ...device.textProperties[x],
            type: 'text'
        })
    }

    for (var x in device.switchProperties) {
        if (!groups[device.switchProperties[x].group]) {
            groups[device.switchProperties[x].group] = [];
        }

        groups[device.switchProperties[x].group].push({
            ...device.switchProperties[x],
            type: 'switch'
        })
    }

    for (var x in device.numberProperties) {
        if (!groups[device.numberProperties[x].group]) {
            groups[device.numberProperties[x].group] = [];
        }

        groups[device.numberProperties[x].group].push({
            ...device.numberProperties[x],
            type: 'number'
        })
    }

    for (var x in device.blobProperties) {
        if (!groups[device.blobProperties[x].group]) {
            groups[device.blobProperties[x].group] = [];
        }

        groups[device.blobProperties[x].group].push({
            ...device.blobProperties[x],
            type: 'blob'
        })
    }

    for (var x in device.lightProperties) {
        if (!groups[device.lightProperties[x].group]) {
            groups[device.lightProperties[x].group] = [];
        }

        groups[device.lightProperties[x].group].push({
            ...device.lightProperties[x],
            type: 'light'
        })
    }

    return groups;
}

class Device extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            activeKey: null,
            groups: []
        };
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        const {device} = nextProps;

        const groups = getGroups(device);

        return {
            groups
        };
    }

    handleClick(groupName) {
        const { activeKey } = this.state;

        const newKey = (activeKey === groupName) ? null : groupName;

        this.setState({ activeKey: newKey });
    }

    render() {
        const { groups, activeKey } = this.state;
        const { onUpdate, clientId, device } = this.props;

        return (
            <Container>
                {Object.entries(groups).sort((a, b) => a[0].localeCompare(b[0])).map((group) => {
                    return <Group key={group[0]} active={activeKey === group[0]} name={group[0]} clientId={clientId} deviceName={device.name} properties={group[1]} onClick={() => this.handleClick(group[0]) } onUpdate={onUpdate} />
                })}
                <List>
                    {device.messages === null ? null : device.messages.map((m) => <List.Item>{m.timestamp} - {m.message}</List.Item>)}
                </List>
            </Container>
        );
    }
}

export default Device;
