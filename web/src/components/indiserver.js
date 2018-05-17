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
        console.log("start server")
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
                    return <DriverGroup key={v[0]} name={v[0]} drivers={v[1]} />
                })}
                </Form>
            </Segment>);
        }
    }
}

export default INDIServer;
