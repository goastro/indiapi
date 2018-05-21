import React from 'react';
import { Segment, Button, Form, Divider } from 'semantic-ui-react'

import DriverGroup from './drivergroup.js';

class INDIServer extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            drivers: [],
            isLoaded: false,
            error: null
        };

        this.handleDriverAdded.bind(this);
        this.handleDriverRemoved.bind(this);
    }

    componentDidMount() {
        fetch('http://localhost:8080/api/server/drivers')
        .then(res => res.json())
        .then((result) => {
            this.setState({
                isLoaded: true,
                drivers: result
            });
        },(error) => {
            this.setState({
                isLoaded: true,
                drivers: [],
                error: error
            });
        });
    }

    startServer() {
        fetch('http://localhost:8080/api/server/start', {
            method: 'POST'
        });
    }

    stopServer() {
        fetch('http://localhost:8080/api/server/stop', {
            method: 'POST'
        });
    }

    startDriver(driver, name) {
        fetch('http://localhost:8080/api/server/drivers/start', {
            method: 'POST',
            body: JSON.stringify({
                driver: driver,
                name: name
            }),
            headers: {
                'content-type': 'application/json'
            },
            mode: 'cors'
        });
    }

    stopDriver(driver, name) {
        fetch('http://localhost:8080/api/server/drivers/stop', {
            method: 'POST',
            body: JSON.stringify({
                driver: driver,
                name: name
            }),
            headers: {
                'content-type': 'application/json'
            },
            mode: 'cors'
        });
    }

    handleDriverAdded(driver, name) {
        this.startDriver(driver, name);
    }

    handleDriverRemoved(driver, name) {
        this.stopDriver(driver, name);
    }

    render() {
        const { error, isLoaded, drivers } = this.state;

        if (error) {
            return <div>Error: {error.Error}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            return(
            <Segment>
                <Form>
                    <Form.Field>
                        <Button onClick={this.startServer}>Start Server</Button>
                    </Form.Field>
                    <Divider />
                {Object.entries(drivers).map(function(v) {
                    return <DriverGroup key={v[0]} name={v[0]} drivers={v[1].sort((a, b) => a.Label.localeCompare(b.Label))} onDriverAdded={this.startDriver} onDriverRemoved={this.stopDriver} />
                }, this)}
                </Form>
            </Segment>);
        }
    }
}

export default INDIServer;
