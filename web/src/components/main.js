import React, { Component } from 'react';
import { Container, Accordion, Icon, Divider } from 'semantic-ui-react'

import INDIServer from './indiserver.js'
import INDIClients from './indiclients.js'

class MainContent extends Component {
    constructor(props) {
        super(props);

        this.state = {
            activeIndex: 0
        };

        this.handleClick = this.handleClick.bind(this);
    }

    handleClick(e, titleProps) {
        const { index } = titleProps;
        const { activeIndex } = this.state;
        const newIndex = activeIndex === index ? -1 : index;

        this.setState({ activeIndex: newIndex });
    }

    render() {
        const { activeIndex } = this.state
        return (
            <Container>
            <Accordion fluid styled>
                <Accordion.Title active={activeIndex === 0} index={0} onClick={this.handleClick}>
                    <Icon name='dropdown' />
                    INDI Server
                </Accordion.Title>
                <Accordion.Content active={activeIndex === 0}>
                    <INDIServer />
                </Accordion.Content>
            </Accordion>
            <Divider />
            <INDIClients />
            </Container>
        );
    }
}

export default MainContent;
