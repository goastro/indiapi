import React from 'react';
import { Checkbox, Input, Form, Button } from 'semantic-ui-react';

class Value extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            clientId: '',
            id: '',
            deviceName: '',
            propName: '',
            readOnly: false,
            value: '',
            valueName: '',
            valueLabel: '',
            type: '',
        };
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        const {clientId, deviceName, propName, value, type, readOnly} = nextProps;

        return {
            clientId,
            id: clientId+deviceName+propName+value.name,
            deviceName,
            propName,
            readOnly,
            value: value.value,
            valueName: value.name,
            valueLabel: value.label,
            type,
        };
    }

    onTextChange(e) {
        const { deviceName, propName, valueName, id } = this.state;
        const { onUpdate } = this.props;

        const newValue = document.getElementById(id).value;

        onUpdate(deviceName, propName, valueName, newValue, 'text');
    }

    onSwitchChange(e, {checked}) {
        const { deviceName, propName, valueName } = this.state;
        const { onUpdate } = this.props;

        const newValue = checked ? "On" : "Off";

        this.setState({
            value: newValue,
        });

        onUpdate(deviceName, propName, valueName, newValue, 'switch');
    }

    render() {
        const { deviceName, propName, valueName, valueLabel, value, readOnly, id, type } = this.state;

        var valueComponent = null;

        if (type === 'text') {
            valueComponent = (
                <Form.Group>
                    <Form.Input label={valueLabel} type='text' defaultValue={value} readOnly={readOnly} id={id} />
                    {readOnly ? null : <Button onClick={(e) => this.onTextChange(e)}>Set</Button>}
                </Form.Group>
            );
        } else if (type === 'switch') {
            valueComponent = <Form.Checkbox label={valueLabel} toggle checked={value === "On"} onChange={this.onSwitchChange.bind(this)} />;
        } else if (type === 'number') {
            valueComponent = (
                <Form.Group>
                    <Form.Input label={valueLabel} type='number' defaultValue={value} readOnly={readOnly} id={id} />
                    {readOnly ? null : <Button onClick={(e) => this.onTextChange(e)}>Set</Button>}
                </Form.Group>
            );
        } else {
            valueComponent = value;
        }

        return (
            <Form.Field>
                {valueComponent}
            </Form.Field>
        );
    }
}

export default Value;
