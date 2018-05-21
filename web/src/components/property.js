import React from 'react';
import { Header, Icon, Form, Divider, Container, List } from 'semantic-ui-react';
import Value from './value.js';

class Property extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            clientId: '',
            deviceName: '',
            property: null,
        };
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        const { clientId, deviceName, property } = nextProps;

        return {
            clientId,
            deviceName,
            property,
        };
    }

    render() {
        const { clientId, deviceName, property } = this.state;
        const { onUpdate } = this.props;

        var color = 'grey';
        if (property.state === 'Warn') {
            color = 'yellow';
        } else if (property.state === 'Alert') {
            color = 'red';
        } else if (property.state === 'Ok') {
            color = 'green';
        }

        return (
            <div>
                <Header as='h5'><Icon name={'circle'} color={color} />{property.label}</Header>
                {Object.entries(property.values).map((value) => {
                    return <Value key={value[0]} value={value[1]} readOnly={property.permissions !== 'rw'} clientId={clientId} deviceName={deviceName} propName={property.name} type={property.type} onUpdate={onUpdate} />
                })}
                <Divider />
                <List>
                    {property.messages.map((m) => <List.Item>{m.timestamp} - {m.message}</List.Item>)}
                </List>
            </div>
        );
    }
}

export default Property;