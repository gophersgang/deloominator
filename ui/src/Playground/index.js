//@flow
import React, { Component } from 'react';
import DocumentTitle from 'react-document-title';

import { Button, Container, Form, Grid } from 'semantic-ui-react';
import { gql, graphql } from 'react-apollo';

import _ from 'lodash';

import QueryResult from './QueryResult';

class PlaygroundPage extends Component {
  state: {
    selectedDataSource?: string;
    currentQuery?: string;
    dataSource?: string;
    query?: string;
    showResult: bool;
  }

  constructor() {
    super();
    this.state = {
      showResult: false
    };
  }

  dataSourcesOptions = (data) => {
    let dataSources = (data.dataSources || []);
    return _.sortBy(dataSources, ['name'], ['asc']).map((s) => ({name: s.name, text: s.name, value: s.name}));
  }

  handleDataSourcesChange = (e, { value }) => {
    this.setState({selectedDataSource: value});
  }

  handleRunClick = (e) => {
    e.preventDefault();
    this.setState({
      showResult: true,
      query: this.state.currentQuery,
      dataSource: this.state.selectedDataSource,
    });
  }

  handleQueryChange = (e, { value }) => {
    this.setState({currentQuery: value});
  }

  render() {
    return (
      <DocumentTitle title='Playground'>
        <Container>
          <Grid.Row>
            <Grid.Column>
              <Form>
                <Form.Group>
                  <Form.Dropdown
                    placeholder='Data Source'
                    search selection
                    onChange={this.handleDataSourcesChange}
                    options={this.dataSourcesOptions(this.props.data)} />
                  <Button icon='play' primary content='Run' disabled={!(this.state.selectedDataSource && this.state.currentQuery)} onClick={this.handleRunClick}/>
                </Form.Group>
                <Form.TextArea placeholder='Write your query here' value={this.state.currentQuery} onChange={this.handleQueryChange} />
              </Form>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              { this.state.showResult && <QueryResult source={this.state.dataSource} input={this.state.query} /> }
            </Grid.Column>
          </Grid.Row>
        </Container>
      </DocumentTitle>
    )
  }
}

const Query = gql`{dataSources { name }}`;

const Playground = graphql(Query)(PlaygroundPage);

export default Playground;
