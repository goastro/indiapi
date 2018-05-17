import React from 'react';
import { Checkbox, Input, Form } from 'semantic-ui-react';

class Value extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
        };
    }

    render() {
        const {value, type, readOnly} = this.props;

        var valueComponent = null;

        if (type === 'text') {
            valueComponent = <Form.Input label={value.label} type='text' defaultValue={value.value} readOnly={readOnly} />;
        } else if (type === 'switch') {
            valueComponent = <Form.Checkbox label={value.label} toggle checked={value === "On"} />;
        } else if (type === 'number') {
            valueComponent = <Form.Input label={value.label} type='text' defaultValue={value.value} readOnly={readOnly} />;
        } else {
            valueComponent = value.value;
        }

        return (
            <Form.Field>
                {valueComponent}
            </Form.Field>
        );
    }
}

export default Value;
