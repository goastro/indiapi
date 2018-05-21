import React from 'react';
import { Accordion, Icon, Form, Header, Segment } from 'semantic-ui-react'
import Property from './property.js'

class Group extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            name: null,
            properties: []
        };
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        const { properties, name } = nextProps;

        return {
            name,
            properties
        };
    }

    render() {
        const { name, properties } = this.state;
        const { active, onClick, onUpdate, clientId, deviceName } = this.props;

        return (
            <Segment>
                <Header as='h3'>{name}</Header>
                <Form>
                    {properties.map((property) => {
                        return <Property key={property.name} clientId={clientId} deviceName={deviceName} property={property} onUpdate={onUpdate} />
                    })}
                </Form>
            </Segment>
        );
    }
}

export default Group;
