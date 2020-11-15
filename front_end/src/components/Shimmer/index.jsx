import React, { Component } from 'react';
import ContentContainer from './ContentContainer';
import AnimatedBox from './StyledPlaceholderBox/index';
import { Grid } from 'grommet';

export class Shimmer extends Component {
  render() {
    return (
      <ContentContainer
        title='Shimmer'
        content='Loading layout placeholder effect before the content is loaded.'
      >
        <Grid
          areas={[
            { name: 'nav', start: [0, 0], end: [0, 0] },
            { name: 'main', start: [1, 0], end: [1, 0] },
            { name: 'side', start: [2, 0], end: [2, 0] },
            { name: 'foot', start: [0, 1], end: [2, 1] },
          ]}
          columns={['small', 'flex', 'small']}
          rows={['medium', 'small']}
          gap='small'
        >
          <AnimatedBox gridArea='nav' background='light-2' />
          <AnimatedBox gridArea='main' background='light-2' />
          <AnimatedBox gridArea='side' background='light-2' />
          <AnimatedBox gridArea='foot' background='light-2' />
        </Grid>
      </ContentContainer>
    );
  }
}
