import React from 'react';
import { Accordion, Container } from 'semantic-ui-react'
import INDIClient from './indiclient';

class INDIClients extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            activeKey: null,
            clients: [],
            isLoaded: false,
            error: null
        };

    }

    getClients() {
        fetch('http://localhost:8080/api/clients')
            .then(res => res.json())
            .then((result) => {
                this.setState({
                    isLoaded: true,
                    clients: result.clients
                });
            }, (error) => {
                this.setState({
                    isLoaded: true,
                    clients: [],
                    error: error
                });
            });

            setTimeout(this.getClients.bind(this), 30000);
        }

    componentDidMount() {
        this.getClients();
    }

    handleClick(c) {
        const { clientId } = c;
        const { activeKey } = this.state;

        const newKey = (activeKey === clientId) ? null : clientId;



        this.setState({ activeKey: newKey });
    }

    render() {
        const { activeKey, clients } = this.state
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
            </Container>
        );
    }
}

export default INDIClients;
