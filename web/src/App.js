import React, { Component } from 'react';
import { Segment, Header, Container } from 'semantic-ui-react';
import MainContent from './components/main.js';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Segment basic inverted>
          <Header as='h1'>INDI Web UI</Header>
        </Segment>
        <Container>
          <MainContent />
          </Container>
      </div>
    );
  }
}

export default App;
