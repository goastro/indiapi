import React from 'react';
import { Accordion, Container, Form } from 'semantic-ui-react'
import INDIClient from './indiclient';

class INDIClients extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            activeKey: null,
            clients: [],
            isLoaded: false,
            error: null,
            network: '',
            address: '',
        };

        this.pendingRequest = null;
    }

    connect(network, address) {
        fetch('http://localhost:8080/api/connect', {
            method: 'POST',
            body: JSON.stringify({
                network: network,
                address: address
            }),
            headers: {
                'content-type': 'application/json'
            },
            mode: 'cors'
        }).then(() => this.getClients());
    }

    getClients() {
        clearTimeout(this.pendingRequest);

        fetch('http://localhost:8080/api/clients')
            .then(res => res.json())
            .then((result) => {
                this.setState({
                    isLoaded: true,
                    clients: result.clients.sort((a, b) => a.clientId.localeCompare(b.clientId))
                });
            }, (error) => {
                this.setState({
                    isLoaded: true,
                    clients: [],
                    error: error
                });
            });

        this.pendingRequest = setTimeout(this.getClients.bind(this), 10000);
    }

    componentDidMount() {
        this.getClients();
    }

    handleSubmit() {
        const { network, address } = this.state;
        this.connect(network, address);
    }

    handleChange(e, { name, value }) {
        this.setState({ [name]: value });
    }

    handleClick(c) {
        const { clientId } = c;
        const { activeKey } = this.state;

        const newKey = (activeKey === clientId) ? null : clientId;

        this.setState({ activeKey: newKey });
    }

    render() {
        const { activeKey, clients } = this.state;
        console.log('render');
        return (
            <Container>
                <Accordion fluid styled>
                    {clients.map(function (c) {
                        if (!c.connected) {
                            return null;
                        }
                        return <INDIClient key={c.clientId} active={c.clientId === activeKey} clientId={c.clientId} onClick={() => this.handleClick(c)} />
                    }, this)}
                </Accordion>
                <Form onSubmit={this.handleSubmit.bind(this)}>
                    <Form.Group>
                    <Form.Input placeholder='tcp' name='network' label='Network' onChange={this.handleChange.bind(this)} />
                    <Form.Input placeholder='localhost:7624' name='address' label='Address' onChange={this.handleChange.bind(this)} />
                    </Form.Group>
                    <Form.Button content='Connect' />
                </Form>
            </Container>
        );
    }
}

export default INDIClients;
