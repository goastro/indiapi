import React from 'react';
import { Header, Icon, Form, Divider, Container } from 'semantic-ui-react';
import Value from './value.js';

class Property extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
        };
    }

    render() {
        const {property} = this.props;

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
                    return <Value key={value[0]} value={value[1]} readOnly={property.permissions !== 'rw'} type={property.type} />
                })}
                <Divider />
            </div>
        );
    }
}

export default Property;