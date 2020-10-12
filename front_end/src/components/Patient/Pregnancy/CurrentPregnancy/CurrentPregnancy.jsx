import { Box } from 'grommet';
import React from 'react';
import Layout from '../../../Layout/Layout';
import PregnancyVitals from '../PregnancyVitals/PregnancyVitals';

const currentPregnancy = {
  encounterId: 2121,
  vitals: {},
};

const CurrentPregnancy = (props) => {
  const { location } = props;
  return (
    <Layout location={location} props={props}>
      <PregnancyVitals />
    </Layout>
  );
};

export default CurrentPregnancy;
