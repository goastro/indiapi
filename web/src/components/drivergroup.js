import React from 'react';
import { Dropdown, Label, Form } from 'semantic-ui-react'

class DriverGroup extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            activeDrivers: []
        };

        this.onChange.bind(this);
    }

    onChange(e, value) {
        const { activeDrivers } = this.state;
        const { onDriverAdded, onDriverRemoved } = this.props;

        for (var i = 0; i < activeDrivers.length; i++) {
            if (!value.includes(activeDrivers[i])) {
                const split = activeDrivers[i].split('|');
                onDriverRemoved(split[0], split[1]);
                activeDrivers.splice(i, 1);
                i--;
            }
        }

        for (var i = 0; i < value.length; i++) {
            if (!activeDrivers.includes(value[i])) {
                const split = value[i].split('|');
                onDriverAdded(split[0], split[1]);
                activeDrivers.push(value[i])
            }
        }

        this.setState({
            activeDrivers
        });
    }

    render() {
        const { name, drivers } = this.props;

        const options = drivers.map((driver) => {
            return {
                key: driver.Driver + '|' + driver.Label,
                text: driver.Label + ' (' + driver.Driver + ')',
                value: driver.Driver + '|' + driver.Label
            }
        });

        return (
            <Form.Field>
                <Form.Dropdown label={name} multiple selection search options={options} onChange={(e, { value }) => this.onChange(e, value)} />
            </Form.Field>
        );
    }
}

export default DriverGroup;
