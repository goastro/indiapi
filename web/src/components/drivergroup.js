import React from 'react';
import { Dropdown, Label, Form } from 'semantic-ui-react'

class DriverGroup extends React.Component {
    render() {
        const {name, drivers} = this.props;

        const options = drivers.map((driver) => {
            return {
                key: driver.Driver + '|' + driver.Label,
                text: driver.Label,
                value: driver.Driver + '|' + driver.Label
            }
        });

        return (
            <Form.Field>
                <Form.Dropdown label={name} multiple selection search options={options} />
            </Form.Field>
        );
    }
}

export default DriverGroup;
